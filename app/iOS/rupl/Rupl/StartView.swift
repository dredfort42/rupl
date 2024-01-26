//
//  StartView.swift
//  rupl
//
//  Created by Dmitry Novikov on 22/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//


import os
import SwiftUI
import HealthKitUI
import HealthKit

struct StartView: View {
	@EnvironmentObject var workoutManager: WorkoutManager
	@State private var isFullScreenCoverActive = false
	@State private var didStartWorkout = false
	@State private var triggerAuthorization = false

	var body: some View {
		NavigationStack {
			VStack {
				Button {
					if !workoutManager.sessionState.isActive {
						//					   startCyclingOnWatch()
					}
					didStartWorkout = true
				} label: {
					let title = workoutManager.sessionState.isActive ? "View ongoing running" : "Start running on watch"
					ButtonLabel(title: title, systemImage: "figure.run")
						.frame(width: 120, height: 120)
						.fontWeight(.medium)
				}
				.clipShape(Circle())
				.overlay {
					Circle().stroke(.white, lineWidth: 4)
				}
				.shadow(radius: 7)
				.buttonStyle(.bordered)
				.tint(.ruplBlue)
				.foregroundColor(.black)
				.frame(width: 400, height: 400)
			}
			.onAppear() {
				triggerAuthorization.toggle()
				//			   workoutManager.retrieveRemoteSession()
			}
			.healthDataAccessRequest(store: workoutManager.healthStore,
									 shareTypes: workoutManager.typesToShare,
									 readTypes: workoutManager.typesToRead,
									 trigger: triggerAuthorization, completion: { result in
				switch result {
					case .success(let success):
						Logger.shared.log("\(success) for authorization")
					case .failure(let error):
						Logger.shared.log("\(error) for authorization")
				}
			})
			.navigationDestination(isPresented: $didStartWorkout) {
				MirroringWorkoutView()
			}
			.navigationBarTitle("Mirroring Workout")
			.navigationBarTitleDisplayMode(.inline)
			.toolbar {
				ToolbarItem(placement: .automatic) {
					Button {
						isFullScreenCoverActive = true
					} label: {
						Label("Workout list", systemImage: "list.bullet")
					}
				}
			}
			.fullScreenCover(isPresented: $isFullScreenCoverActive) {
				WorkoutListView()
			}
		}
	}

	//	private func startCyclingOnWatch() {
	//		Task {
	//			do {
	//				try await workoutManager.startWatchWorkout(workoutType: .cycling)
	//			} catch {
	//				Logger.shared.log("Failed to start running on the paired watch.")
	//			}
	//		}
	//	}
}
