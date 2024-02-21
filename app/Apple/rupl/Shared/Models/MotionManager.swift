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

	var isAvailable: Bool = false

	static let shared = MotionManager()

	private func requestAuthorization() {
		if !isAvailable {
#if DEBUG
			print("MotionManager.requestAuthorization()")
#endif
			if CMMotionActivityManager.isActivityAvailable() {
				isAvailable = true
			} else {
				Logger.shared.log("Motion activity tracking is not available on this device")
			}
		}
	}

	func start() {
		requestAuthorization()
		startActivityUpdates()
	}

	func stop() {
		motionActivityManager.stopActivityUpdates()
	}

	func startActivityUpdates() {
		motionActivityManager.startActivityUpdates(to: OperationQueue.main) { (activity: CMMotionActivity?) in
			if let activity = activity {
#if DEBUG
				print("motionActivityManager.activity: ", activity)
#endif
				if activity.running {
					self.autoPauseState = false
				} else {
					self.autoPauseState = true
				}
			}
		}
	}
}
