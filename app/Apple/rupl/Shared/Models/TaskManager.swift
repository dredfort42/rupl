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
	var intervalID: Int = -1
	var intervalTimeLeft: Int = 0 { didSet { runTask() } }
	var intervalDistanceLeft: Double = 0 { didSet { runTask() } }
	var intervalEndDistance: Double = 0
	lazy var intervalHeartRateZone: (maxHeartRate: Int, minHeartRate: Int) = getHeartRateInterval(pz: AppSettings.shared.runningTaskHeartRate)
	var intervalSpeed: Double = 0


	private struct Interval: Codable {
		var id: Int
		var description: String
		var speed: Double
		var pulse_zone: Int
		var distance: Int
		var duration: Int
	}

	private struct Task: Codable {
		var id: Int
		var description: String
		var intervals: [Interval]
	}

	enum HeartRateZones: String, CaseIterable, Identifiable {
		case any, pz1, pz2, pz3, pz4, pz5
		var id: Self { self }
	}

	private var task: Task?

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
		DispatchQueue.global().async {
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
						self.task = runTask
					} catch {
						Logger.shared.log("Error parsing JSON: \(error)")
					}
				}
#if DEBUG
				self.printTask()
#endif
				completion("getTask() compleated")
			}
			task.resume()
		}
	}

	private func printTask() {
		if task == nil {
			return
		}

		print("ID: \(task!.id), Description: \(task!.description)")

		for i in task!.intervals {
			print("### INTERVAL \(i.id) - \(i.description) ###")
			print("# Speed:\t\(i.speed)")
			print("# Pulse:\t\(i.pulse_zone)")
			print("# Distance:\t\(i.distance)")
			print("# Duration:\t\(i.duration)")
		}
	}

	func runTask() {
		if intervalID != -2 && intervalTimeLeft == 0 && intervalDistanceLeft == 0 {
			getNextInterval()
		}
	}

	private func getNextInterval() {
		if intervalID != -2 {
			intervalID += 1
		} else {
			return
		}

		intervalTimeLeft = task?.intervals[intervalID].duration ?? 0
		intervalDistanceLeft = Double(task?.intervals[intervalID].distance ?? 0)
		intervalHeartRateZone = getHeartRateInterval(pz: HeartRateZones.allCases[(task?.intervals[intervalID].pulse_zone ?? 0)].rawValue)
		intervalSpeed = task?.intervals[intervalID].speed ?? 0

#if DEBUG
		printNextInterval()
#endif

		if intervalID + 1 == task?.intervals.count {
			intervalID = -2
		}
	}

	private func printNextInterval() {
		print("### INTERVAL \(intervalID) ###")
		print("# Speed:\t\(intervalSpeed)")
		print("# Pulse:\t\(intervalHeartRateZone)")
		print("# Distance:\t\(intervalDistanceLeft)")
		print("# Duration:\t\(intervalTimeLeft)")
	}

}
