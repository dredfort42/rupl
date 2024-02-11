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
	@Published var distance: Double = 0
	@Published var speed: Double = 0

	//	HealthKit data types to share
	let typesToShare: Set = [
		HKQuantityType.workoutType(),
		HKSeriesType.workoutRoute()
	]

	//	HealthKit data types to read
	let typesToRead: Set = [
		HKQuantityType(.heartRate), //	count/s, Discrete (Temporally Weighted)
		HKQuantityType(.distanceWalkingRunning), //	m, Cumulative
		HKQuantityType(.runningSpeed), //	m/s, Discrete (Arithmetic)
	]

	//	Settings
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

	//	Environment
	let sounds = SoundEffects()
	let timerManager = TimerManager()
	let locationManager = LocationManager()
	let motionManager = MotionManager()
	let healthStore = HKHealthStore()

	var session: HKWorkoutSession?
	var routeBuilder: HKWorkoutRouteBuilder?
#if os(watchOS)
	var builder: HKLiveWorkoutBuilder?
#else
	var contextDate: Date?
#endif

	//	Timer
	var heartRateNotificationTimer: Int = 0
	var lastSegmentViewPresentTimer: Int = 0

	var isPauseSetWithButton: Bool = false
	var pauseStartTime: Date = Date()
	var stopTime: Date = Date()

	//	Array for store last 10 speed measurements to colculate average speed
	var last10SpeedMeasurements: [Double] = []
	var last10SpeedMeasurementsSum: Double = 0
	var last10SpeedAverage: Double = 0

	//	Segment data
	var lastSegment: Int = 0
	var lastSegmentStartTime: Date = Date()
	var lastSegmentStopTime: Date = Date()
	var lastSegmentHeartRatesSum: UInt64 = 0
	var lastSegmentHeartRatesCount: UInt = 0

	//	Summary data
	var heartRateSum: UInt64 = 0
	var heartRateCount: UInt = 0
	var averageHeartRate: Int = 0
	var averageSpeedMetersPerSecond: Double = 0

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

		switch change.newState {
			case .paused:
				pauseStartTime = Date()
			case .running:
				lastSegmentStartTime -= pauseStartTime.timeIntervalSinceNow
			case .stopped:
				stopTime = Date()
				timerManager.stop()
				locationManager.stop()
				motionManager.stop()

				if heartRateCount > 0 {
					averageHeartRate = Int(Double(heartRateSum / UInt64(heartRateCount)) + 0.5)
				}

				let time: Double = stopTime.timeIntervalSince(session?.startDate ?? Date())
				if  time > 0 {
					averageSpeedMetersPerSecond = distance / time
				}

				await finishWorkout()
			default:
				return
		}
	}
}

//	MARK: - Workout session reset workout
//
extension WorkoutManager {
	func resetWorkout() {
		sessionState = .notStarted
		heartRate = 0
		distance = 0
		speed = 0

		timerManager.start(timeInterval: 1, repeats: true, action: timerActions)
		locationManager.start()
		motionManager.start()

		session = nil
		routeBuilder = nil
#if os(watchOS)
		builder = nil
#endif

		//	Timer
		heartRateNotificationTimer = 0
		lastSegmentViewPresentTimer = 0

//		isSessionEnded = false
		isPauseSetWithButton = false
		pauseStartTime = Date()
		stopTime = Date()

		//	Array for store last 10 speed measurements to colculate average speed
		last10SpeedMeasurements = []
		last10SpeedMeasurementsSum = 0
		last10SpeedAverage = 0

		//	Segment data
		lastSegment = 0
		lastSegmentStartTime = Date()
		lastSegmentStopTime = Date()
		lastSegmentHeartRatesSum = 0
		lastSegmentHeartRatesCount = 0

		//	Summary data
		heartRateSum = 0
		heartRateCount = 0
		averageHeartRate = 0
		averageSpeedMetersPerSecond = 0
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

			case HKQuantityType.quantityType(forIdentifier: .distanceWalkingRunning):
				let distanceUnit = HKUnit.meter()
				distance = statistics.sumQuantity()?.doubleValue(for: distanceUnit) ?? 0
				checkLastSegment()

			case HKQuantityType.quantityType(forIdentifier: .runningSpeed):
				let speedUnit = HKUnit.meter().unitDivided(by: HKUnit.second())
				speed = statistics.mostRecentQuantity()?.doubleValue(for: speedUnit) ?? 0
				checkSpeed()

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

