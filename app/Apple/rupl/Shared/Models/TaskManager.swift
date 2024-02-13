//
//  TaskManager.swift
//  rupl
//
//  Created by Dmitry Novikov on 13/02/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation

class TaskManager {
	struct PaceZone {
		//	seconds per segment
		var maxPace: Int = 0
		var minPace: Int = 0
	}
	
	struct HeartRateZone {
		//	bpm
		var maxHeartRate: Int = 0
		var minHeartRate: Int = 0
	}

	var intervalHeartRateZone = HeartRateZone()
	
	static let shared = TaskManager()

	init() {
		intervalHeartRateZone.maxHeartRate = AppSettings.shared.pz3FatBurning
		intervalHeartRateZone.minHeartRate = AppSettings.shared.pz2Easy
	}

}
