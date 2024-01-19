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
	}
}

//	MARK: - Auto pause logic
//
extension WorkoutManager {
	func autoPause() {
		if sessionState.isActive && !isPauseSetWithButton {
			if sessionState == .running {
				if locationManager.autoPauseState {
					sessionState = .paused
					session?.pause()
					sounds.stopSound?.play()
#if os(watchOS)
					Vibration.vibrate()
#endif
				}
			} else {
				if !locationManager.autoPauseState {
					sessionState = .running
					session?.resume()
					sounds.startSound?.play()
#if os(watchOS)
					Vibration.vibrate()
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
		guard !locationManager.filteredLocations.isEmpty else { return }

		routeBuilder.insertRouteData(locationManager.filteredLocations) { (success, error) in
			if !success {
				Logger.shared.log("Failed to add locations to the route: \(error))")
			}
		}
	}
}
