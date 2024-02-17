//
//  WorkoutManager+watchOS.swift
//  rupl
//
//  Created by Dmitry Novikov on 04/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os
import HealthKit

// MARK: - Workout session management
//
extension WorkoutManager {

	func startWorkout(workoutConfiguration: HKWorkoutConfiguration) async throws {
		session = try HKWorkoutSession(healthStore: healthStore, configuration: workoutConfiguration)
		builder = session?.associatedWorkoutBuilder()
		session?.delegate = self
		builder?.delegate = self
		builder?.dataSource = HKLiveWorkoutDataSource(healthStore: healthStore, workoutConfiguration: workoutConfiguration)
		routeBuilder = HKWorkoutRouteBuilder(healthStore: healthStore, device: nil)

		session?.startActivity(with: .now)
		try await builder?.beginCollection(at: session?.startDate ?? Date())
	}
}

//	MARK: - HKLiveWorkoutBuilderDelegate
//	HealthKit calls the delegate methods on an anonymous serial background queue.
//	The methods need to be nonisolated explicitly.
//
extension WorkoutManager: HKLiveWorkoutBuilderDelegate {

	nonisolated func workoutBuilder(_ workoutBuilder: HKLiveWorkoutBuilder, didCollectDataOf collectedTypes: Set<HKSampleType>) {
		Task { @MainActor in
			var allStatistics: [HKStatistics] = []

			for type in collectedTypes {
				if let quantityType = type as? HKQuantityType, let statistics = workoutBuilder.statistics(for: quantityType) {
					updateForStatistics(statistics)
					allStatistics.append(statistics)
				}
			}
		}
	}

	nonisolated func workoutBuilderDidCollectEvent(_ workoutBuilder: HKLiveWorkoutBuilder) {
	}
}
