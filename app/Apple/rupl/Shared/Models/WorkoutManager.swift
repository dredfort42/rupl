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

	//	Environment
	let timerManager = TimerManager()
	let healthStore = HKHealthStore()
	var session: HKWorkoutSession?
	var routeBuilder: HKWorkoutRouteBuilder?
#if os(watchOS)
	var builder: HKLiveWorkoutBuilder?
#else
	var contextDate: Date?
#endif

	//	Array for store last 10 speed measurements to colculate average speed
	var last10SpeedMeasurements: [Double] = []
	var last10SpeedMeasurementsSum: Double = 0
	var last10SpeedAverage: Double = 0

	//	Summary dataSession state changed
	var summaryHeartRateSum: UInt64 = 0
	var summaryHeartRateCount: UInt = 0

	//	Segment data
	var segmentStartTime: Date = Date()
	var segmentFinishTime: Date = Date()
	var segmentNumber: Int = 0
	var segmentHeartRatesSum: UInt64 = 0
	var segmentHeartRatesCount: UInt = 0

	//	Pause data
	var isPauseSetWithButton: Bool = false
	var pauseStartTime: Date = Date()

	let asynStreamTuple = AsyncStream.makeStream(of: SessionSateChange.self, bufferingPolicy: .bufferingNewest(1))

	static let shared = WorkoutManager()

	private override init() {
		super.init()

		LocationManager.shared.start()
		MotionManager.shared.start()

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
				segmentStartTime -= pauseStartTime.timeIntervalSinceNow
			default:
				return
		}
	}
}

//	MARK: - Workout session reset workout
//
extension WorkoutManager {
	func resetWorkout() {
#if DEBUG
		print("resetWorkout()")
#endif
		sessionState = .notStarted
		heartRate = 0
		distance = 0
		speed = 0

		LocationManager.shared.start()
		MotionManager.shared.start()
		timerManager.start(timeInterval: 1, repeats: true, action: autoPause)
		timerManager.start(timeInterval: 1, repeats: true, action: addLocationsToRoute)
		timerManager.start(timeInterval: TimeInterval(AppSettings.shared.soundNotificationTimeOut), repeats: true, action: checkHeartRate)

		session = nil
		routeBuilder = nil
#if os(watchOS)
		builder = nil
#endif

		//	Array for store last 10 speed measurements to colculate average speed
		last10SpeedMeasurements = []
		last10SpeedMeasurementsSum = 0
		last10SpeedAverage = 0

		//	Summary data
		summaryHeartRateSum = 0
		summaryHeartRateCount = 0

		//	Segment data
		segmentStartTime = session?.startDate ?? Date()
		segmentFinishTime = session?.startDate ?? Date()
		segmentNumber = 0
		segmentHeartRatesSum = 0
		segmentHeartRatesCount = 0

		//	Pause data
		isPauseSetWithButton = false
		pauseStartTime = session?.startDate ?? Date()
	}
}

//	MARK: - Start workout
//
extension WorkoutManager {
	func startWorkout() {
#if DEBUG
		print("startWorkout()")
#endif
		resetWorkout()
#if os(watchOS)
		Task {
			do {
				let configuration = HKWorkoutConfiguration()
				configuration.activityType = .running
				configuration.locationType = .outdoor
				try await startWorkout(workoutConfiguration: configuration)
			} catch {
				Logger.shared.log("Failed to start workout \(error))")
			}
		}
#endif
	}
}

//	MARK: - Stop workout
//
extension WorkoutManager {
	func finishWorkout() {
#if DEBUG
		print("finishWorkout()")
#endif
		sessionState = .stopped
		session?.stopActivity(with: .now)
	}
}

//	MARK: - Save workout
//
extension WorkoutManager{
	func saveWorkout() {
#if DEBUG
		print("saveWorkout()")
#endif
		timerManager.stop()
		LocationManager.shared.stop()
		MotionManager.shared.stop()
		session?.end()

#if os(watchOS)
		Task {
			do {
				try await builder?.endCollection(at: session?.endDate ?? Date())
				if let finishedWorkout = try await builder?.finishWorkout() {
					try await routeBuilder?.finishRoute(with: finishedWorkout, metadata: nil)
				}
			} catch {
				Logger.shared.log("Failed to end workout: \(error))")
			}
		}
#endif
		 DispatchQueue.main.asyncAfter(deadline: .now() + 0.25) {
			 self.sessionState = .notStarted
		 }
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
#if DEBUG
		print("Session state changed from \(fromState.rawValue) to \(toState.rawValue)")
#endif

		let sessionSateChange = SessionSateChange(newState: toState, date: date)
		asynStreamTuple.continuation.yield(sessionSateChange)
	}

	nonisolated func workoutSession(_ workoutSession: HKWorkoutSession,
									didFailWithError error: Error) {
		Logger.shared.log("\(#function): \(error)")
	}
}

//	MARK: - A structure for synchronizing the elapsed time.
//
struct WorkoutElapsedTime: Codable {
	var timeInterval: TimeInterval
	var date: Date
}
