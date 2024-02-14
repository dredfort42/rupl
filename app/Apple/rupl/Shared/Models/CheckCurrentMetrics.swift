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
			if self.heartRate != 0 && self.session?.state == .running {
				self.segmentHeartRatesSum += UInt64(self.heartRate + 0.5)
				self.segmentHeartRatesCount += 1
				self.summaryHeartRateSum += UInt64(self.heartRate + 0.5)
				self.summaryHeartRateCount += 1
			}

			if self.session?.state == .running && (-self.workoutStartTime.timeIntervalSinceNow > 300){
				if !self.checkCriticalHeartRate() {
					self.checkHeartRateZone()
				}
			}
		}
	}

	private func checkCriticalHeartRate() -> Bool {
		if Int(self.heartRate) > AppSettings.shared.criticalHeartRate {
#if targetEnvironment(simulator)
			print("* Alarm sound")
#else
			self.sounds.alarmSound?.play()
#if os(watchOS)
			Vibration.vibrate(type: .underwaterDepthCriticalPrompt)
#endif
#endif
			return true
		}

		return false
	}

	private func checkHeartRateZone() {
		if Int(self.heartRate) > TaskManager.shared.intervalHeartRateZone.maxHeartRate {
#if targetEnvironment(simulator)
			print("* Run slower sound")
#else
			self.sounds.runSlower?.play()
#if os(watchOS)
			Vibration.vibrate(type: .directionDown)
#endif
#endif
		} else if Int(self.heartRate) < TaskManager.shared.intervalHeartRateZone.minHeartRate {
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

// MARK: - Check last segment
//
extension WorkoutManager {
	func checkLastSegment() {
		DispatchQueue.global().async {
			let segment = Int(self.distance / 1000)
			if self.segmentNumber != segment {
				self.segmentFinishTime = Date()
				self.segmentNumber = segment
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
			if self.speed < AppSettings.shared.paceForAutoPause {
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

