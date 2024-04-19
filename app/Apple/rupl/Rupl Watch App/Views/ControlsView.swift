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
import AVFoundation

struct ControlsView: View {
	@EnvironmentObject var workoutManager: WorkoutManager

	@State private var isSettingsActive = false
	@State private var isTaskActive = false
	@State private var isNewTask: Bool = false
	@State private var taskRequestsTimer: Timer?

	var body: some View {
		VStack {
			if (workoutManager.sessionState != .ended && workoutManager.sessionState != .stopped) {
				if workoutManager.sessionState == .notStarted {
					StartView()
				} else {
					StopView()
				}
			} else {
				LoadingIndicatorView()
			}
		}
		.sheet(isPresented: $isSettingsActive) {
//			print("Close settings screen")
		} content: {
			SettingsView()
		}
		.sheet(isPresented: $isTaskActive) {
			isNewTask = false
		} content: {
			TaskView()
		}
		.onAppear() {
			var taskRequestsTimer = Timer.scheduledTimer(withTimeInterval: 5, repeats: true) { _ in
				if TaskManager.shared.task == nil {
					TaskManager.shared.getTask { result in
						isNewTask = result
					}
				}
			}
		}
		.onDisappear() {
			taskRequestsTimer?.invalidate()
		}
	}

	@ViewBuilder
	private func StartView() -> some View {
		HStack {
			// Start
			GetButtonView(size: 130, color: .ruplBlue, image: "", title: "Start") {
				isNewTask = false

				if TaskManager.shared.task != nil {
					TaskManager.shared.isTaskStarted = true
				} else {
					TaskManager.shared.intervalHeartRateZone = TaskManager.shared.getHeartRateInterval(pz: AppSettings.shared.runningTaskHeartRate)
				}

				workoutManager.startWorkout()
			}
			Spacer()
		}
		.padding(.horizontal)

		Spacer()

		HStack {
			// Settings
			GetButtonView(size: 30, color: .ruplGray, image: "gear", title: "") {
				isSettingsActive = true
			}
				.padding(.top, 10)
			Spacer()
			// Task
			if isNewTask {
				GetButtonView(size: 40, color: .ruplRed, image: "figure.stand", title: "") {
					isTaskActive = true
				}
			} else {
				GetButtonView(size: 40, color: .ruplGreen, image: "figure.run", title: "") {
					isTaskActive = true
				}
			}
		}
		.padding(.horizontal)
		.padding(.top, -20)
	}

	@ViewBuilder
	private func StopView() -> some View {
		HStack {
			// Pause
			let title = workoutManager.sessionState == .running ? "Pause" : "Resume"
			let image = workoutManager.sessionState == .running ? "pause" : "play"
			GetButtonView(size: 110, color: .ruplYellow, image: image, title: title) {
				workoutManager.isPauseSetWithButton = workoutManager.sessionState == .running
				workoutManager.sessionState = workoutManager.sessionState == .running ? .paused : .running
				workoutManager.session?.state == .running ? workoutManager.session?.pause() : workoutManager.session?.resume()
				workoutManager.session?.state == .running ? SoundEffects.shared.playStopSound() : SoundEffects.shared.playStartSound()
			}
			Spacer()
		}
		.padding(.horizontal)

		Spacer()

		HStack {
			// Settings
			GetButtonView(size: 30, color: .ruplGray, image: "gear", title: "") {
				isSettingsActive = true
			}
				.padding(.top, 30)
			Spacer()
			// End
			GetButtonView(size: 60, color: .ruplRed, image: "", title: "Stop") {
				TaskManager.shared.task = nil

				workoutManager.finishWorkout()
			}
		}
		.padding(.horizontal)
		.padding(.top, -20)
	}

	@ViewBuilder
	private func GetButtonView(size: CGFloat, color: Color, image: String, title: String, action: @escaping () -> Void) -> some View {
		Button(action: action) {
			let imageSize = size > 50 ? size / 4 : size * 2 / 3
			ZStack {
				Circle()
					.frame(width: size, height: size)
					.foregroundColor(color)
					.opacity(0.01)
				VStack {
					if image != "" {
						Image(systemName: image)
							.resizable()
							.scaledToFit()
							.frame(height: imageSize)
					}
					if title != "" {
						Text(title)
							.padding(.horizontal)
							.fontWeight(.medium)
					}
				}
				.foregroundColor(color)
			}
		}
		.clipShape(Circle())
		.overlay {
			Circle().stroke(color.opacity(0.8), lineWidth: 2)
		}
		.buttonStyle(.bordered)
		.frame(width: size, height: size)
	}

}
