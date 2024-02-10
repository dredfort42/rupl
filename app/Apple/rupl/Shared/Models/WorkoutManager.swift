//
//	 WorkoutManager.swift
//	 rupl
//
//	 Created by Dmitry Novikov on 04/01/2024.
//	 Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os
import HealthKit
import SwiftUI

@MainActor
class WorkoutManager: NSObject, ObservableObject {

	struct SessionSateChange {
		let newState: HKWorkoutSessionState
		let date: Date
	}

	//	The workout session live states that the UI observes.
	@Published var sessionState: HKWorkoutSessionState = .notStarted
	@Published var heartRate: Double = 0
//	@Published var activeEnergy: Double = 0
	@Published var distance: Double = 0
//	@Published var power: Double = 0
	@Published var speed: Double = 0
//	@Published var strideLength: Double = 0
//	@Published var verticalOscillation: Double = 0
//	@Published var groundContactTime: Double = 0
//	@Published var vo2Max: Double = 0
//	@Published var stepCount: Double = 0
//	@Published var cadence: Double = 0
//	@Published var elapsedTimeInterval: TimeInterval = 0


	//	HealthKit data types to share
	let typesToShare: Set = [
		HKQuantityType.workoutType(),
		HKSeriesType.workoutRoute()
	]

	//	HealthKit data types to read
	let typesToRead: Set = [
		HKQuantityType(.heartRate), //	count/s, Discrete (Temporally Weighted)
		HKQuantityType(.activeEnergyBurned), //	kcal, Cumulative
		HKQuantityType(.distanceWalkingRunning), //	m, Cumulative
		HKQuantityType(.runningPower), //	W, Discrete (Arithmetic)
		HKQuantityType(.runningSpeed), //	m/s, Discrete (Arithmetic)
		HKQuantityType(.runningStrideLength), //	m, Discrete (Arithmetic)
		HKQuantityType(.runningVerticalOscillation), //	cm, Discrete (Arithmetic)
		HKQuantityType(.runningGroundContactTime), //	ms, Discrete (Arithmetic)
		HKQuantityType(.vo2Max), //	ml/(kg*min), Discrete (Arithmetic)
		HKQuantityType(.stepCount), //	count, Cumulative
		HKQuantityType.workoutType(),
		HKObjectType.activitySummaryType(),
		HKSeriesType.workoutRoute()
	]

//	let parameters = WorkoutParameters()
	let permissibleHorizontalAccuracy: Double = AppSettings.shared.permissibleHorizontalAccuracy
	var useAutoPause: Bool = AppSettings.shared.useAutoPause
	var paceForAutoPause: Double = AppSettings.shared.paceForAutoPause
	var paceForAutoResume: Double = AppSettings.shared.paceForAutoResume
	var pz1NotInZone: Int = AppSettings.shared.pz1NotInZone
	var pz2Easy: Int = AppSettings.shared.pz2Easy
	var pz3FatBurning: Int = AppSettings.shared.pz3FatBurning
	var pz4Aerobic: Int = AppSettings.shared.pz4Aerobic
	var pz5Anaerobic: Int = AppSettings.shared.pz5Anaerobic
	var timeForShowLastSegmentView: Int = AppSettings.shared.timeForShowLastSegmentView
	let sounds = SoundEffects()
	var isTimerStarted: Bool = false
	let timer = DispatchSource.makeTimerSource(queue: DispatchQueue.main)
	let locationManager = LocationManager()
	let motionManager = MotionManager()
	let healthStore = HKHealthStore()
	var workout: HKWorkout?
	var session: HKWorkoutSession?
	var routeBuilder: HKWorkoutRouteBuilder?
#if os(watchOS)
	var builder: HKLiveWorkoutBuilder?
#else
	var contextDate: Date?
#endif
	var isSessionEnded: Bool = false
	var isPauseSetWithButton: Bool = false
	var pauseStartTime: Date = Date()
	var stopTime: Date = Date()

	//	Array for store last 10 speed measurements to colculate average speed
	var last10SpeedMeasurements: [Double] = []
	var last10SpeedMeasurementsSum: Double = 0
	var last10SpeedAverage: Double = 0

	var averageSpeedMetersPerSecond: Double = 0

	var lastSegment: Int = 0
	var lastSegmentStartTime: Date = Date()
	var lastSegmentStopTime: Date = Date()
	var lastSegmentHeartRatesSum: UInt64 = 0
	var lastSegmentHeartRatesCount: UInt = 0
	var lastSegmentViewPresentTime: Int = 0

	var heartRateSum: UInt64 = 0
	var heartRateCount: UInt = 0
	var averageHeartRate: Int = 0
	var heartRateNotificationTimer: Int = 0

	let asynStreamTuple = AsyncStream.makeStream(of: SessionSateChange.self, bufferingPolicy: .bufferingNewest(1))

	static let shared = WorkoutManager()

	private override init() {
		super.init()

		Task {
			do {
				try await healthStore.requestAuthorization(toShare: typesToShare, read: typesToRead)
			} catch {
				Logger.shared.log("Failed to request authorization: \(error)")
			}
		}
		
		Task {
			for await value in asynStreamTuple.stream {
				await consumeSessionStateChange(value)
			}
		}
	}

	private func consumeSessionStateChange(_ change: SessionSateChange) async {
		sessionState = change.newState

#if os(watchOS)
		switch change.newState {
			case .paused:
				pauseStartTime = Date()
			case .running:
					lastSegmentStartTime -= pauseStartTime.timeIntervalSinceNow
			case .stopped:
				timer.cancel()
				locationManager.locationManager.stopUpdatingLocation()

				averageSpeedMetersPerSecond = distance / (builder?.elapsedTime(at: Date()) ?? 1)
				if heartRateCount > 0 {
					averageHeartRate = Int(Double(heartRateSum / UInt64(heartRateCount)) + 0.5)
				}

				let finishedWorkout: HKWorkout?
				stopTime = Date()
				do {
					try await builder?.endCollection(at: change.date)
					finishedWorkout = try await builder?.finishWorkout()
					session?.end()
				} catch {
					Logger.shared.log("Failed to end workout: \(error))")
					return
				}
				workout = finishedWorkout

				guard finishedWorkout != nil else {
					return
				}

				do {
					try await routeBuilder?.finishRoute(with: finishedWorkout!, metadata: nil)
				} catch {
					Logger.shared.log("Failed to associate the route with the workout: \(error)")
					return
				}
			default:
				return
		}
#endif
	}
}

//	MARK: - Workout session management
//
extension WorkoutManager {
	func resetWorkout() {
		sessionState = .notStarted
		heartRate = 0
//		activeEnergy = 0
		distance = 0
//		power = 0
		speed = 0
//		strideLength = 0
//		verticalOscillation = 0
//		groundContactTime = 0
//		vo2Max = 0
//		stepCount = 0
//		cadence = 0
//		elapsedTimeInterval = 0
		workout = nil
		session = nil
		isSessionEnded = false
		isPauseSetWithButton = false
		pauseStartTime = Date()
		stopTime = session?.startDate ?? Date()
		last10SpeedMeasurements = []
		last10SpeedMeasurementsSum = 0
		last10SpeedAverage = 0
		averageSpeedMetersPerSecond = 0
		lastSegment = 0
		lastSegmentStartTime = Date()
		lastSegmentStopTime = Date()
		lastSegmentHeartRatesSum = 0
		lastSegmentHeartRatesCount = 0
		lastSegmentViewPresentTime = 0
		heartRateSum = 0
		heartRateCount = 0
		averageHeartRate = 0
		heartRateNotificationTimer = 0
		routeBuilder = nil
#if os(watchOS)
		builder = nil
#endif
	}
}

//	MARK: - Workout statistics
//
extension WorkoutManager {
	func updateForStatistics(_ statistics: HKStatistics) {
		switch statistics.quantityType {
			case HKQuantityType.quantityType(forIdentifier: .heartRate):
				let heartRateUnit = HKUnit.count().unitDivided(by: .minute())
				heartRate = statistics.mostRecentQuantity()?.doubleValue(for: heartRateUnit) ?? 0
				checkHeartRate()

//			case HKQuantityType.quantityType(forIdentifier: .activeEnergyBurned):
//				let energyUnit = HKUnit.kilocalorie()
//				activeEnergy = statistics.sumQuantity()?.doubleValue(for: energyUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .distanceWalkingRunning):
				let distanceUnit = HKUnit.meter()
				distance = statistics.sumQuantity()?.doubleValue(for: distanceUnit) ?? 0
				checkLastSegment()

//			case HKQuantityType.quantityType(forIdentifier: .runningPower):
//				let powerUnit = HKUnit.watt()
//				power = statistics.mostRecentQuantity()?.doubleValue(for: powerUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .runningSpeed):
				let speedUnit = HKUnit.meter().unitDivided(by: HKUnit.second())
				speed = statistics.mostRecentQuantity()?.doubleValue(for: speedUnit) ?? 0
				checkSpeed()
//				Logger.shared.log("speed: \(self.speed) | avgSpeed: \(self.last10SpeedAverage) | clSpeed: \(self.locationManager.speed)")

//			case HKQuantityType.quantityType(forIdentifier: .runningStrideLength):
//				let lengthUnit = HKUnit.meter()
//				strideLength = statistics.averageQuantity()?.doubleValue(for: lengthUnit) ?? 0
//
//			case HKQuantityType.quantityType(forIdentifier: .runningVerticalOscillation):
//				let verticalOscillationhUnit = HKUnit.meter()
//				verticalOscillation = statistics.sumQuantity()?.doubleValue(for: verticalOscillationhUnit) ?? 0
//
//			case HKQuantityType.quantityType(forIdentifier: .runningGroundContactTime):
//				let contactTimeUnit = HKUnit.secondUnit(with: .milli)
//				groundContactTime = statistics.averageQuantity()?.doubleValue(for: contactTimeUnit) ?? 0
//
//			case HKQuantityType.quantityType(forIdentifier: .vo2Max):
//				let vo2MaxUnit = HKUnit.literUnit(with: .milli).unitDivided(by: HKUnit.gramUnit(with: .kilo).unitDivided(by: HKUnit.minute()))
//				vo2Max = statistics.averageQuantity()?.doubleValue(for: vo2MaxUnit) ?? 0
//
//			case HKQuantityType.quantityType(forIdentifier: .stepCount):
//				stepCount = statistics.mostRecentQuantity()?.doubleValue(for: HKUnit.count()) ?? 0

			default:
				return
		}
	}
}

//	MARK: - HKWorkoutSessionDelegate
//	HealthKit calls the delegate methods on an anonymous serial background queue,
//	so the methods need to be nonisolated explicitly.
//
extension WorkoutManager: HKWorkoutSessionDelegate {
	nonisolated func workoutSession(_ workoutSession: HKWorkoutSession,
									didChangeTo toState: HKWorkoutSessionState,
									from fromState: HKWorkoutSessionState,
									date: Date) {
		Logger.shared.log("Session state changed from \(fromState.rawValue) to \(toState.rawValue)")

		let sessionSateChange = SessionSateChange(newState: toState, date: date)
		asynStreamTuple.continuation.yield(sessionSateChange)
	}

	nonisolated func workoutSession(_ workoutSession: HKWorkoutSession,
									didFailWithError error: Error) {
		Logger.shared.log("\(#function): \(error)")
	}

	nonisolated func workoutSession(_ workoutSession: HKWorkoutSession,
									didDisconnectFromRemoteDeviceWithError error: Error?) {
		Logger.shared.log("\(#function): \(error)")
	}

	nonisolated func workoutSession(_ workoutSession: HKWorkoutSession,
									didReceiveDataFromRemoteWorkoutSession data: [Data]) {
		Logger.shared.log("\(#function): \(data.debugDescription)")
	}
}

//	MARK: - A structure for synchronizing the elapsed time.
//
struct WorkoutElapsedTime: Codable {
	var timeInterval: TimeInterval
	var date: Date
}

//	MARK: - Convenient workout state
//
extension HKWorkoutSessionState {
	var isActive: Bool {
		self != .notStarted && self != .ended
	}
}

