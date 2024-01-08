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
	let maxMirroringErrors: UInt8 = 100
	let paceForAutoPause: Double = 0.8 // m/s
	let paceForAutoResume: Double = 1.6 // m/s
//	let paceForAutoPause: Double = 3.9 // m/s
//	let paceForAutoResume: Double = 4.2 // m/s
}
