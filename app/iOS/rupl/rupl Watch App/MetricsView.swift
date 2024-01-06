//
//  MetricsView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 04/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI
import HealthKit

struct MetricsView: View {
	@EnvironmentObject var workoutManager: WorkoutManager

	var body: some View {
		TimelineView(MetricsTimelineSchedule(from: workoutManager.session?.startDate ?? Date(),
											 isPaused: workoutManager.sessionState == .paused)) { context in
			VStack(alignment: .leading) {
				ElapsedTimeView(elapsedTime: elapsedTime(with: context.date), showSubseconds: context.cadence == .live)
					.foregroundColor(.ruplYellow)
					.font(.system(.headline, design: .rounded).monospacedDigit().lowercaseSmallCaps())
				Text(convertToMinutesAndSeconds(seconds: workoutManager.speed*60) + " m/km")
					.foregroundStyle(.ruplBlue)
				Text((workoutManager.distance / 1000).formatted(.number.precision(.fractionLength(0))) + " km")
//				Text(workoutManager.power.formatted(.number.precision(.fractionLength(0))) + " w")
				Text(workoutManager.cadence.formatted(.number.precision(.fractionLength(0))) + " rpm")
//				Text(workoutManager.water.formatted(.number.precision(.fractionLength(0))) + " oz")
			}
			.font(.system(.title, design: .rounded).monospacedDigit().lowercaseSmallCaps())
			.frame(maxWidth: .infinity, alignment: .leading)
			.ignoresSafeArea(edges: .bottom)
			.scenePadding()
			.padding([.top], 30)
		}
	}

	func elapsedTime(with contextDate: Date) -> TimeInterval {
		return workoutManager.builder?.elapsedTime(at: contextDate) ?? 0
	}

	func convertToMinutesAndSeconds(seconds: Double) -> String {
		let minutes = Int(seconds/60)
		let seconds = Int((seconds - Double(minutes)*60))

		return seconds < 10 ? "\(minutes):0\(seconds)" : "\(minutes):\(seconds)"
	}
}
