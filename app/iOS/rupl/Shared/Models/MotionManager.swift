//
//  MotionManager.swift
//  rupl
//
//  Created by Dmitry Novikov on 26/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os
import CoreMotion

class MotionManager {
	let motionActivityManager = CMMotionActivityManager()
	var autoPauseState: Bool = true

	init() {
		if CMMotionActivityManager.isActivityAvailable() {
			startMotionUpdates()
		} else {
			Logger.shared.log("Motion activity tracking is not available on this device")
		}
	}

	func startMotionUpdates() {
		motionActivityManager.startActivityUpdates(to: OperationQueue.main) { (activity: CMMotionActivity?) in
			if let activity = activity {
				if activity.running {
					self.autoPauseState = false
				} else {
					self.autoPauseState = true
				}
			}
		}
	}
}
