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
		ScrollView {
			summaryListView()
				.scenePadding()
		}
		.navigationTitle("Summary")
		.navigationBarTitleDisplayMode(.inline)
	}

	@ViewBuilder
	private func summaryListView() -> some View {
		let totalWorkoutTime: TimeInterval = workoutManager.stopTime.timeIntervalSince(workoutManager.session?.startDate ?? Date())
		let runDuration: TimeInterval = workoutManager.builder?.elapsedTime(at: Date()) ?? 0
		let distance: Double = workoutManager.distance / 1000
		let averageSpeedMetersPerSecond: Double = workoutManager.averageSpeedMetersPerSecond
		let averageHeartRate: Int = workoutManager.averageHeartRate

		VStack(alignment: .leading) {
			SummaryMetricView(title: "Total workout time",
							  value: workoutManager.formatDuration(seconds: totalWorkoutTime)
			)

			SummaryMetricView(title: "Run duration",
							  value: workoutManager.formatDuration(seconds: runDuration)
			).foregroundStyle(.ruplYellow)

			SummaryMetricView(title: "Distance",
							  value: distance.formatted(.number.precision(.fractionLength(2)))
			)

			SummaryMetricView(title: "Average pace",
							  value: workoutManager.convertToMinutesPerKilometer(metersPerSecond: averageSpeedMetersPerSecond)
			).foregroundColor(.ruplBlue)
	
			SummaryMetricView(title: "Average heart rate",
							  value: averageHeartRate.formatted(.number.precision(.fractionLength(0)))
			).foregroundStyle(.ruplRed)

			Button {
				dismiss()
			} label: {
				Text("Done")
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
//
//#Preview {
//	SummaryView()
//}
