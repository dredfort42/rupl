//
//  TimerActions.swift
//  rupl
//
//  Created by Dmitry Novikov on 18/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os

//	MARK: - Auto pause logic
//
extension WorkoutManager {
	func autoPause() {
		if AppSettings.shared.useAutoPause && self.sessionState == .running && !self.isPauseSetWithButton {
			var isPaused: Bool = false

#if targetEnvironment(simulator)
			MotionManager.shared.autoPauseState = false
#endif

			if	LocationManager.shared.autoPauseState || MotionManager.shared.autoPauseState {
				isPaused = true
			}

			if self.sessionState == .running {
				if isPaused {
					self.session?.pause()
#if targetEnvironment(simulator)
					print("* Stop sound")
#else
					self.sounds.stopSound?.play()
#if os(watchOS)
					Vibration.vibrate(type: .notification)
#endif
#endif
				}
			} else {
				if !isPaused {
					self.session?.resume()
#if targetEnvironment(simulator)
					print("* Start sound")
#else
					self.sounds.startSound?.play()
#if os(watchOS)
					Vibration.vibrate(type: .notification)
#endif
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
		DispatchQueue.main.async {
			if !LocationManager.shared.filteredLocations.isEmpty {
				self.routeBuilder?.insertRouteData(LocationManager.shared.filteredLocations) { (success, error) in
					if !success {
						Logger.shared.log("Failed to add locations to the route: \(error))")
					}
				}
			}
		}
	}
}
