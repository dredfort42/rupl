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
	@State var lastSegmentHeartRatesSum: UInt64 = 0
	@State var lastSegmentHeartRatesCount: UInt = 0

	
	var body: some View {
		ScrollView {
			LastSegmentSummaryView()
				.scenePadding()
		}
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
	private func LastSegmentSummaryView() -> some View {
		let pace: TimeInterval = lastSegmentStopTime.timeIntervalSince(lastSegmentStartTime)
		let averageHeartRate: Double = Double(lastSegmentHeartRatesSum / UInt64(lastSegmentHeartRatesCount))

		VStack(alignment: .leading) {
			HStack(alignment: .lastTextBaseline) {
				Spacer()
				Text("\(workoutManager.lastSegment)")
					.font(.system(size: 40, weight: .regular, design: .rounded).monospacedDigit().lowercaseSmallCaps())
				Text(" km")
					.font(.system(size: 25, weight: .medium, design: .rounded).monospacedDigit().lowercaseSmallCaps())
			}
			.padding([.top], -18)
			.padding([.bottom], -2)
			
			Divider()
			
			HStack(alignment: .lastTextBaseline) {
				Spacer()
				Text(workoutManager.formatDuration(seconds: pace))
					.font(.system(.title, design: .rounded).monospacedDigit().lowercaseSmallCaps())
				
				Text(" /km")
					.font(.system(size: 20, weight: .medium, design: .rounded).monospacedDigit().lowercaseSmallCaps())
			}
			.foregroundColor(.ruplBlue)
			.padding([.top], -2)
			
			HStack(alignment: .lastTextBaseline) {
				Spacer()
				Text(averageHeartRate.formatted(.number.precision(.fractionLength(0))))
					.font(.system(.title, design: .rounded).monospacedDigit().lowercaseSmallCaps())
				Text(" bpm")
					.font(.system(size: 20, weight: .medium, design: .rounded).monospacedDigit().lowercaseSmallCaps())
			}
			.foregroundColor(.ruplRed)
		}
	}
}
