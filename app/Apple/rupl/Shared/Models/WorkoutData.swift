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
					"h_acc": location.horizontalAccuracy,
					"altitude": location.altitude,
					"ellipsoidal_altitude": location.ellipsoidalAltitude,
					"v_acc": location.verticalAccuracy,
					"course": location.course,
					"crs_acc": location.courseAccuracy,
					"speed": location.speed,
					"spd_acc": location.speedAccuracy
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

								self.workoutData[type.description] = route
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

							self.workoutData[type.description] = sampleData
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
			let temporaryDirectoryURL = FileManager.default.temporaryDirectory
			let fileName = UUID().uuidString + ".rupl"
			let fileURL = temporaryDirectoryURL.appendingPathComponent(fileName)

			try jData.write(to: fileURL)

			//			print(fileURL)
		} catch {
			Logger.shared.log("Error: \(error)")
			return
		}
	}

	private func forDispatch() -> [String] {
		var filesToSend: [String] = []
		let temporaryDirectoryURL = FileManager.default.temporaryDirectory

		do {
			let fileURLs = try FileManager.default.contentsOfDirectory(at: temporaryDirectoryURL, includingPropertiesForKeys: nil, options: [])

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

		//		print(String(data: jsonData, encoding: .utf8)!)

		guard let apiUrl = URL(string: "\(AppSettings.shared.sessionURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)") else {
			Logger.shared.log("Invalid URL")
			completion(false)
			return
		}

		var request = URLRequest(url: apiUrl)
		request.httpMethod = "POST"
		request.httpBody = jsonData
		//		request.addValue("application/json", forHTTPHeaderField: "Content-Type")

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
			let temporaryDirectoryURL = FileManager.default.temporaryDirectory

			for s in sessions {
				let fileURL = temporaryDirectoryURL.appendingPathComponent(s)

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
		}
	}

	func sendSessionDataController(retryCount: Int) {
		sendSessionData() { success in
			if !success && retryCount > 0 {
				Logger.shared.log("Retrying in 10 seconds...")
				DispatchQueue.main.asyncAfter(deadline: .now() + 10) {
					self.sendSessionDataController(retryCount: retryCount - 1)
				}
			} else if !success {
				Logger.shared.log("Failed to send JSON data after retries.")
			}
		}
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
