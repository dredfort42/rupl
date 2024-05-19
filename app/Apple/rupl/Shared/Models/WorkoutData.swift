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
	var compleatedQuery: UInt8 = 0
	var workoutData: [String: Any] = [:]

	private func getWorkoutDataType(type: HKSampleType) -> String {
		switch type {
			case HKSeriesType.workoutRoute():
				return "route_data"
			case HKQuantityType(.heartRate):
				return "heart_rate"
			case HKQuantityType(.distanceWalkingRunning):
				return "distance"
			case HKQuantityType(.runningSpeed):
				return "speed"
			case HKQuantityType(.activeEnergyBurned):
				return "energy_burned"
			case HKQuantityType(.runningPower):
				return "running_power"
			case HKQuantityType(.runningGroundContactTime):
				return "ground_contact_time"
			case HKQuantityType(.runningStrideLength):
				return "stride_length"
			case HKQuantityType(.runningVerticalOscillation):
				return "vertical_oscillation"
			case HKQuantityType(.stepCount):
				return "step_count"
			case HKQuantityType(.vo2Max):
				return "vo2_max"
			default:
				return "unknown"
		}
	}

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

			let sessionData: [String: Any] = [
				"session_start_time": Int64(lastRun.startDate.timeIntervalSince1970),
				"session_end_time": Int64(lastRun.endDate.timeIntervalSince1970),
				"email": AppSettings.shared.userEmail
			]

			self.workoutData["session"] = sessionData

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
					"horizontal_accuracy": location.horizontalAccuracy,
					"altitude": location.altitude,
					"vertical_accuracy": location.verticalAccuracy,
					"course": location.course,
					"course_accuracy": location.courseAccuracy,
					"speed": location.speed,
					"speed_accuracy": location.speedAccuracy
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

								self.workoutData[self.getWorkoutDataType(type: type)] = route
								self.compleatedQuery += 1
							}

						case HKQuantityType.workoutType():
							self.compleatedQuery += 1

						default:
							guard let samples = results as? [HKQuantitySample] else {
								print("No data found during the run")
								return
							}

							var sampleData: [[String: Any]] = []

							for sample in samples {
								let data: [String: Any] = [
									"timestamp": Int64(sample.startDate.timeIntervalSince1970),
									"quantity": sample.quantity.description
								]
								sampleData.append(data)
							}

							self.workoutData[self.getWorkoutDataType(type: type)] = sampleData
							self.compleatedQuery += 1
					}
				}
				self.healthStore.execute(workoutQuery)
			}
		}
	}

	private func dumpJson() {
		getWorkoutData()

		while self.compleatedQuery < WorkoutManager.shared.typesToRead.count {
			sleep(1)
		}

		if workoutData.isEmpty {
			return
		}

		guard let jData = try? JSONSerialization.data(withJSONObject: workoutData, options: [.withoutEscapingSlashes, .fragmentsAllowed]) else {
			Logger.shared.log("Failed to convert convert object to JSON")
			return
		}

		do {
			let documentsDirectory = try FileManager.default.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: true)
			let fileName = UUID().uuidString + ".rupl"
			let fileURL = documentsDirectory.appendingPathComponent(fileName)

			try jData.write(to: fileURL)

			//						print(fileURL)
		} catch {
			Logger.shared.log("Error: \(error)")
			return
		}
	}

	private func forDispatch() -> [String] {
		var filesToSend: [String] = []

		do {
			let documentsDirectory = try FileManager.default.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: true)
			let fileURLs = try FileManager.default.contentsOfDirectory(at: documentsDirectory, includingPropertiesForKeys: nil, options: [])

			for fileURL in fileURLs {
				if fileURL.lastPathComponent.hasSuffix(".rupl") {
					filesToSend.append(fileURL.lastPathComponent)
				}
			}
		} catch {
			Logger.shared.log("Error: \(error)")
		}

		return filesToSend
	}

	private func sendData(jsonURL: URL, completion: @escaping (Bool) -> Void) {
		print("sendData", jsonURL)

		var jsonData: Data = Data()

		do {
			jsonData = try Data(contentsOf: jsonURL)
		} catch {
			Logger.shared.log("Error: \(error)")
			completion(false)
			return
		}

		guard let apiUrl = URL(string: "\(AppSettings.shared.sessionURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)") else {
			Logger.shared.log("Invalid URL")
			completion(false)
			return
		}

		var request = URLRequest(url: apiUrl)
		request.httpMethod = "POST"
		request.addValue("application/json", forHTTPHeaderField: "Content-Type")
		request.httpBody = jsonData

		let task = URLSession.shared.dataTask(with: request) { data, response, error in
			if let error = error {
				Logger.shared.log("Error: \(error)")
				completion(false)
				return
			}

			guard let httpResponse = response as? HTTPURLResponse else {
				completion(false)
				return
			}

			// TMP
			guard let data = data else {
				print("No data received")
				return
			}

			// Convert data to a string (for demonstration purposes)
			if let responseString = String(data: data, encoding: .utf8) {
				print("Response: \(responseString)")
			} else {
				print("Unable to convert data to string")
			}
			// ---

			if httpResponse.statusCode != 200 {
				completion(false)
				return
			}

			completion(true)
		}

		task.resume()
	}

	func sendSessionData(completion: @escaping (Bool) -> Void) {
		let sessions: [String] = forDispatch()

		if sessions.isEmpty {
			completion(true)
			return
		}

		DispatchQueue.global().async {
			do {
				let documentsDirectory = try FileManager.default.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: true)

				for s in sessions {
					let fileURL = documentsDirectory.appendingPathComponent(s)

					self.sendData(jsonURL: fileURL) { success in
						if success {
							do {
								try FileManager.default.removeItem(at: fileURL)
							} catch {
								Logger.shared.log("Error: \(error)")
								completion(false)
								return
							}
						} else {
							completion(false)
							return
						}
					}
				}
				completion(true)
			} catch {
				Logger.shared.log("Error: \(error)")
				completion(false)
				return
			}
		}
	}

	func sendSessionDataController(retryCount: Int) {
		var success = false
		var retry = retryCount == 0 ? 1 : retryCount

		DispatchQueue.global().async {
			while !success && retry > 0 {
				self.sendSessionData { s in
					success = s
				}

				if !success {
					Logger.shared.log("Retrying in 10 seconds...")
					retry -= 1
					sleep(10)
				} else if retryCount == 0 {
					Logger.shared.log("Failed to send JSON data after retries.")
				}
			}
		}

		// sendSessionData() { success in
		// 	if !success && retryCount > 0 {
		// 		Logger.shared.log("Retrying in 10 seconds...")
		// 		DispatchQueue.main.asyncAfter(deadline: .now() + 10) {
		// 			self.sendSessionDataController(retryCount: retryCount - 1)
		// 		}
		// 	} else if !success {
		// 		Logger.shared.log("Failed to send JSON data after retries.")
		// 	}
		// }
	}

	func postWorkout() {
#if DEBUG
		print("postWorkout()")
#endif

		dumpJson()
		sendSessionDataController(retryCount: 60) // 10 minutes
	}

	static let shared = WorkoutData()
}
