//
//  TaskManager.swift
//  rupl
//
//  Created by Dmitry Novikov on 13/02/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os

class TaskManager: ObservableObject {
	struct Interval: Codable, Hashable {
		var id: Int
		var description: String
		var speed: Double
		var pulse_zone: Int
		var distance: Int
		var duration: Int
	}

	struct Task: Codable {
		var id: Int
		var description: String
		var intervals: [Interval]
		var compleated: Bool
	}

	enum HeartRateZones: String, CaseIterable, Identifiable {
		case any, pz1, pz2, pz3, pz4, pz5
		var id: Self { self }
	}


	private var interval: Interval?
	private var intervalID: Int = -1

	@Published var task: Task?
	@Published var isRunTaskDownloaded: Bool = false
	@Published var isNewRunTaskAvailable: Bool = false
	@Published var isRunTaskAccepted: Bool?
	@Published var isRunTaskStarted: Bool = false
	@Published var intervalTimeLeft: Int = 0
	@Published var intervalDistanceLeft: Double = 0
	@Published var intervalEndDistance: Double = 0
	@Published var intervalHeartRateZone: (maxHeartRate: Int, minHeartRate: Int) = (0, AppSettings.shared.criticalHeartRate) // bpm
	@Published var intervalSpeedZone: (maxSpeed: Double, minSpeed: Double) = (0, 11) // mps

	static let shared = TaskManager()

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
						self.task = try JSONDecoder().decode(Task.self, from: data)
						self.isNewRunTaskAvailable = true
						self.isRunTaskDownloaded = true
						completion("getTask() successfully completed")
#if DEBUG
						self.printTask(self.task!)
#endif
					} catch {
						Logger.shared.log("Error parsing JSON: \(error)")
						completion("getTask() completed with error")
					}
				}
			}
			task.resume()
		}
	}

	func runSession() {
		if isRunTaskAccepted == false {
			return
		}

		if intervalID == -1 {
			intervalID = 0
			interval = getInterval(intervalID)
			intervalTimeLeft = interval?.duration ?? 0
			intervalDistanceLeft = Double(interval?.distance ?? 0)
			intervalHeartRateZone = getHeartRateInterval(pz: HeartRateZones.allCases[(interval?.pulse_zone ?? 0)].rawValue)
			intervalSpeedZone = getSpeedInterval(speed: interval?.speed ?? 0)
#if DEBUG
			printInterval()
#endif
			return
		}

		if !(task?.compleated ?? false) && isRunTaskStarted && intervalTimeLeft == 0 && intervalDistanceLeft == 0 {
			intervalID += 1
			interval = getInterval(intervalID)

			if interval == nil {
				task?.compleated = true
#if DEBUG
				print("### SESSION COMPLEATED ###")
#endif
				return
			}

			intervalTimeLeft = interval?.duration ?? 0
			intervalDistanceLeft = Double(interval?.distance ?? 0)
			intervalHeartRateZone = getHeartRateInterval(pz: HeartRateZones.allCases[(interval?.pulse_zone ?? 0)].rawValue)
			intervalSpeedZone = getSpeedInterval(speed: interval?.speed ?? 0)
#if DEBUG
			printInterval()
#endif
		}
	}

	private func getInterval(_ num: Int) -> Interval? {
		if num >= 0 && num < task?.intervals.count ?? 0 {
			return (task?.intervals[num])
		}

		return (nil)
	}

	// MARK: - Get measurements intervals
	//
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

	private func getSpeedInterval(speed: Double) -> (maxSpeed: Double, minSpeed: Double) {
		return (speed * 0.9, speed * 1.2)
	}

	// MARK: - Printers
	//
	private func printTask(_ task: Task) {
		print("ID: \(task.id), Description: \(task.description)")

		for i in task.intervals {
			print("### INTERVAL \(i.id) - \(i.description) ###")
			print("# Speed:\t\(i.speed)")
			print("# Pulse:\t\(i.pulse_zone)")
			print("# Distance:\t\(i.distance)")
			print("# Duration:\t\(i.duration)")
		}
	}

	private func printInterval() {
		print("### INTERVAL \(intervalID) ###")
		print("# Speed:\t\(intervalSpeedZone)")
		print("# Pulse:\t\(intervalHeartRateZone)")
		print("# Distance:\t\(intervalDistanceLeft)")
		print("# Duration:\t\(intervalTimeLeft)")
	}
}
