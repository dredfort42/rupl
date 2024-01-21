//
//  WorkoutParameters.swift
//  rupl
//
//  Created by Dmitry Novikov on 08/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation

//	MARK: - A structure for declare workout constants
//
struct WorkoutParameters {
	let maxMirroringErrors: UInt8 = 100 // count
	let minHorizontalAccuracy:  Double = 30 // m
	let paceForAutoPause: Double = 1.85 // m/s
	let paceForAutoResume: Double = 2.22 // m/s
//	let paceForAutoPause: Double = 3.5 // m/s
//	let paceForAutoResume: Double = 3.9 // m/s
	let timeForShowLastSegmentView = 15 // s
}
