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
	let motionManager = CMMotionActivityManager()
	var autoPauseState: Bool = true

	var isAvailable: Bool = false

	static let shared = MotionManager()

	private func requestAuthorization() {
		if !isAvailable {
			print("Motion manager request authorization")
			if CMMotionActivityManager.isActivityAvailable() {
				isAvailable = true
			} else {
				Logger.shared.log("Motion activity tracking is not available on this device")
			}
		}
	}

	func start() {
		requestAuthorization()
		motionManager.startActivityUpdates(to: OperationQueue.main) { (activity: CMMotionActivity?) in
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
		motionManager.stopActivityUpdates()
	}
}
