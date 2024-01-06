//
//  ControlsView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 04/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import os
import SwiftUI
import HealthKit

struct ControlsView: View {
	@EnvironmentObject var workoutManager: WorkoutManager

	var body: some View {
		VStack {
			Button {
				startWorkout()
			} label: {
				ButtonLabel(title: "Start", systemImage: "figure.run")
			}
			.disabled(workoutManager.sessionState.isActive)
			.tint(.ruplBlue)

			Button {
				workoutManager.sessionState == .running ? workoutManager.session?.pause() : workoutManager.session?.resume()
			} label: {
				let title = workoutManager.sessionState == .running ? "Pause" : "Resume"
				let systemImage = workoutManager.sessionState == .running ? "pause" : "play"
				ButtonLabel(title: title, systemImage: systemImage)
			}
			.disabled(!workoutManager.sessionState.isActive)
			.tint(.ruplYellow)

			Button {
				workoutManager.session?.stopActivity(with: .now)
			} label: {
				ButtonLabel(title: "End", systemImage: "xmark")
			}
			.tint(.ruplRed)
			.disabled(!workoutManager.sessionState.isActive)
		}
	}

	private func startWorkout() {
		Task {
			do {
				let configuration = HKWorkoutConfiguration()
				configuration.activityType = .running
				configuration.locationType = .outdoor
				try await workoutManager.startWorkout(workoutConfiguration: configuration)
			} catch {
				Logger.shared.log("Failed to start workout \(error))")
			}
		}
	}
}
