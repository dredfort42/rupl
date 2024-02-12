//
//  SegmentView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 21/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import HealthKit
import SwiftUI

struct SegmentView: View {
	@EnvironmentObject var workoutManager: WorkoutManager
	@Environment(\.dismiss) var dismiss

	@State var pace: TimeInterval = 0
	@State var lastSegmentAverageHeartRates: Int = 0
	@State var timer: Timer?

	var body: some View {
		ScrollView {
			LastSegmentSummaryView()
				.scenePadding()
		}
		.onAppear() {
			pace = workoutManager.segmentFinishTime.timeIntervalSince(workoutManager.segmentStartTime)
			lastSegmentAverageHeartRates = workoutManager.segmentHeartRatesCount > 0 ? Int(Double(workoutManager.segmentHeartRatesSum / UInt64(workoutManager.segmentHeartRatesCount)) + 0.5) : 0
			timer = Timer.scheduledTimer(withTimeInterval: TimeInterval(AppSettings.shared.viewNotificationTimeOut), repeats: false) { _ in
				dismiss()
			}
			workoutManager.segmentStartTime = workoutManager.segmentFinishTime
			workoutManager.segmentHeartRatesSum = 0
			workoutManager.segmentHeartRatesCount = 0
		}
		.onDisappear() {
			timer?.invalidate()
		}
	}
	
	@ViewBuilder
	private func LastSegmentSummaryView() -> some View {


		VStack(alignment: .leading) {
			HStack(alignment: .lastTextBaseline) {
				Spacer()
				Text("\(workoutManager.segmentNumber)")
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
				Text(String(lastSegmentAverageHeartRates))
					.font(.system(.title, design: .rounded).monospacedDigit().lowercaseSmallCaps())
				Text(" bpm")
					.font(.system(size: 20, weight: .medium, design: .rounded).monospacedDigit().lowercaseSmallCaps())
			}
			.foregroundColor(.ruplRed)
		}
	}
}
