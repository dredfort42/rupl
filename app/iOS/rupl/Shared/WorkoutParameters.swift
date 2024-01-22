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
	//	Mirroring
	let maxMirroringErrors: UInt8 = 100 // count

	// 	Location
	let minHorizontalAccuracy:  Double = 30 // m

	//	Auto pause
	let paceForAutoPause: Double = 1.85 // m/s
	let paceForAutoResume: Double = 2.25 // m/s
	//	Parameters for simulator
//	let paceForAutoPause: Double = 3.7 // m/s
//	let paceForAutoResume: Double = 3.9 // m/s

	//	Last segmetn metrics
	let timeForShowLastSegmentView = 20 // s

	//	Pulse zones
	let pz1NotInZone: Int = 126 // bpm
	let pz2Easy: Int = 137 // bpm
	let pz3FatBurning: Int = 147 // bpm
	let pz4Aerobic: Int = 158 // bpm
	let pz5Anaerobic: Int = 168 // bpm
}
