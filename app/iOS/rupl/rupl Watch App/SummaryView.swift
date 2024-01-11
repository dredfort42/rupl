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
		let totalWorkoutTime: TimeInterval = workoutManager.stopTime.timeIntervalSince(workoutManager.session?.startDate ?? Date())
		let runDuration: TimeInterval = workout.duration
		let distance: Double = workout.statistics(for: HKQuantityType(.distanceWalkingRunning))?.sumQuantity()?.doubleValue(for: HKUnit.meterUnit(with: .kilo)) ?? 0
		let speedMetersPerSecond: Double = distance * 1000 / runDuration
		let averagePower: Double = workout.statistics(for: HKQuantityType(.runningPower))?.averageQuantity()?.doubleValue(for: HKUnit.watt()) ?? 0
		let calories: Double = workout.statistics(for: HKQuantityType(.activeEnergyBurned))?.sumQuantity()?.doubleValue(for: HKUnit.largeCalorie()) ?? 0
		let averageHeartRate: Double = workout.statistics(for: HKQuantityType(.heartRate))?.averageQuantity()?.doubleValue(for: HKUnit.count().unitDivided(by: .minute())) ?? 0

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
							  value: workoutManager.convertToMinutesPerKilometer(speedMetersPerSecond: speedMetersPerSecond)
			).foregroundColor(.ruplBlue)

			SummaryMetricView(title: "Average power",
							  value: averagePower.formatted(.number.precision(.fractionLength(0)))
			)

			SummaryMetricView(title: "Calories",
							  value: calories.formatted(.number.precision(.fractionLength(0)))
			)

			SummaryMetricView(title: "Average heart rate",
							  value: averageHeartRate.formatted(.number.precision(.fractionLength(0)))
			).foregroundStyle(.ruplRed)

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
