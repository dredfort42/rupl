//
//  TimerActions.swift
//  rupl
//
//  Created by Dmitry Novikov on 18/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os

// MARK: - Timer actions
//
extension WorkoutManager {
	func timerActions() {
		autoPause()
		addLocationsToRoute()
		heartRateTimer()
		lastSegmentViewPresentTimer()
	}
}

//	MARK: - Auto pause logic
//
extension WorkoutManager {
	func autoPause() {
		if self.sessionState.isActive && !self.isPauseSetWithButton {
			if self.sessionState == .running {
				if self.locationManager.autoPauseState || 
					(self.locationManager.accuracy > (self.parameters.minHorizontalAccuracy * 0.8) &&
					 self.motionManager.autoPauseState) {
					self.sessionState = .paused
					self.session?.pause()
					self.sounds.stopSound?.play()
#if os(watchOS)
					Vibration.vibrate(type: .notification)
#endif
				}
			} else {
				if !self.locationManager.autoPauseState || 
					(self.locationManager.accuracy > (self.parameters.minHorizontalAccuracy * 0.8) &&
					 !self.motionManager.autoPauseState) {
					self.sessionState = .running
					self.session?.resume()
					self.sounds.startSound?.play()
#if os(watchOS)
					Vibration.vibrate(type: .notification)
#endif
				}
			}
		}
	}
}

//	MARK: - Add locations to the route
//
extension WorkoutManager {
	func addLocationsToRoute() {
		DispatchQueue.global().async {
			guard !self.locationManager.filteredLocations.isEmpty else { return }
			
			self.routeBuilder.insertRouteData(self.locationManager.filteredLocations) { (success, error) in
				if !success {
					Logger.shared.log("Failed to add locations to the route: \(error))")
				}
			}
		}
	}
}

//	MARK: - Check heart rate zone
//
extension WorkoutManager {
	func heartRateTimer() {
		DispatchQueue.main.async {
			if self.heartRateNotificationTimer != 0 {
				self.heartRateNotificationTimer -= 1
			}
		}
	}
}

//	MARK: - LastSegmentView present timer
//
extension WorkoutManager {
	func lastSegmentViewPresentTimer() {
		DispatchQueue.main.async {
			if self.lastSegmentViewPresentTime != 0 {
				self.lastSegmentViewPresentTime -= 1
			}
		}
	}
}
