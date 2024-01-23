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
	@State private var isSegmentSheetActive = false

	var body: some View {
		TimelineView(MetricsTimelineSchedule(from: workoutManager.session?.startDate ?? Date(),
											 isPaused: workoutManager.sessionState == .paused)) { context in
			VStack(alignment: .leading) {
				ElapsedTimeView(elapsedTime: elapsedTime(with: context.date), showSubseconds: context.cadence == .live)
					.foregroundColor(.ruplYellow)
//				Text(workoutManager.speed.description + " /km")
//					.foregroundStyle(.ruplBlue)
				Text(workoutManager.convertToMinutesPerKilometer(metersPerSecond: workoutManager.last10SpeedAverage) + " /km")
					.foregroundStyle(.ruplBlue)
				Text((workoutManager.distance / 1000).formatted(.number.precision(.fractionLength(2))) + " km")
				Text((workoutManager.heartRate).formatted(.number.precision(.fractionLength(0))) + " bpm")
					.foregroundStyle(.ruplRed)
//				Text(workoutManager.power.formatted(.number.precision(.fractionLength(0))) + " w")
//				Text(workoutManager.cadence.formatted(.number.precision(.fractionLength(0))) + " rpm")
//				Text(workoutManager.water.formatted(.number.precision(.fractionLength(0))) + " oz")
			}
			.font(.system(.title, design: .rounded).monospacedDigit().lowercaseSmallCaps())
			.frame(maxWidth: .infinity, alignment: .leading)
			.ignoresSafeArea(edges: .bottom)
			.scenePadding()
			.padding([.top], 30)
			.onChange(of: workoutManager.lastSegment) {
				isSegmentSheetActive = true
			}
			.sheet(isPresented: $isSegmentSheetActive) {
//				print("segment changed")
			} content: {
				LastSegmentView()
			}
		}

	}

	func elapsedTime(with contextDate: Date) -> TimeInterval {
		return workoutManager.builder?.elapsedTime(at: contextDate) ?? 0
	}
}
