//
//  AppDeligate.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 04/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import os
import WatchKit
import HealthKit
import SwiftUI

class AppDelegate: NSObject, WKApplicationDelegate {

	func handle(_ workoutConfiguration: HKWorkoutConfiguration) {
		Task {
			do {
				WorkoutManager.shared.resetWorkout()
				try await WorkoutManager.shared.startWorkout(workoutConfiguration: workoutConfiguration)
				Logger.shared.log("Successfully started workout")
			} catch {
				Logger.shared.log("Failed started workout")
			}
		}
	}
}
