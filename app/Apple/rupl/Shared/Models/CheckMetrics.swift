//
//  CheckMetrics.swift
//  rupl
//
//  Created by Dmitry Novikov on 23/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os

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

			if self.session?.state == .running && (-(self.session?.startDate ?? Date()).timeIntervalSinceNow > 300){
				if !self.checkCriticalHeartRate() {
					self.checkHeartRateZone()
				}
			}
		}
	}

	private func checkCriticalHeartRate() -> Bool {
		if Int(self.heartRate) > AppSettings.shared.criticalHeartRate {
			SoundEffects.shared.playAlarmSound()
			return true
		}

		return false
	}

	private func checkHeartRateZone() {
		if Int(self.heartRate) > TaskManager.shared.intervalHeartRateZone.maxHeartRate {
			SoundEffects.shared.playRunSlowerSound()
		} else if Int(self.heartRate) < TaskManager.shared.intervalHeartRateZone.minHeartRate {
			SoundEffects.shared.playRunFasterSound()
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
				SoundEffects.shared.playSegmentSound()
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

//	MARK: - Auto pause logic
//
extension WorkoutManager {
	func autoPause() {
		if AppSettings.shared.useAutoPause && !isPauseSetWithButton && (sessionState == .running || sessionState == .paused) {
#if targetEnvironment(simulator)
			MotionManager.shared.autoPauseState = true
#endif

			guard let LMAutoPauseState = LocationManager.shared.autoPauseState else {
				setPauseState(isPaused: MotionManager.shared.autoPauseState)
				return
			}

			if LMAutoPauseState && MotionManager.shared.autoPauseState {
				setPauseState(isPaused: true)
			} else if !LMAutoPauseState && !MotionManager.shared.autoPauseState {
				setPauseState(isPaused: false)
			}
		}
	}

	func setPauseState(isPaused: Bool) {
		if self.sessionState == .running {
			if isPaused {
				self.session?.pause()
				SoundEffects.shared.playStopSound()
			}
		} else {
			if !isPaused {
				self.session?.resume()
				SoundEffects.shared.playStartSound()
			}
		}
	}
}

//	MARK: - Add locations to the route
//
extension WorkoutManager {
	func addLocationsToRoute() {
		if !LocationManager.shared.filteredLocations.isEmpty {
			self.routeBuilder?.insertRouteData(LocationManager.shared.filteredLocations) { (success, error) in
				if !success {
					Logger.shared.log("Failed to add locations to the route: \(error))")
				}
			}
		}
	}
}
