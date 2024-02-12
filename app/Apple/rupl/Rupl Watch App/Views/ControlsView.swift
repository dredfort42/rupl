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

	var body: some View {
		VStack {
			if (workoutManager.session?.state != .ended) {

				if !workoutManager.sessionState.isActive {
					Button {
						startWorkout()
					} label: {
						ZStack {
							Circle()
								.frame(width: 140, height: 140)
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
					.frame(width: 140, height: 140)

					Spacer()

					HStack {
						Button {
							isSettingsActive = true
							print("Show settings screen")
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
						.padding(.leading)
						Spacer()
					}
					.padding(.top, -10)
				} else {
					HStack {
						Button {
							workoutManager.isPauseSetWithButton = workoutManager.sessionState == .running
							workoutManager.sessionState = workoutManager.sessionState == .running ? .paused : .running
							workoutManager.session?.state == .running ? workoutManager.session?.pause() : workoutManager.session?.resume()
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
						.padding(.leading)

						Spacer()
					}

					Spacer()

					HStack {

						Spacer()
						
						Button {
							workoutManager.session?.stopActivity(with: .now)
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
						.padding(.leading)
					}
					.padding(.top, -20)
				}
			} else {
				LoadingIndicatorView()
			}
		}
		.sheet(isPresented: $isSettingsActive) {
			print("Close settings screen")
		} content: {
			SettingsView()
		}
	}

	private func startWorkout() {
		workoutManager.resetWorkout()
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
