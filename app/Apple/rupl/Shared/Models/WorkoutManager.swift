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
		HKQuantityType(.heartRate),
		HKQuantityType(.distanceWalkingRunning),
		HKQuantityType(.runningSpeed),
		HKQuantityType(.activeEnergyBurned),
		HKQuantityType(.bodyMass),
		HKQuantityType(.bodyFatPercentage),
		HKQuantityType(.bodyMassIndex),
		HKQuantityType(.bodyTemperature),
		HKQuantityType(.height),
		HKQuantityType(.oxygenSaturation),
		HKQuantityType(.runningPower),
		HKQuantityType(.runningGroundContactTime),
		HKQuantityType(.runningStrideLength),
		HKQuantityType(.runningVerticalOscillation),
		HKQuantityType(.stepCount),
		HKQuantityType(.vo2Max),
		HKQuantityType(.flightsClimbed)
	]

	//	Environment
	let timerManager = TimerManager()
	let healthStore = HKHealthStore()
	var session: HKWorkoutSession?
#if os(watchOS)
	var builder: HKLiveWorkoutBuilder?
#else
	var contextDate: Date?
#endif
	var routeBuilder: HKWorkoutRouteBuilder?

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
		timerManager.start(timeInterval: 1, repeats: true, action: TaskManager.shared.runSession)
		timerManager.start(timeInterval: 1, repeats: true, action: checkIntervalTimeLeft)
		timerManager.start(timeInterval: TimeInterval(AppSettings.shared.soundNotificationTimeOut), repeats: true, action: checkHeartRate)

		session = nil
#if os(watchOS)
		builder = nil
#endif
		routeBuilder = nil

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

		guard let endDate = session?.endDate else {
			Logger.shared.log("End date for the workout session is nil")
			return
		}

#if os(watchOS)
		Task {
			do {
				try await builder?.endCollection(at: endDate)

				if let finishedWorkout = try await builder?.finishWorkout() {
					try await routeBuilder?.finishRoute(with: finishedWorkout, metadata: nil)
				}

				// Change state after route saving is complete
				self.sessionState = .notStarted
			} catch {
				Logger.shared.log("Failed to end workout: \(error)")
			}
		}
#endif
	}
}

//	MARK: - Post workout
//
extension WorkoutManager{
	func postWorkout() async {
#if DEBUG
		print("postWorkout()")
#endif



//		let query = HKSampleQuery(sampleType: HKObjectType.workoutType(), predicate: nil, limit: 1, sortDescriptors: [NSSortDescriptor(key: HKSampleSortIdentifierEndDate, ascending: false)]) { (query, samples, error) in
//			guard let workoutSamples = samples as? [HKWorkout], let runWorkout = workoutSamples.first(where: { $0.workoutActivityType == .running }) else {
//				print("No running workout data available.")
//				return
//			}
//
//			// Extract relevant data from the workout
//			let workoutStartDate = runWorkout.startDate
//			let workoutEndDate = runWorkout.endDate
//			let totalDistance = runWorkout.totalDistance?.doubleValue(for: .meter())
//			let totalEnergyBurned = runWorkout.totalEnergyBurned?.doubleValue(for: .kilocalorie())
//
//			print(workoutStartDate)
//			print(workoutEndDate)
//			print(totalDistance ?? 0.0)
//			print(totalEnergyBurned ?? 0.0)
//
//
//			//			// Prepare the data to send to the server
//			//			let dataToSend: [String: Any] = [
//			//				"startDate": workoutStartDate,
//			//				"endDate": workoutEndDate,
//			//				"totalDistance": totalDistance ?? 0.0,
//			//				"totalEnergyBurned": totalEnergyBurned ?? 0.0
//			//				// Add more data as needed
//			//			]
//			//
//			//			// Send data to server
//			//			print("sendDataToServer(data: dataToSend)")
//			let runningObjectQuery = HKQuery.predicateForObjects(from: runWorkout)
//
//
//			let routeQuery = HKAnchoredObjectQuery(type: HKSeriesType.workoutRoute(), predicate: runningObjectQuery, anchor: nil, limit: HKObjectQueryNoLimit) { (query, samples, deletedObjects, anchor, error) in
//
//				guard error == nil else {
//					// Handle any errors here.
//					fatalError("The initial query failed.")
//				}
//
//				// Process the initial route data here.
//			}
//
//
//			routeQuery.updateHandler = { (query, samples, deleted, anchor, error) in
//
//				guard error == nil else {
//					// Handle any errors here.
//					fatalError("The update failed.")
//				}
//
//				// Process updates or additions here.
//			}
//
//
//			self.healthStore.execute(routeQuery)
//		}
//		healthStore.execute(query)




//		let routeQuery = HKSampleQuery(sampleType: HKSeriesType.workoutRoute(), predicate: nil, limit: 1, sortDescriptors: [NSSortDescriptor(key: HKSampleSortIdentifierStartDate, ascending: false)]) { (query, samples, error) in
//			guard let workoutRouteSamples = samples as? [HKWorkoutRoute], let workoutRoute = workoutRouteSamples.first else {
//				print("No workout route data available.")
//				return
//			}
//
//			let routeQueryDescriptor = HKWorkoutRouteQueryDescriptor(workoutRoute)
////
////			// Get the AsyncSequence that returns individual locations.
//			let locations = routeQueryDescriptor.results(for: self.healthStore)
////
////
//// Access each location.
//			for try await location in locations {
//
//				// Process the results here.
//				print(location.coordinate)
//				print(location.timestamp)
//			}
////			print("---")
////			print(locations)
//
//		}
//		healthStore.execute(routeQuery)

		let workouts = await readWorkouts()
//		guard let workout = workouts?.first else {
//			return
//		}
		print("---workout---")
		print(workouts ?? "")

//		let routes = await getWorkoutRoute(workout: workout)
//		guard let route = routes?.first else {
//			return
//		}
//		print("---route---")
//		print(route)

//		let locations = await getLocationDataForRoute(givenRoute: route)
//		print("---locations---")
//		print(locations)

	}

	func readWorkouts() async -> [HKWorkout]? {
		let cycling = HKQuery.predicateForWorkouts(with: .cycling)

		let samples = try! await withCheckedThrowingContinuation { (continuation: CheckedContinuation<[HKSample], Error>) in
			healthStore.execute(HKSampleQuery(sampleType: .workoutType(), predicate: cycling, limit: HKObjectQueryNoLimit,sortDescriptors: [.init(keyPath: \HKSample.startDate, ascending: false)], resultsHandler: { query, samples, error in
				if let hasError = error {
					continuation.resume(throwing: hasError)
					return
				}

				guard let samples = samples else {
					Logger.shared.log("Can't read workout from Health Storage")
					return
				}

				continuation.resume(returning: samples)
			}))
		}

		guard let workouts = samples as? [HKWorkout] else {
			return nil
		}

		return workouts
	}

	func getWorkoutRoute(workout: HKWorkout) async -> [HKWorkoutRoute]? {
		let byWorkout = HKQuery.predicateForObjects(from: workout)

		let samples = try! await withCheckedThrowingContinuation { (continuation: CheckedContinuation<[HKSample], Error>) in
			healthStore.execute(HKAnchoredObjectQuery(type: HKSeriesType.workoutRoute(), predicate: byWorkout, anchor: nil, limit: HKObjectQueryNoLimit, resultsHandler: { (query, samples, deletedObjects, anchor, error) in
				if let hasError = error {
					continuation.resume(throwing: hasError)
					return
				}

				guard let samples = samples else {
					Logger.shared.log("Can't read workout route from Health Storage")
					return
				}

				continuation.resume(returning: samples)
			}))
		}

		guard let routs = samples as? [HKWorkoutRoute] else {
			return nil
		}

		return routs
	}

//	func getLocationDataForRoute(givenRoute: HKWorkoutRoute) async -> [CLLocation] {
//		let locations = try! await withCheckedThrowingContinuation { (continuation: CheckedContinuation<[CLLocation], Error>) in
//			var allLocations: [CLLocation] = []
//
//			// Create the route query.
//			let query = HKWorkoutRouteQuery(route: givenRoute) { (query, locationsOrNil, done, errorOrNil) in
//
//				if let error = errorOrNil {
//					continuation.resume(throwing: error)
//					return
//				}
//
//				guard let currentLocationBatch = locationsOrNil else {
//					Logger.shared.log("Can't read locations from Workout Route")
//					return
//				}
//
//				allLocations.append(contentsOf: currentLocationBatch)
//
//				if done {
//					continuation.resume(returning: allLocations)
//				}
//			}
//
//			healthStore.execute(query)
//		}
//
//		return locations
//	}

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
				checkIntervalDistanceLeft()

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
