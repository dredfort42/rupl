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
			start()
		} else {
			Logger.shared.log("Motion activity tracking is not available on this device")
		}
	}

	func start() {
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

	func stop() {
		motionActivityManager.stopActivityUpdates()
	}
}
