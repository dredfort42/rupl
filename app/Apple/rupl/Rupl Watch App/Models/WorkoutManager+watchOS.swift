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
		if !isTimerStarted {
			isTimerStarted = true
			timer.schedule(deadline: .now(), repeating: .seconds(1))
			timer.setEventHandler {self.timerActions()}
			timer.resume()
		}

		locationManager.locationManager.startUpdatingLocation()
		session = try HKWorkoutSession(healthStore: healthStore, configuration: workoutConfiguration)
		builder = session?.associatedWorkoutBuilder()
		session?.delegate = self
		builder?.delegate = self
		builder?.dataSource = HKLiveWorkoutDataSource(healthStore: healthStore, workoutConfiguration: workoutConfiguration)
		routeBuilder = HKWorkoutRouteBuilder(healthStore: healthStore, device: nil)

		let startDate = Date()
		session?.startActivity(with: startDate)
		sessionState = .paused
		session?.pause()
		try await builder?.beginCollection(at: startDate)
	}

	func finishWorkout() async {
		let finishedWorkout: HKWorkout?

		do {
			try await builder?.endCollection(at: stopTime)
			finishedWorkout = try await builder?.finishWorkout()
			session?.end()
		} catch {
			Logger.shared.log("Failed to end workout: \(error))")
			return
		}

		if (finishedWorkout != nil) {
			do {
				try await routeBuilder?.finishRoute(with: finishedWorkout!, metadata: nil)
			} catch {
				Logger.shared.log("Failed to associate the route with the workout: \(error)")
				return
			}
		}
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
