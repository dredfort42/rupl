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
	@AppStorage(AppSettings.isSoundNotificationOnKey) var isSoundNotificationOn = AppSettings.shared.isSoundNotificationOn

	@EnvironmentObject var workoutManager: WorkoutManager

	@State private var isSegmentSheetActive = false

	var body: some View {
		TimelineView(MetricsTimelineSchedule(from: workoutManager.session?.startDate ?? Date(),
											 isPaused: workoutManager.sessionState == .paused)) { context in
			HStack {
				Button {
					isSoundNotificationOn = !isSoundNotificationOn
				} label: {
					HStack {
						Image(systemName: isSoundNotificationOn ? "speaker.wave.2" : "speaker.slash")
						Text(isSoundNotificationOn ? "On" : "Off")
							.padding(.horizontal, 2)
							.fontWeight(.light)
							.font(.footnote)
					}
					.foregroundColor(isSoundNotificationOn ? .ruplGreen : .ruplRed)
				}
				.controlSize(.mini)
				.frame(width: 70)

				Spacer()
			}
			.padding(.horizontal)
			.padding(.bottom, -20)

			Spacer()

			VStack(alignment: .leading) {
				ElapsedTimeView(elapsedTime: elapsedTime(with: context.date), showSubseconds: context.cadence == .live)
					.foregroundColor(.ruplYellow)
				Text(workoutManager.convertToMinutesPerKilometer(metersPerSecond: workoutManager.last10SpeedAverage) + " /km")
					.foregroundStyle(.ruplBlue)
				Text((workoutManager.distance / 1000).formatted(.number.precision(.fractionLength(2))) + " km")
				Text((workoutManager.heartRate).formatted(.number.precision(.fractionLength(0))) + " bpm")
					.foregroundStyle(.ruplRed)
			}
			.font(.system(.title, design: .rounded).monospacedDigit().lowercaseSmallCaps())
			.frame(maxWidth: .infinity, alignment: .leading)
			.ignoresSafeArea(edges: .bottom)
			.scenePadding()
			.padding([.top], 30)
			.onChange(of: workoutManager.segmentNumber) {
				if workoutManager.segmentNumber > 0 {
					isSegmentSheetActive = true
				}
			}
			.sheet(isPresented: $isSegmentSheetActive) {
//				print("last segment view closed")
			} content: {
				SegmentView()
			}

			Spacer()
		}

	}

	func elapsedTime(with contextDate: Date) -> TimeInterval {
		return workoutManager.builder?.elapsedTime(at: contextDate) ?? 0
	}
}
