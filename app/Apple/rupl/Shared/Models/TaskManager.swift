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
		var distance: Int
		var duration: Int
		var speed: Speed
		var pulse: Pulse
	}

	struct Speed: Codable, Hashable {
		var min: Float32
		var max: Float32
	}

	struct Pulse: Codable, Hashable {
		var min: Int
		var max: Int
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

	enum ParametersToControl {
		case all, heartRate, speed
	}

	private var interval: Interval?
	private var intervalID: Int = -1

	var isTaskStarted: Bool = false
	var task: Task?
	var intervalTimeLeft: Int = 0
	var intervalDistanceLeft: Double = 0
	var intervalEndDistance: Double = 0
	var intervalHeartRateZone: (maxHeartRate: Int, minHeartRate: Int) = (0, AppSettings.shared.criticalHeartRate) // bpm
	var intervalSpeedZone: (maxSpeed: Double, minSpeed: Double) = (0, 11) // mps
//	var controlParaneters: ParametersToControl = .all
	var controlParaneters: ParametersToControl = .heartRate

	static let shared = TaskManager()

	// MARK: - API Methods
	//
	func getTask(completion: @escaping (Bool) -> Void) {
		if !AppSettings.shared.connectedToRupl {
			completion(false)
			return
		}

		DispatchQueue.global().async {
			let apiUrl = URL(string: "\(AppSettings.shared.taskURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)")!
			var request = URLRequest(url: apiUrl)

			request.httpMethod = "GET"

			let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
				if let error = error {
					Logger.shared.log("Error: \(error)")
					return
				}

				guard let httpResponse = response as? HTTPURLResponse else {
					completion(false)
					return
				}

				if httpResponse.statusCode != 200  {
					completion(false)
					return
				}

				if let data = data {
					do {
						self.task = try JSONDecoder().decode(Task.self, from: data)
						completion(true)
#if DEBUG
						//						self.printTask(self.task!)
#endif
					} catch {
						Logger.shared.log("Error parsing JSON: \(error)")
						completion(false)
					}
				}

			}
			task.resume()
		}
	}

	func declineTask(completion: @escaping (Bool) -> Void) {
		if !AppSettings.shared.connectedToRupl {
			self.task = nil
			completion(true)
			return
		}

		DispatchQueue.global().async {
			let apiUrl = URL(string: "\(AppSettings.shared.taskURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)&task_id=\(self.task?.id ?? 0)")!
			var request = URLRequest(url: apiUrl)

			request.httpMethod = "POST"

			let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
				if let error = error {
					Logger.shared.log("Error: \(error)")
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

				self.task = nil
				self.intervalHeartRateZone = (0, AppSettings.shared.criticalHeartRate) // bpm
				self.intervalSpeedZone = (0, 11) // mps
				self.controlParaneters = .heartRate

				completion(true)
			}
			task.resume()
		}
	}

	// MARK: - Running session
	//
	func runSession() {
		if task == nil {
			return
		}

		if isTaskStarted && intervalTimeLeft == 0 && intervalDistanceLeft == 0 {
			intervalID += 1
			let interval: Interval? = intervalID < (task?.intervals.count ?? 0) ? task?.intervals[intervalID] : nil

			if interval == nil {
#if DEBUG
				print("### SESSION COMPLEATED ###")
#endif
				task = nil
				return
			}

			if interval?.speed.max != 0 && interval?.pulse.max == 0 {
				controlParaneters = .speed
			} else if interval?.speed.max == 0 && interval?.pulse.max != 0 {
				controlParaneters = .heartRate
			} else {
				controlParaneters = .all
			}

			intervalTimeLeft = interval?.duration ?? 0
			intervalDistanceLeft = Double(interval?.distance ?? 0)
			intervalHeartRateZone = (interval?.pulse.max ?? 0, interval?.pulse.min ?? 0)
			intervalSpeedZone = (Double(interval?.speed.max ?? 0), Double(interval?.speed.min ?? 0))
#if DEBUG
			//			printInterval()
#endif
		}
	}

	// MARK: - Get interval parameters
	//
	func getHeartRateInterval(pz: String) -> (maxHeartRate: Int, minHeartRate: Int) {
		switch pz {
			case HeartRateZones.pz1.rawValue:
				return (AppSettings.shared.pz1NotInZone, 60)
			case HeartRateZones.pz2.rawValue:
				return (AppSettings.shared.pz2Easy, AppSettings.shared.pz1NotInZone)
			case HeartRateZones.pz3.rawValue:
				return (AppSettings.shared.pz3FatBurning, AppSettings.shared.pz2Easy)
			case HeartRateZones.pz4.rawValue:
				return (AppSettings.shared.pz4Aerobic, AppSettings.shared.pz3FatBurning)
			case HeartRateZones.pz5.rawValue:
				return (AppSettings.shared.pz5Anaerobic, AppSettings.shared.pz4Aerobic)
			default:
				return (AppSettings.shared.criticalHeartRate, 60)
		}
	}

	private func getSpeedInterval(speed: Double) -> (maxSpeed: Double, minSpeed: Double) {
		return (speed * 0.8, speed * 1.2)
	}

	// MARK: - Printers
	//
	private func printTask(_ task: Task) {
		print("ID: \(task.id), Description: \(task.description)")

		for i in task.intervals {
			print("### INTERVAL \(i.id) - \(i.description) ###")
			print("# Speed:\t\(i.speed)")
			print("# Pulse:\t\(i.pulse)")
			print("# Distance:\t\(i.distance)")
			print("# Duration:\t\(i.duration)")
		}
	}

	private func printInterval() {
		print("### INTERVAL \(intervalID + 1) ###")
		print("# Speed:\t\(intervalSpeedZone)")
		print("# Pulse:\t\(intervalHeartRateZone)")
		print("# Distance:\t\(intervalDistanceLeft)")
		print("# Duration:\t\(intervalTimeLeft)")
	}
}
