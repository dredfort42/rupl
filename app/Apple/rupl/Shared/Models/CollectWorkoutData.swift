//
//  CollectWorkoutData.swift
//  rupl
//
//  Created by Dmitry Novikov on 04/05/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import HealthKit
import os

//	MARK: - Collect workout route data
//
extension WorkoutManager{
	func getWorkoutRoute(results: [HKSample]?) -> [[String: Any]] {
		guard let routes = results as? [HKWorkoutRoute], let route = routes.first else {
			print("No workout routes found or error occurred")
			return []
		}

		var routeData: [[String: Any]] = []

		let routeQuery = HKWorkoutRouteQuery(route: route) { (query, locationsOrNil, done, errorOrNil) in
			guard let locations = locationsOrNil else {
				print("Error fetching locations for route")
				return
			}

			for location in locations {
				let locationData: [String: Any] = [
					"timestamp": location.timestamp.timeIntervalSince1970,
					"latitude": location.coordinate.latitude,
					"longitude": location.coordinate.longitude,
					"horizontalAccuracy": location.horizontalAccuracy,
					"altitude": location.altitude,
					"verticalAccuracy": location.verticalAccuracy,
					"course": location.course,
					"courseAccuracy": location.courseAccuracy,
					"speed": location.speed,
					"speedAccuracy": location.speedAccuracy,
					"ellipsoidalAltitude": location.ellipsoidalAltitude,
					"floor": location.floor ?? 0,
					"sourceInformation": location.sourceInformation ?? ""
				]

				print(location)
				print(locationData)
				routeData.append(locationData)

			}
		}

		self.healthStore.execute(routeQuery)

		print(">> func getWorkoutRoute(results: [HKSample]?) -> [[String: Any]]")
		for r in routeData {
			print(r)
		}

		return routeData
	}
}

//	MARK: - Collect workout data
//
extension WorkoutManager{
	func getWorkoutData() -> [String: Any] {
		var workoutData: [String: Any] = [:]

		let query = HKSampleQuery(
			sampleType: HKObjectType.workoutType(),
			predicate: HKQuery.predicateForWorkouts(with: .running),
			limit: HKObjectQueryNoLimit,
			sortDescriptors: [NSSortDescriptor(key: HKSampleSortIdentifierStartDate, ascending: false)]) { (query, samples, error) in
				guard let samples = samples as? [HKWorkout], let lastRun = samples.first else {
					print("No running workouts found")
					return
				}

				let predicate = HKQuery.predicateForSamples(withStart: lastRun.startDate, end: lastRun.endDate, options: .strictStartDate)

				for type in self.typesToRead {
					let workoutQuery = HKSampleQuery(sampleType: type, predicate: predicate, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { (query, results, error) in
						if error != nil {
							print("Workout query executed with error")
							return
						}

						switch type {
							case HKSeriesType.workoutRoute():
								guard let routes = results as? [HKWorkoutRoute], let route = routes.first else {
									print("No workout routes found or error occurred")
									return
								}

								let routeQuery = HKWorkoutRouteQuery(route: route) { (query, locationsOrNil, done, errorOrNil) in
									guard let locations = locationsOrNil else {
										print("Error fetching locations for route")
										return
									}

									var routeData: [[String: Any]] = []

									for location in locations {
										let locationData: [String: Any] = [
											"timestamp": location.timestamp.timeIntervalSince1970,
											"latitude": location.coordinate.latitude,
											"longitude": location.coordinate.longitude,
											"horizontalAccuracy": location.horizontalAccuracy,
											"altitude": location.altitude,
											"verticalAccuracy": location.verticalAccuracy,
											"course": location.course,
											"courseAccuracy": location.courseAccuracy,
											"speed": location.speed,
											"speedAccuracy": location.speedAccuracy,
											"ellipsoidalAltitude": location.ellipsoidalAltitude,
											"floor": location.floor ?? 0,
											"sourceInformation": location.sourceInformation ?? ""
										]

										routeData.append(locationData)
									}

									if done {
										workoutData["route"] = routeData

										print("\n\n\n")
										print(type.description)
										print(workoutData["route"] ?? "workoutData[\"route\"] empty")
									}
								}

								self.healthStore.execute(routeQuery)

							default:
								guard let samples = results as? [HKQuantitySample] else {
									print("No data found during the run")
									return
								}

								var samplesData: [[String: Any]] = []

								for sample in samples {
									let sampleData: [String: Any] = [
										"timestamp": Int64(sample.startDate.timeIntervalSince1970),
										"quantity": sample.quantity,
										"quantityType": sample.quantityType,
										"description": sample.description,
										"sampleType": sample.sampleType,
										"count": sample.count
									]
									samplesData.append(sampleData)
								}

								workoutData[type.description] = samplesData

								print("\n\n\n")
								print(type.description)
								print(workoutData[type.description] ?? "workoutData[type.description] empty")
						}

					}

					self.healthStore.execute(workoutQuery)
				}
			}

		healthStore.execute(query)


		if let jData = try? JSONSerialization.data(withJSONObject: workoutData, options: .prettyPrinted) {
			if let jsonString = String(data: jData, encoding: .utf8) {
				print("Raw JSON data:")
				print(jsonString)
			} else {
				print("Failed to convert JSON data to string")
			}
		}

		return workoutData
	}
}


//	MARK: - Post workout json data
//
extension WorkoutManager{
	func postWorkoutJsonData(_ jsonData: Data) {
		guard let apiUrl = URL(string: "\(AppSettings.shared.sessionURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)") else {
			print("Invalid URL")
			return
		}
		var request = URLRequest(url: apiUrl)
		request.httpMethod = "POST"
		request.setValue("application/json", forHTTPHeaderField: "Content-Type")

		request.httpBody = jsonData

		let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
			guard let data = data else {
				Logger.shared.log("Error: \(error)")
				return
			}

#if DEBUG
			if let response = response as? HTTPURLResponse {
				print("Response status code: \(response.statusCode)")
				// Handle response if needed
			}

			print(String(data: data, encoding: .utf8)!)
#endif
		}

		task.resume()
	}
}
