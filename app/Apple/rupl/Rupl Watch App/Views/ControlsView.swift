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
	@State private var isNewTaskAvailable = TaskManager.shared.isNewRunTaskAvailable

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
//			print("Close settings screen")
		} content: {
			TaskView()
		}
		.onAppear() {
			TaskManager.shared.getTask { result in
				isNewTaskAvailable = TaskManager.shared.isNewRunTaskAvailable
			}
		}
	}

	@ViewBuilder
	private func StartView() -> some View {
		HStack {
			Button {
				workoutManager.startWorkout()
				TaskManager.shared.isRunTaskStarted = true
			} label: {
				ZStack {
					Circle()
						.frame(width: 130, height: 130)
						.foregroundColor(.ruplBlue)
						.opacity(0.1)
					VStack {
						Image(systemName: "figure.run")
						Text("Start")
							.padding(.horizontal)
							.fontWeight(.medium)
					}
					.foregroundColor(.ruplBlue)
				}
			}
			.clipShape(Circle())
			.overlay {
				Circle().stroke(.ruplBlue.opacity(0.8), lineWidth: 2)
			}
			.buttonStyle(.bordered)
			.frame(width: 130, height: 130)

			Spacer()
		}
		.padding(.horizontal)

		Spacer()

		HStack {
			SettingsButtonView()
				.padding(.top, 10)
			Spacer()
			TaskButtonView()
		}
		.padding(.horizontal)
		.padding(.top, -20)
	}

	@ViewBuilder
	private func StopView() -> some View {
		HStack {
			Button {
				workoutManager.isPauseSetWithButton = workoutManager.sessionState == .running
				workoutManager.sessionState = workoutManager.sessionState == .running ? .paused : .running
				workoutManager.session?.state == .running ? workoutManager.session?.pause() : workoutManager.session?.resume()
				workoutManager.session?.state == .running ? SoundEffects.shared.playStopSound() : SoundEffects.shared.playStartSound()
			} label: {
				let title = workoutManager.sessionState == .running ? "Pause" : "Resume"
				let systemImage = workoutManager.sessionState == .running ? "pause" : "play"

				ZStack {
					Circle()
						.frame(width: 110, height: 110)
						.foregroundColor(.ruplYellow)
						.opacity(0.1)
					VStack {
						Image(systemName: systemImage)
						Text(title)
							.padding(.horizontal)
							.fontWeight(.medium)
					}
					.foregroundColor(.ruplYellow)
				}
			}
			.clipShape(Circle())
			.overlay {
				Circle().stroke(.ruplYellow.opacity(0.8), lineWidth: 2)
			}
			.buttonStyle(.bordered)
			.frame(width: 110, height: 110)

			Spacer()
		}
		.padding(.horizontal)

		Spacer()

		HStack {
			SettingsButtonView()
				.padding(.top, 30)
			Spacer()

			Button {
				workoutManager.finishWorkout()
			} label: {
				ZStack {
					Circle()
						.frame(width: 60, height: 60)
						.foregroundColor(.ruplRed)
						.opacity(0.1)
					VStack {
						Image(systemName: "xmark")
						Text("End")
							.padding(.horizontal)
							.fontWeight(.medium)
					}
					.foregroundColor(.ruplRed)
				}
			}
			.clipShape(Circle())
			.overlay {
				Circle().stroke(.ruplRed.opacity(0.8), lineWidth: 2)
			}
			.buttonStyle(.bordered)
			.frame(width: 60, height: 60)
		}
		.padding(.horizontal)
		.padding(.top, -20)
	}

	@ViewBuilder
	private func SettingsButtonView() -> some View {
		Button {
			isSettingsActive = true
		} label: {
			ZStack {
				Circle()
					.frame(width: 30, height: 30)
					.foregroundColor(.ruplGray)
					.opacity(0.1)
				Image(systemName: "gear")
					.foregroundColor(.ruplGray)
			}
		}
		.clipShape(Circle())
		.overlay {
			Circle().stroke(.ruplGray.opacity(0.8), lineWidth: 2)
		}
		.buttonStyle(.bordered)
		.frame(width: 30, height: 30)
	}

	@ViewBuilder
	private func TaskButtonView() -> some View {
		var color: Color = isNewTaskAvailable ? .ruplYellow : .ruplGreen

		Button {
			isTaskActive = true
		} label: {

			ZStack {
				Circle()
					.frame(width: 40, height: 40)
					.foregroundColor(color)
					.opacity(0.1)
				Image(systemName: "list.clipboard")
					.foregroundColor(color)
			}
		}
		.clipShape(Circle())
		.overlay {
			Circle().stroke(color.opacity(0.8), lineWidth: 2)
		}
		.buttonStyle(.bordered)
		.frame(width: 40, height: 40)
	}
}
