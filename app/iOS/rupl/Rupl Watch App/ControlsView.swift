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
	@State private var isAnimating = false

	var body: some View {
		VStack {
			if (!workoutManager.isSessionEnded) {

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
						Circle().stroke(.ruplBlue, lineWidth: 1)
					}
					.buttonStyle(.bordered)
					.frame(width: 140, height: 140)

					Spacer()

					HStack {
						Button {
							print("show settings screen ->")
						} label: {
							ZStack {
								Circle()
									.frame(width: 30, height: 30)
									.foregroundColor(.gray)
									.opacity(0.1)
								Image(systemName: "gear")
									.foregroundColor(.gray)
							}
						}
						.clipShape(Circle())
						.overlay {
							Circle().stroke(.gray, lineWidth: 1)
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
									.frame(width: 120, height: 120)
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
							Circle().stroke(.ruplYellow, lineWidth: 1)
						}
						.buttonStyle(.bordered)
						.frame(width: 120, height: 120)
						.padding(.leading)

						Spacer()
					}


					HStack {
						Spacer()

						Button {
							workoutManager.isSessionEnded = true
							workoutManager.session?.stopActivity(with: .now)
						} label: {
							ZStack {
								Circle()
									.frame(width: 80, height: 80)
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
							Circle().stroke(.ruplRed, lineWidth: 1)
						}
						.buttonStyle(.bordered)
						.frame(width: 80, height: 80)
						.padding(.leading)
					}
					.padding(.top, -20)
				}
			} else {
				LoadingIndicatorView()
			}
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
