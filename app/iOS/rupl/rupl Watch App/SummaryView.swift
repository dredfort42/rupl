//
//  SummaryView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 04/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import HealthKit
import SwiftUI

struct SummaryView: View {
	@EnvironmentObject var workoutManager: WorkoutManager
	@Environment(\.dismiss) var dismiss

	var body: some View {
		if let workout = workoutManager.workout {
			ScrollView {
				summaryListView(workout: workout)
					.scenePadding()
			}
			.navigationTitle("Summary")
			.navigationBarTitleDisplayMode(.inline)
		} else {
			ProgressView("Saving Workout")
				.navigationBarHidden(true)
		}
	}

	@ViewBuilder
	private func summaryListView(workout: HKWorkout) -> some View {
		VStack(alignment: .leading) {
			SummaryMetricView(title: "Total time",
							  value: workoutManager.formatDuration(seconds: ((workoutManager.stopTime ?? Date()).timeIntervalSince(workoutManager.startTime ?? Date()))))

			SummaryMetricView(title: "Run duration",
							  value: workoutManager.formatDuration(seconds: workout.duration))
			.foregroundStyle(.yellow)

			Group {
				Text("Activity rings")
				ActivityRingsView(healthStore: workoutManager.healthStore)
					.frame(width: 100, height: 100)
				Button {
					dismiss()
				} label: {
					Text("Done")
				}
			}
		}
	}
}

struct SummaryMetricView: View {
	var title: String
	var value: String

	var body: some View {
		Text(title)
			.foregroundStyle(.foreground)
		Text(value)
			.font(.system(.title2, design: .rounded).lowercaseSmallCaps())
		Divider()
	}
}
