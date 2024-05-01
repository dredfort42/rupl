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
import CoreLocation

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
	var workoutRoute: HKWorkoutRoute?

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
		workoutRoute = nil

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
		print("Save workout start")

		Task {
			do {
				try await builder?.endCollection(at: endDate)

				if let finishedWorkout = try await builder?.finishWorkout() {
					try await routeBuilder?.finishRoute(with: finishedWorkout, metadata: nil)
				}

				postWorkout()
			} catch {
				Logger.shared.log("Failed to end workout: \(error)")
			}
		}


		//		builder?.endCollection(withEnd: endDate) { (_, error) in
		//			print("Ended data collection")
		//			if error != nil {
		//				Logger.shared.log("Failed to end workout: \(error)")
		//			} else {
		//				self.builder?.finishWorkout() { (newWorkout, error) in
		//					print("Finished workout")
		//					guard newWorkout != nil else {
		//						Logger.shared.log("Failed to create workout")
		//						return
		//					}
		//					self.routeBuilder?.finishRoute(with: newWorkout!, metadata: nil) { (newRoute, error) in
		//						print("Finished workout route")
		//						guard newRoute != nil else {
		//							Logger.shared.log("Failed to create workout route")
		//							return
		//						}
		//
		//						let query = HKWorkoutRouteQuery(route: newRoute!) { (query, locationsOrNil, done, errorOrNil) in
		//							if let error = errorOrNil {
		//								Logger.shared.log("Failed to execute workout route query")
		//								return
		//							}
		//
		//							guard let locations = locationsOrNil else {
		//								Logger.shared.log("Failed to get workout route locations")
		//								return
		//							}
		//
		//							if done {
		//								// The query returned all the location data associated with the route.
		//								// Do something with the complete data set.
		//								print(locations)
		//								self.postWorkout(locations: locations)
		//							}
		//
		//							// You can stop the query by calling:
		//							// store.stop(query)
		//
		//						}
		//						self.healthStore.execute(query)
		//
		//					}
		//				}
		//			}
		//		}
#endif
	}
}


//func sampleDataDict(from sample: HKSample) -> [String: Any] {
//	var sampleDict: [String: Any] = [:]
//
//	// Add sample metadata
//	sampleDict["startDate"] = sample.startDate
//	sampleDict["endDate"] = sample.endDate
//
//	// Extract specific data depending on the type
//	if let quantitySample = sample as? HKQuantitySample {
//		sampleDict["quantity"] = quantitySample.quantity.doubleValue(for: HKUnit.count())
//		sampleDict["unit"] = quantitySample.quantity
//	}
//
//	// Add more properties as needed for different types of samples
//
//	return sampleDict
//}

//func fetchHealthData(completion: @escaping ([String: Any]?, Error?) -> Void) {
//	let healthStore = HKHealthStore()
//
//	// Define the types of data you want to retrieve
//	let types: Set<HKObjectType> = [
//		HKObjectType.workoutType(),
////		HKSeriesType.workoutRoute(),
////		HKObjectType.quantityType(forIdentifier: .heartRate)!,
////		HKObjectType.quantityType(forIdentifier: .distanceWalkingRunning)!,
////		HKObjectType.quantityType(forIdentifier: .runningSpeed)!,
////		HKObjectType.quantityType(forIdentifier: .activeEnergyBurned)!,
////		HKObjectType.quantityType(forIdentifier: .bodyMass)!,
////		HKObjectType.quantityType(forIdentifier: .bodyFatPercentage)!,
////		HKObjectType.quantityType(forIdentifier: .bodyMassIndex)!,
////		HKObjectType.quantityType(forIdentifier: .bodyTemperature)!,
////		HKObjectType.quantityType(forIdentifier: .height)!,
////		HKObjectType.quantityType(forIdentifier: .oxygenSaturation)!,
////		HKObjectType.quantityType(forIdentifier: .runningPower)!,
////		HKObjectType.quantityType(forIdentifier: .runningGroundContactTime)!,
////		HKObjectType.quantityType(forIdentifier: .runningStrideLength)!,
////		HKObjectType.quantityType(forIdentifier: .runningVerticalOscillation)!,
////		HKObjectType.quantityType(forIdentifier: .stepCount)!,
////		HKObjectType.quantityType(forIdentifier: .vo2Max)!,
////		HKObjectType.quantityType(forIdentifier: .flightsClimbed)!
//	]
//
//
//	var healthData: [String: Any] = [:]
//
//	// Fetch data for each type
//	for type in types {
//		let query = HKSampleQuery(sampleType: type as! HKSampleType, predicate: nil, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { (query, samples, error) in
//			if let samples = samples {
//				var sampleData: [[String: Any]] = []
//				for sample in samples {
//					// Convert each sample to a dictionary representation
//												let sampleDict = sampleDataDict(from: sample)
//												sampleData.append(sampleDict)
//				}
//
//				print(type.identifier)
//				print(sampleData)
//
//				healthData[type.identifier] = sampleData
//			} else {
//				healthData[type.identifier] = []
//			}
//		}
//		healthStore.execute(query)
//	}
//
//	completion(healthData, nil)
//}






//	MARK: - Post workout
//
extension WorkoutManager{
	func postWorkout() {
#if DEBUG
		print("postWorkout()")
#endif

		let query = HKSampleQuery(
			sampleType: HKObjectType.workoutType(),
			predicate: HKQuery.predicateForWorkouts(with: .running),
			limit: HKObjectQueryNoLimit,
			sortDescriptors: [NSSortDescriptor(key: HKSampleSortIdentifierStartDate, ascending: false)]) { (query, samples, error) in
				guard let samples = samples as? [HKWorkout], let lastRun = samples.first else {
					print("No running workouts found")
					return
				}

				// Access data from the last run
//				print("Last run:")
//				print("Start Date: \(lastRun.startDate)")
//				print("End Date: \(lastRun.endDate)")
//				print("Distance: \(lastRun.totalDistance?.doubleValue(for: HKUnit.meter()) ?? 0) meters")
//				print("Duration: \(lastRun.duration)")

				let predicate = HKQuery.predicateForSamples(withStart: lastRun.startDate, end: lastRun.endDate, options: .strictStartDate)

				// Define query for heart rate samples during the run
				for type in self.typesToRead {
					let heartRateQuery = HKSampleQuery(sampleType: type, predicate: predicate, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { (query, samples, error) in
						guard let samples = samples as? [HKQuantitySample] else {
							print("No data found during the run")
							return
						}

						print(type.description)
						print("Data during the run:")
						for sample in samples {
							print("Date: \(Int64(sample.startDate.timeIntervalSince1970)), data: \(sample.quantity)")
						}
					}

					// Execute the heart rate query
					self.healthStore.execute(heartRateQuery)
				}
			}

		// Execute the query
		healthStore.execute(query)
		//		let query = HKSampleQuery(sampleType: HKQuantityType(.heartRate), predicate: nil, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { (query, samples, error) in


		//			var healthData: [String: Any] = [:]
		//
		//			// Fetch data for each type
		//			for type in typesToRead {
		//				let query = HKSampleQuery(sampleType: type, predicate: nil, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { (query, samples, error) in
		//					if let samples = samples {
		//						var sampleData: [[String: Any]] = []
		//						for sample in samples {
		//							// Convert each sample to a dictionary representation
		//														let sampleDict = sampleDataDict(from: sample)
		//														sampleData.append(sampleDict)
		//						}
		//
		//						print(type.identifier)
		//						print(sampleData)
		//
		//						healthData[type.identifier] = sampleData
		//					} else {
		//						healthData[type.identifier] = []
		//					}
		//				}
		//				healthStore.execute(query)
		//			}
		//
		//			print(healthData)

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
