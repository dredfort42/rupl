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
		locationManager.locationManager.startUpdatingLocation()

		if !isTimerStarted {
			isTimerStarted = true
			timer.schedule(deadline: .now(), repeating: .seconds(1))
			timer.setEventHandler {self.timerActions()}
			timer.resume()
		}

		session = try HKWorkoutSession(healthStore: healthStore, configuration: workoutConfiguration)
		builder = session?.associatedWorkoutBuilder()
		session?.delegate = self
		builder?.delegate = self
		builder?.dataSource = HKLiveWorkoutDataSource(healthStore: healthStore, workoutConfiguration: workoutConfiguration)

//	--------------------
//		if isMirroring {
//			//	Start mirroring the session to the companion device
//			try await session?.startMirroringToCompanionDevice()
//		}
//	--------------------

		//	Start the workout session activity
		let startDate = Date()
		session?.startActivity(with: startDate)
//		sessionState = .paused
//		session?.pause()
		try await builder?.beginCollection(at: startDate)
	}

//	--------------------
//	func handleReceivedData(_ data: Data) throws {
//		guard let decodedQuantity = try NSKeyedUnarchiver.unarchivedObject(ofClass: HKQuantity.self, from: data) else {
//			return
//		}
//
//		let sampleDate = Date()
//		Task {
//			let waterSample = [HKQuantitySample(type: HKQuantityType(.dietaryWater), quantity: decodedQuantity, start: sampleDate, end: sampleDate)]
//			try await builder?.addSamples(waterSample)
//		}
//	}
//	--------------------

}

//	MARK: - HKLiveWorkoutBuilderDelegate
//	HealthKit calls the delegate methods on an anonymous serial background queue.
//	The methods need to be nonisolated explicitly.
//
extension WorkoutManager: HKLiveWorkoutBuilderDelegate {

	nonisolated func workoutBuilder(_ workoutBuilder: HKLiveWorkoutBuilder, didCollectDataOf collectedTypes: Set<HKSampleType>) {

		//	HealthKit calls this method on an anonymous serial background queue.
		//	Use Task to provide an asynchronous context so MainActor can come to play.
		Task { @MainActor in
			var allStatistics: [HKStatistics] = []

			for type in collectedTypes {
				if let quantityType = type as? HKQuantityType, let statistics = workoutBuilder.statistics(for: quantityType) {
					updateForStatistics(statistics)
					allStatistics.append(statistics)
				}
			}

//	--------------------
//			if isMirroring {
//				let archivedData = try? NSKeyedArchiver.archivedData(withRootObject: allStatistics, requiringSecureCoding: true)
//				guard let archivedData = archivedData, !archivedData.isEmpty else {
//					Logger.shared.log("Encoded running data is empty")
//					return
//				}
//
//				//	Send a Data object to the connected remote workout session.
//				await sendData(archivedData)
//			}
//	--------------------
			
		}
	}

	nonisolated func workoutBuilderDidCollectEvent(_ workoutBuilder: HKLiveWorkoutBuilder) {
	}
}
