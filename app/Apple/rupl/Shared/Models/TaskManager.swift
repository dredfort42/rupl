//
//  TaskManager.swift
//  rupl
//
//  Created by Dmitry Novikov on 13/02/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation

class TaskManager {
//	struct PaceZone {
//		//	seconds per segment
//		var maxPace: Int = 0
//		var minPace: Int = 0
//	}

//	struct HeartRateZone {
//		//	bpm
//		var maxHeartRate: Int = 0
//		var minHeartRate: Int = 0
//	}

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

}
