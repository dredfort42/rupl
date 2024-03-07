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

	@State var timer: Timer?

	var body: some View {
		ScrollView {
			summaryListView()
				.scenePadding()
		}
		.navigationTitle("Summary")
		.navigationBarTitleDisplayMode(.inline)
		.onAppear() {
			timer = Timer.scheduledTimer(withTimeInterval: TimeInterval(AppSettings.shared.viewNotificationTimeOut * 3), repeats: false) { _ in
				dismiss()
			}
		}
	}

	@ViewBuilder
	private func summaryListView() -> some View {
		let totalWorkoutTime: TimeInterval = (workoutManager.session?.endDate ?? Date()).timeIntervalSince(workoutManager.session?.startDate ?? Date())
		let runDuration: TimeInterval = workoutManager.builder?.elapsedTime(at: workoutManager.session?.endDate ?? Date()) ?? 0
		let distance: Double = workoutManager.distance
		let averageSpeedMetersPerSecond: Double = runDuration > 0 ? distance / runDuration : 0
		let averageHeartRate: Int = workoutManager.summaryHeartRateCount > 0 ? Int(Double(workoutManager.summaryHeartRateSum / UInt64(workoutManager.summaryHeartRateCount)) + 0.5) : 0

		VStack(alignment: .leading) {
			SummaryMetricView(title: "Total workout time",
							  value: workoutManager.formatDuration(seconds: totalWorkoutTime)
			)

			SummaryMetricView(title: "Run duration",
							  value: workoutManager.formatDuration(seconds: runDuration)
			).foregroundStyle(.ruplYellow)

			SummaryMetricView(title: "Distance",
							  value: (distance / 1000).formatted(.number.precision(.fractionLength(2)))
			)

			SummaryMetricView(title: "Average pace",
							  value: workoutManager.convertToMinutesPerKilometer(metersPerSecond: averageSpeedMetersPerSecond)
			).foregroundColor(.ruplBlue)
			
			SummaryMetricView(title: "Average heart rate",
							  value: String(averageHeartRate)
			).foregroundStyle(.ruplRed)

			Button {
				dismiss()
			} label: {
				Text("Save")
			}
			.padding(.vertical)
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

//#Preview {
//	SummaryView()
//}
