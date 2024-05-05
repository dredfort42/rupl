//
//  WorkoutData.swift
//  rupl
//
//  Created by Dmitry Novikov on 04/05/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import HealthKit
import os

class WorkoutData {
	let isHealthDataAvailable: Bool = HKHealthStore.isHealthDataAvailable()
	let healthStore = HKHealthStore()
	var predicate: NSPredicate?


	private func getPredicate(completion: @escaping (NSPredicate?, Error?) -> Void) {
		if !isHealthDataAvailable {
			completion(nil, nil)
			return
		}

		let query = HKSampleQuery(
			sampleType: HKObjectType.workoutType(),
			predicate: HKQuery.predicateForWorkouts(with: .running),
			limit: HKObjectQueryNoLimit,
			sortDescriptors: [NSSortDescriptor(key: HKSampleSortIdentifierStartDate, ascending: false)]
		) { (query, samples, error) in
			if let error = error {
				Logger.shared.log("Error: \(error)")
				completion(nil, error)
				return
			}

			guard let samples = samples as? [HKWorkout], let lastRun = samples.first else {
				Logger.shared.log("No samples found")
				completion(nil, nil)
				return
			}

			let predicate = HKQuery.predicateForSamples(withStart: lastRun.startDate, end: lastRun.endDate, options: .strictStartDate)
			completion(predicate, nil)
		}

		healthStore.execute(query)
	}

	private func getWorkoutRoute(results: [HKSample]?, completion: @escaping ([[String: Any]]?) -> Void) {
		guard let routes = results as? [HKWorkoutRoute], let route = routes.first else {
			Logger.shared.log("No workout routes found or error occurred")
			completion(nil)
			return
		}

		var routeData: [[String: Any]] = []

		let routeQuery = HKWorkoutRouteQuery(route: route) { (query, locationsOrNil, done, errorOrNil) in
			guard let locations = locationsOrNil else {
				Logger.shared.log("Error fetching locations for route")
				completion(nil)
				return
			}

			for location in locations {
				let locationData: [String: Any] = [
					"timestamp": Int64(location.timestamp.timeIntervalSince1970),
					"latitude": location.coordinate.latitude,
					"longitude": location.coordinate.longitude,
					"horizontalAccuracy": location.horizontalAccuracy,
					"altitude": location.altitude,
					"verticalAccuracy": location.verticalAccuracy,
					"course": location.course,
					"courseAccuracy": location.courseAccuracy,
					"speed": location.speed,
					"speedAccuracy": location.speedAccuracy,
					"ellipsoidalAltitude": location.ellipsoidalAltitude
				]

				routeData.append(locationData)
			}

			completion(routeData)
		}

		self.healthStore.execute(routeQuery)
	}

	private func getWorkoutData() {
		getPredicate { predicate, error in
			if let error = error {
				Logger.shared.log("Error: \(error)")
				return
			}

			guard let predicate = predicate else {
				Logger.shared.log("No predicate")
				return
			}

			for type in WorkoutManager.shared.typesToRead {
				let workoutQuery = HKSampleQuery(sampleType: type, predicate: predicate, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { (query, results, error) in
					if error != nil {
						Logger.shared.log("Error: \(error)")
						return
					}

					switch type {
						case HKSeriesType.workoutRoute():
							self.getWorkoutRoute(results: results) { route in
								guard let route = route else {
									Logger.shared.log("No workout routes found or error occurred")
									return
								}

								print("\n\n\n")
								print(type.description)
//								print(route)
								self.serialization(data: [type.description: route])
							}
							
//						case HKQuantityType(.activeEnergyBurned):
//							guard let samples = results as? [HKQuantitySample] else {
//								print("No data found during the run")
//								return
//							}
//
//							var samplesData: [[String: Any]] = []
//
//							for sample in samples {
//
//								let sampleData: [String: Any] = [
//									"timestamp": Int64(sample.startDate.timeIntervalSince1970),
//									"quantity": sample.quantity.description,
//									"quantityType": sample.quantityType.description,
//									"description": sample.description,
//									"sampleType": sample.sampleType.description,
//									"count": sample.count
//								]
//								samplesData.append(sampleData)
//							}
//
//							print("\n\n\n")
//							print(type.description)
////							print(samplesData)
//							self.serialization(data: [type.description: samplesData])

						default:
							guard let samples = results as? [HKQuantitySample] else {
								print("No data found during the run")
								return
							}

							var samplesData: [[String: Any]] = []

							for sample in samples {
								let sampleData: [String: Any] = [
									"timestamp": Int64(sample.startDate.timeIntervalSince1970),
									"quantity": sample.quantity.description
//									"quantityType": sample.quantityType.description,
//									"description": sample.description,
//									"sampleType": sample.sampleType.description,
//									"count": sample.count
								]
								samplesData.append(sampleData)
							}

							//							workoutData[type.description] = samplesData

							print("\n\n\n")
							print(type.description)
//							print(samplesData)
							self.serialization(data: [type.description: samplesData])

					}

				}

				self.healthStore.execute(workoutQuery)
			}
		}
	}

	private func serialization(data: [String: Any]) -> Data? {
		if data.isEmpty {
			return nil
		}

		if let jData = try? JSONSerialization.data(withJSONObject: data, options: .prettyPrinted) {
			if let jsonString = String(data: jData, encoding: .utf8) {
				print("Raw JSON data:")
				print(jsonString)
			} else {
				print("Failed to convert JSON data to string")
			}

			return jData
		}

		return nil
	}

	func postWorkout() {
#if DEBUG
		print("postWorkout()")
#endif
		getWorkoutData()
	}

	static let shared = WorkoutData()
}

//
////	MARK: - Collect workout route data
////
//extension WorkoutManager{

//}
//
//// MARK: - 	Get start and stop Date of last run workout
////
//extension WorkoutManager{
////	func getLastRunPredicate(completion: @escaping (NSPredicate?, Error?) -> Void) async {
////		let query = HKSampleQuery(
////			sampleType: HKObjectType.workoutType(),
////			predicate: HKQuery.predicateForWorkouts(with: .running),
////			limit: HKObjectQueryNoLimit,
////			sortDescriptors: [NSSortDescriptor(key: HKSampleSortIdentifierStartDate, ascending: false)]) { (query, samples, error) in
////				guard let samples = samples as? [HKWorkout], let lastRun = samples.first, error == nil else {
////					completion(nil, error)
////					return
////				}
////
////				completion(HKQuery.predicateForSamples(withStart: lastRun.startDate, end: lastRun.endDate, options: .strictStartDate), nil)
////			}
////
////		healthStore.execute(query)
////	}
//
//	// Asynchronous function to perform the workout query and return the result
//	func fetchWorkoutsAsync(completion: @escaping ([HKWorkout]?, Error?) -> Void) {
//		// Predicate to specify any filtering criteria (optional)
//		let predicate = NSPredicate(format: "duration > %@", NSNumber(value: 3)) // Example: Filter workouts with duration greater than 300 seconds (5 minutes)
//
//		// Create a query to retrieve workouts matching the specified type and predicate
//		let workoutQuery = HKSampleQuery(sampleType: HKObjectType.workoutType(), predicate: predicate, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { (query, samples, error) in
//			guard let retrievedWorkouts = samples as? [HKWorkout], error == nil else {
//				completion(nil, error)
//				return
//			}
//
//			// Return the retrieved workouts
//			completion(retrievedWorkouts, nil)
//		}
//
//		// Execute the workout query
//		healthStore.execute(workoutQuery)
//	}
//
//}
//
////	MARK: - Collect workout data
////
//extension WorkoutManager{
//	func getWorkoutData() -> [String: Any] {
//
////		var predicate: NSPredicate?
//		var workoutData: [String: Any] = [:]
//
//		var workouts: [HKWorkout]?
//
//		fetchWorkoutsAsync { retrievedWorkouts, error in
//			if let error = error {
//				print("Error retrieving workouts: \(error)")
//				return
//			}
//
//			workouts = retrievedWorkouts
//		}
//
//
//		print("---PRINTER---")
//		// Process the retrieved workouts
//		for workout in workouts ?? [] {
//			print("Workout: \(workout)")
//
//		}
//
////		Task {
////			do {
////				try await getLastRunPredicate() { (data, error) in
////					if error != nil {
////						Logger.shared.log("Error in getLastRunPredicate(): \(error)")
////						return
////					}
////					predicate = data
////				}
////			} catch {
//////				Logger.shared.log("Error: \(error)")
////			}
////		}
//
////		let query = HKSampleQuery(
////			sampleType: HKObjectType.workoutType(),
////			predicate: HKQuery.predicateForWorkouts(with: .running),
////			limit: HKObjectQueryNoLimit,
////			sortDescriptors: [NSSortDescriptor(key: HKSampleSortIdentifierStartDate, ascending: false)]) { (query, samples, error) in
////				guard let samples = samples as? [HKWorkout], let lastRun = samples.first else {
////					print("No running workouts found")
////					return
////				}
////
////				let predicate = HKQuery.predicateForSamples(withStart: lastRun.startDate, end: lastRun.endDate, options: .strictStartDate)
////
////				for type in self.typesToRead {
////					let workoutQuery = HKSampleQuery(sampleType: type, predicate: predicate, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { (query, results, error) in
////						if error != nil {
////							print("Workout query executed with error")
////							return
////						}
////
////						switch type {
////							case HKSeriesType.workoutRoute():
//////								workoutData["route"] = self.getWorkoutRoute(results: results)
////
////								guard let routes = results as? [HKWorkoutRoute], let route = routes.first else {
////									print("No workout routes found or error occurred")
////									return
////								}
////
////								let routeQuery = HKWorkoutRouteQuery(route: route) { (query, locationsOrNil, done, errorOrNil) in
////									guard let locations = locationsOrNil else {
////										print("Error fetching locations for route")
////										return
////									}
////
////									var routeData: [[String: Any]] = []
////
////									for location in locations {
////										let locationData: [String: Any] = [
////											"timestamp": location.timestamp.timeIntervalSince1970,
////											"latitude": location.coordinate.latitude,
////											"longitude": location.coordinate.longitude,
////											"horizontalAccuracy": location.horizontalAccuracy,
////											"altitude": location.altitude,
////											"verticalAccuracy": location.verticalAccuracy,
////											"course": location.course,
////											"courseAccuracy": location.courseAccuracy,
////											"speed": location.speed,
////											"speedAccuracy": location.speedAccuracy,
////											"ellipsoidalAltitude": location.ellipsoidalAltitude,
////											"floor": location.floor ?? 0,
////											"sourceInformation": location.sourceInformation ?? ""
////										]
////
////										routeData.append(locationData)
////									}
////
////									if done {
////										workoutData["route"] = routeData
////
//////										print("\n\n\n")
//////										print(type.description)
//////										print(workoutData["route"] ?? "workoutData[\"route\"] empty")
////									}
////								}
////
////								self.healthStore.execute(routeQuery)
////
////							default:
////								guard let samples = results as? [HKQuantitySample] else {
////									print("No data found during the run")
////									return
////								}
////
////								var samplesData: [[String: Any]] = []
////
////								for sample in samples {
////									let sampleData: [String: Any] = [
////										"timestamp": Int64(sample.startDate.timeIntervalSince1970),
////										"quantity": sample.quantity,
////										"quantityType": sample.quantityType,
////										"description": sample.description,
////										"sampleType": sample.sampleType,
////										"count": sample.count
////									]
////									samplesData.append(sampleData)
////								}
////
////								workoutData[type.description] = samplesData
////
//////								print("\n\n\n")
//////								print(type.description)
//////								print(workoutData[type.description] ?? "workoutData[type.description] empty")
////						}
////
////					}
////
////					self.healthStore.execute(workoutQuery)
////				}
////			}
////
////		healthStore.execute(query)
//
//
//
////
//
//		return workoutData
//	}
//}
//
//
////	MARK: - Post workout json data
////
//extension WorkoutManager{
//	func postWorkoutJsonData(_ jsonData: Data) {
//		guard let apiUrl = URL(string: "\(AppSettings.shared.sessionURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)") else {
//			print("Invalid URL")
//			return
//		}
//		var request = URLRequest(url: apiUrl)
//		request.httpMethod = "POST"
//		request.setValue("application/json", forHTTPHeaderField: "Content-Type")
//
//		request.httpBody = jsonData
//
//		let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
//			guard let data = data else {
//				Logger.shared.log("Error: \(error)")
//				return
//			}
//
//#if DEBUG
//			if let response = response as? HTTPURLResponse {
//				print("Response status code: \(response.statusCode)")
//				// Handle response if needed
//			}
//
//			print(String(data: data, encoding: .utf8)!)
//#endif
//		}
//
//		task.resume()
//	}
//}
