//
//  SummaryView.swift
//  rupl
//
//  Created by Dmitry Novikov on 22/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import HealthKit
import SwiftUI

struct SummaryView: View {
	@Binding var workout: HKWorkout?
	

	var body: some View {
//		let totalWorkoutTime: TimeInterval = workoutManager.stopTime.timeIntervalSince(workoutManager.session?.startDate ?? Date())
//		let runDuration: TimeInterval = workout.duration
//		let distance: Double = workout.statistics(for: HKQuantityType(.distanceWalkingRunning))?.sumQuantity()?.doubleValue(for: HKUnit.meterUnit(with: .kilo)) ?? 0
//		let speedMetersPerSecond: Double = distance * 1000 / runDuration
//		let averagePower: Double = workout.statistics(for: HKQuantityType(.runningPower))?.averageQuantity()?.doubleValue(for: HKUnit.watt()) ?? 0
//		let calories: Double = workout.statistics(for: HKQuantityType(.activeEnergyBurned))?.sumQuantity()?.doubleValue(for: HKUnit.largeCalorie()) ?? 0
//		let averageHeartRate: Double = workout.statistics(for: HKQuantityType(.heartRate))?.averageQuantity()?.doubleValue(for: HKUnit.count().unitDivided(by: .minute())) ?? 0

		if let workout = workout {
			LazyVGrid(columns: [GridItem(.adaptive(minimum: 150))]) {
				GridItemView(title: "Total Time", value: workout.duration.description)
					.foregroundStyle(.yellow)

//				GridItemView(title: "Total Distance", value: workout.totalDistance)
//					.foregroundStyle(.orange)

//				GridItemView(title: "Total Energy", value: workout.totalEnergy)
//					.foregroundStyle(.pink)
//
//				GridItemView(title: "Average Speed", value: workout.averageCyclingSpeed)
//					.foregroundStyle(.green)
//
//				GridItemView(title: "Average Power", value: workout.averageCyclingPower)
//					.foregroundStyle(.pink)
//
//				GridItemView(title: "Average Cadence", value: workout.averageCyclingCadence)
//					.foregroundStyle(.black)
			}
		}
	}
}

private struct GridItemView: View {
	var title: String
	var value: String

	var body: some View {
		VStack {
			Text(title)
				.foregroundStyle(.foreground)
			Text(value)
				.font(.system(.title2, design: .rounded).lowercaseSmallCaps())
		}
	}
}
