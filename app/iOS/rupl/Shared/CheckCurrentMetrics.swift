//
//  CheckCurrentMetrics.swift
//  rupl
//
//  Created by Dmitry Novikov on 23/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation

// MARK: - Check heart rate
//
extension WorkoutManager {
	func checkHeartRate() {
		DispatchQueue.global().async {
			if self.heartRate != 0 {
				self.lastSegmentHeartRatesSum += self.heartRate
				self.lastSegmentHeartRatesCount += 1
			}

			let pulse: Int = Int(self.heartRate)

			//	Checking the critical heart rate level
			if pulse > self.parameters.pz5Anaerobic {
				self.heartRateNotificationTimer = 0
				self.sounds.alarmSound?.play()
#if os(watchOS)
				Vibration.vibrate(type: .underwaterDepthCriticalPrompt)
#endif
			}

			// 	MARK: - TMP Checking the puls zone
			if self.sessionState.isActive {
				if pulse > self.parameters.pz3FatBurning {
					self.heartRateNotificationTimer = 10
					self.sounds.runSlower?.play()
#if os(watchOS)
					Vibration.vibrate(type: .directionDown)
#endif
				} else if pulse < self.parameters.pz2Easy && (-(self.session?.startDate?.timeIntervalSinceNow ?? 0) > 600) {
					self.heartRateNotificationTimer = 10
					self.sounds.runFaster?.play()
#if os(watchOS)
					Vibration.vibrate(type: .directionUp)
#endif
				}
			}
		}
	}
}

// MARK: - Check last segment
//
extension WorkoutManager {
	func checkLastSegment() {
		DispatchQueue.global().async {
			let segment = Int(self.distance / 1000)
			if self.lastSegment != segment {
				self.lastSegmentStopTime = Date()
				self.lastSegmentViewPresentTime = self.parameters.timeForShowLastSegmentView
				self.lastSegment = segment
				self.sounds.segmentSound?.play()
#if os(watchOS)
				Vibration.vibrate(type: .success)
#endif
			}
		}
	}
}

// MARK: - Check speed
//
extension WorkoutManager {
	func checkSpeed() {
		DispatchQueue.global().async {
			//	Calculate average speed from last 10 measurements
			if self.speed < self.parameters.paceForAutoPause {
				return
			}

			while self.last10SpeedMeasurements.count > 10 {
				self.last10SpeedMeasurementsSum -= self.last10SpeedMeasurements[0]
				self.last10SpeedMeasurements.remove(at: 0)
			}

			self.last10SpeedMeasurementsSum += self.speed
			self.last10SpeedMeasurements.append(self.speed)

			self.last10SpeedAverage = self.last10SpeedMeasurementsSum / Double(self.last10SpeedMeasurements.count)
		}
	}
}

