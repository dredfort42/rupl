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
	private let motionManager = CMMotionManager()
	private let acceletationsCount: UInt8 = 50
	private var lastAccelerations: [Double] = []
	private var lastAccelerationsSum: Double = 0
	private var accelerationAverage: Double = 0

	var autoPauseState: Bool = true

	static let shared = MotionManager()

	func start() {
		startMotionUpdates()
	}

	func stop() {
		motionManager.stopDeviceMotionUpdates()
	}

	private func startMotionUpdates() {
		motionManager.startDeviceMotionUpdates(to: OperationQueue.main) { (motionData, error) in
			guard let motionData = motionData else {
				return
			}

			let acceleration = motionData.userAcceleration
			let accelerationMagnitude = sqrt(pow(acceleration.x, 2) + pow(acceleration.y, 2) + pow(acceleration.z, 2))

			while self.lastAccelerations.count > self.acceletationsCount {
				self.lastAccelerationsSum -= self.lastAccelerations[0]
				self.lastAccelerations.remove(at: 0)
			}
			
			self.lastAccelerationsSum += accelerationMagnitude
			self.lastAccelerations.append(accelerationMagnitude)
			
			self.accelerationAverage = self.lastAccelerationsSum / Double(self.acceletationsCount)

			if self.accelerationAverage < AppSettings.shared.accelerationForAutoPause {
				self.autoPauseState = true
			} else if self.accelerationAverage > AppSettings.shared.accelerationForAutoResume {
				self.autoPauseState = false
			}

			print(self.autoPauseState)
		}
	}
}
