//
//  LastSegmentView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 21/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import HealthKit
import SwiftUI

struct LastSegmentView: View {
	@EnvironmentObject var workoutManager: WorkoutManager
	@Environment(\.dismiss) var dismiss

	@State var lastSegmentStartTime: Date = Date()
	@State var lastSegmentStopTime: Date = Date()
	@State var lastSegmentHeartRatesSum: Double = 0
	@State var lastSegmentHeartRatesCount: Int = 0


	var body: some View {
			ScrollView {
				summaryListView()
					.scenePadding()
			}
//			.navigationTitle("Segment")
//			.navigationBarTitleDisplayMode(.inline)
			.onAppear() {
				lastSegmentStartTime = workoutManager.lastSegmentStartTime
				lastSegmentStopTime = workoutManager.lastSegmentStopTime
				lastSegmentHeartRatesSum = workoutManager.lastSegmentHeartRatesSum
				lastSegmentHeartRatesCount = workoutManager.lastSegmentHeartRatesCount

				workoutManager.lastSegmentStartTime = workoutManager.lastSegmentStopTime
				workoutManager.lastSegmentHeartRatesSum = 0
				workoutManager.lastSegmentHeartRatesCount = 0
			}
			.onChange(of: workoutManager.lastSegmentViewPresentTime) {_, time in
				if time == 0 {
					dismiss()
				}
			}
	}

	@ViewBuilder
	private func summaryListView() -> some View {
		let pace: TimeInterval = lastSegmentStopTime.timeIntervalSince(lastSegmentStartTime)
		let averageHeartRate: Double = lastSegmentHeartRatesSum / Double(lastSegmentHeartRatesCount)

		VStack(alignment: .leading) {

			Text("Distance: \(workoutManager.lastSegment) km")
				.font(.system(.title3, design: .rounded).monospacedDigit().lowercaseSmallCaps())
				.padding([.top], 15)

			Divider()

			Text(workoutManager.formatDuration(seconds: pace) + " /km")
				.font(.system(.title, design: .rounded).monospacedDigit().lowercaseSmallCaps())
				.foregroundColor(.ruplBlue)

			Text(averageHeartRate.formatted(.number.precision(.fractionLength(0))) + " bpm")
				.font(.system(.title, design: .rounded).monospacedDigit().lowercaseSmallCaps())
				.foregroundColor(.ruplRed)

		}
	}
}
