//
//  TaskManager.swift
//  rupl
//
//  Created by Dmitry Novikov on 13/02/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os

class TaskManager {
	struct Interval: Codable {
		var id: Int
		var description: String
		var speed: Int
		var pulse_zone: Int
		var distance: Int
		var duration: Int
	}

	struct Task: Codable {
		var id: Int
		var description: String
		var intervals: [Interval]
	}

	enum HeartRateZones: String, CaseIterable, Identifiable {
		case any, pz1, pz2, pz3, pz4, pz5
		var id: Self { self }
	}

	lazy var intervalHeartRateZone: (maxHeartRate: Int, minHeartRate: Int) = getHeartRateInterval(pz: AppSettings.shared.runningTaskHeartRate) {
		didSet {
			print(self.intervalHeartRateZone)
		}
	}

	static let shared = TaskManager()

	func getHeartRateInterval(pz: String) -> (maxHeartRate: Int, minHeartRate: Int) {
		switch pz {
			case HeartRateZones.pz1.rawValue:
				return (AppSettings.shared.pz1NotInZone, 0)
			case HeartRateZones.pz2.rawValue:
				return (AppSettings.shared.pz2Easy, AppSettings.shared.pz1NotInZone)
			case HeartRateZones.pz3.rawValue:
				return (AppSettings.shared.pz3FatBurning, AppSettings.shared.pz2Easy)
			case HeartRateZones.pz4.rawValue:
				return (AppSettings.shared.pz4Aerobic, AppSettings.shared.pz3FatBurning)
			case HeartRateZones.pz5.rawValue:
				return (AppSettings.shared.pz5Anaerobic, AppSettings.shared.pz4Aerobic)
			default:
				return (AppSettings.shared.criticalHeartRate, 0)
		}
	}

	func getTask(completion: @escaping (String) -> Void) {
		let apiUrl = URL(string: "\(AppSettings.shared.taskURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)")!
		var request = URLRequest(url: apiUrl)

		request.httpMethod = "GET"

		let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
			if let error = error {
				Logger.shared.log("Error: \(error)")
				return
			}

			if let data = data {
				do {
					let decoder = JSONDecoder()
					let runTask = try decoder.decode(Task.self, from: data)
					print("ID: \(runTask.id), Description: \(runTask.description)")
					for i in runTask.intervals {
						print(i.id)
						print(i.description)
						print(i.speed)
						print(i.pulse_zone)
						print(i.distance)
						print(i.duration)
					}
				} catch {
					Logger.shared.log("Error parsing JSON: \(error)")
				}
			}
			completion("OK")
		}
		task.resume()
	}

}
