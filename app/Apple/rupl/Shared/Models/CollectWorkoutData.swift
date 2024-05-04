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
		guard let routes = results as? [HKWorkoutRoute] else {
			print("No workout routes found or error occurred")
			return []
		}

		var routeData: [[String: Any]] = []

		let routeQuery = HKWorkoutRouteQuery(route: routes[0]) { (query, locationsOrNil, done, errorOrNil) in
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
				routeData.append(locationData)
			}
		}

		self.healthStore.execute(routeQuery)
		return routeData
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
