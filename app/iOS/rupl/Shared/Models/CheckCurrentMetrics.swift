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
			if pulse > self.pz5Anaerobic {
#if targetEnvironment(simulator)
				print("* Alarm sound")
#else
				self.sounds.alarmSound?.play()
#if os(watchOS)
				Vibration.vibrate(type: .underwaterDepthCriticalPrompt)
#endif
#endif
			}

			// 	MARK: - TMP Checking the puls zone
			if self.session?.state != .paused && self.heartRateNotificationTimer == 0 {
				if pulse > self.pz2Easy {
					self.heartRateNotificationTimer = 10
#if targetEnvironment(simulator)
					print("* Run slower sound")
#else
					self.sounds.runSlower?.play()
#if os(watchOS)
					Vibration.vibrate(type: .directionDown)
#endif
#endif
				} else if pulse < self.pz1NotInZone && (-(self.session?.startDate?.timeIntervalSinceNow ?? 0) > 600) {
					self.heartRateNotificationTimer = 10
#if targetEnvironment(simulator)
					print("* Run faster sound")
#else
					self.sounds.runFaster?.play()
#if os(watchOS)
					Vibration.vibrate(type: .directionUp)
#endif
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
				self.lastSegmentViewPresentTime = self.timeForShowLastSegmentView
				self.lastSegment = segment
#if targetEnvironment(simulator)
					print("* Segment sound")
#else
				self.sounds.segmentSound?.play()
#if os(watchOS)
				Vibration.vibrate(type: .success)
#endif
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
			if self.speed < self.paceForAutoPause {
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

