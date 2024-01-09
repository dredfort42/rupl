//
//	 WorkoutManager.swift
//	 rupl
//
//	 Created by Dmitry Novikov on 04/01/2024.
//	 Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os
import HealthKit
import SwiftUI

@MainActor
class WorkoutManager: NSObject, ObservableObject {

	struct SessionSateChange {
		let newState: HKWorkoutSessionState
		let date: Date
	}

	//	The workout session live states that the UI observes.
	@Published var sessionState: HKWorkoutSessionState = .notStarted
	@Published var heartRate: Double = 0
	@Published var activeEnergy: Double = 0
	@Published var distance: Double = 0
	@Published var power: Double = 0
	@Published var speed: Double = 0
	@Published var strideLength: Double = 0
	@Published var verticalOscillation: Double = 0
	@Published var groundContactTime: Double = 0
	@Published var vo2Max: Double = 0
	@Published var stepCount: Double = 0
	@Published var cadence: Double = 0
	@Published var elapsedTimeInterval: TimeInterval = 0

	//	SummaryView (watchOS) changes from Saving Workout to the metric
	//	summary view when a workout changes from nil to a valid value.
	@Published var workout: HKWorkout?

	//	HealthKit data types to share
	let typesToShare: Set = [HKQuantityType.workoutType()]

	//	HealthKit data types to read
	let typesToRead: Set = [
		HKQuantityType(.heartRate), //	count/s, Discrete (Temporally Weighted)
		HKQuantityType(.activeEnergyBurned), //	kcal, Cumulative
		HKQuantityType(.distanceWalkingRunning), //	m, Cumulative
		HKQuantityType(.runningPower), //	W, Discrete (Arithmetic)
		HKQuantityType(.runningSpeed), //	m/s, Discrete (Arithmetic)
		HKQuantityType(.runningStrideLength), //	m, Discrete (Arithmetic)
		HKQuantityType(.runningVerticalOscillation), //	cm, Discrete (Arithmetic)
		HKQuantityType(.runningGroundContactTime), //	ms, Discrete (Arithmetic)
		HKQuantityType(.vo2Max), //	ml/(kg*min), Discrete (Arithmetic)
		HKQuantityType(.stepCount), //	count, Cumulative
		HKQuantityType.workoutType(),
		HKObjectType.activitySummaryType()
	]

	let parameters = WorkoutParameters()
	let healthStore = HKHealthStore()
	var session: HKWorkoutSession?

	var isMirroring: Bool = false
	var mirroringErrorsCounter: UInt8 = 0

	var isPauseSetWithButton: Bool = false
	var startTime: Date = Date()
	var stopTime: Date = Date()

	//	Array for store last 10 speed measurements to colculate average speed
	var last10SpeedMeasurements: [Double] = []
	var last10SpeedMeasurementsSum: Double = 0
	var last10SpeedAverage: Double = 0

	#if os(watchOS)
		var builder: HKLiveWorkoutBuilder?
	#else
		var contextDate: Date?
	#endif

	let asynStreamTuple = AsyncStream.makeStream(of: SessionSateChange.self, bufferingPolicy: .bufferingNewest(1))

	static let shared = WorkoutManager()

	//	Kick off a task to consume the async stream.
	//	The next value in the stream can't start processing until
	//	"await consumeSessionStateChange(value)" returns and the loop enters
	//	the next iteration, which serializes the asynchronous operations
	private override init() {
		super.init()
		Task {
			for await value in asynStreamTuple.stream {
				await consumeSessionStateChange(value)
			}
		}
	}

	//	Consume the session state change from the async stream to update sessionState
	//	and finish the workout
	private func consumeSessionStateChange(_ change: SessionSateChange) async {
		sessionState = change.newState

		//	Wait for the session to transition states before ending the builder
		#if os(watchOS)

			//	Send the elapsed time to the iOS side
			let elapsedTimeInterval = session?.associatedWorkoutBuilder().elapsedTime(at: change.date) ?? 0
			let elapsedTime = WorkoutElapsedTime(timeInterval: elapsedTimeInterval, date: change.date)

			if isMirroring {
				if let elapsedTimeData = try? JSONEncoder().encode(elapsedTime) {
					await sendData(elapsedTimeData)
				}
			}

			guard change.newState == .stopped, let builder else {
				return
			}

			let finishedWorkout: HKWorkout?
			do {
				stopTime = Date()
				try await builder.endCollection(at: change.date)
				finishedWorkout = try await builder.finishWorkout()
				session?.end()
			} catch {
				Logger.shared.log("Failed to end workout: \(error))")
				return
			}
			workout = finishedWorkout
		#endif
	}
}

//	MARK: - Workout session management
//
extension WorkoutManager {
	func resetWorkout() {
		sessionState = .notStarted
		heartRate = 0
		activeEnergy = 0
		distance = 0
		power = 0
		speed = 0
		strideLength = 0
		verticalOscillation = 0
		groundContactTime = 0
		vo2Max = 0
		stepCount = 0
		cadence = 0
		elapsedTimeInterval = 0
		workout = nil
		session = nil
		mirroringErrorsCounter = 0
		isPauseSetWithButton = false
		startTime = Date()
		stopTime = startTime
		last10SpeedMeasurements = []
		last10SpeedMeasurementsSum = 0
		last10SpeedAverage = 0

		#if os(watchOS)
			builder = nil
		#endif
	}

	func sendData(_ data: Data) async {
		if mirroringErrorsCounter < parameters.maxMirroringErrors {
			do {
				try await session?.sendToRemoteWorkoutSession(data: data)
				mirroringErrorsCounter = 0
			} catch {
				Logger.shared.log("[\(self.mirroringErrorsCounter)] Failed to send data: \(error)")
				mirroringErrorsCounter += 1
			}
		}
	}
}

//	MARK: - Workout measurements conversions
//
extension WorkoutManager {

	//	Calculate average speed from last 10 measurements
	func calculateLast10SpeedAverage(lastSpeedMeasurement: Double) {
		if lastSpeedMeasurement < parameters.paceForAutoPause {
			return
		}

		if last10SpeedMeasurements.count > 10 {
			last10SpeedMeasurementsSum -= last10SpeedMeasurements[0]
			last10SpeedMeasurements.remove(at: 0)
		}

		last10SpeedMeasurementsSum += lastSpeedMeasurement
		last10SpeedMeasurements.append(lastSpeedMeasurement)

		last10SpeedAverage = last10SpeedMeasurementsSum / Double(last10SpeedMeasurements.count)
	}

	//	Convert speed from meters per second to minutes per kilometer
	func convertToMinutesPerKilometer(speedMetersPerSecond: Double) -> Double {
		return speedMetersPerSecond > 0 ? (1 / (speedMetersPerSecond * (60 / 1000))) : 0
	}

	func convertToMinutesPerKilometer(speedMetersPerSecond: Double) -> String {
		if speedMetersPerSecond == 0 {
			return "00:00"
		}

		let secondsPerKilometer: Double = 1 / (speedMetersPerSecond * (1 / 1000))
		
		return formatDuration(seconds: secondsPerKilometer)
	}

	//	Format duration in seconds into HH:MM:SS
	func formatDuration(seconds: TimeInterval) -> String {
		let hours = Int(seconds / 3600)
		let minutes = Int((seconds.truncatingRemainder(dividingBy: 3600)) / 60)
		let seconds = Int(seconds.truncatingRemainder(dividingBy: 60))

		var formattedString: String
		if hours > 0 {
			formattedString = String(format: "%02d:%02d:%02d", hours, minutes, seconds)
		} else {
			formattedString = String(format: "%02d:%02d", minutes, seconds)
		}

		return formattedString
	}
}

//	MARK: - Auto pause logic
//
extension WorkoutManager {
	func autoPause(lastSpeedMeasurement: Double) {
		if !isPauseSetWithButton {
			if sessionState == .running {
				if lastSpeedMeasurement < parameters.paceForAutoPause {
					sessionState = .paused
					session?.pause()
//					print("auto pause on", "button pressed:", isPauseSetWithButton)
				}
			} else {
				if lastSpeedMeasurement > parameters.paceForAutoResume {
					sessionState = .running
					session?.resume()
//					print("auto pause off", "button pressed:", isPauseSetWithButton)
				}
			}
		}
	}
}

//	MARK: - Workout statistics
//
extension WorkoutManager {

	func updateForStatistics(_ statistics: HKStatistics) {
		switch statistics.quantityType {
			case HKQuantityType.quantityType(forIdentifier: .heartRate):
				let heartRateUnit = HKUnit.count().unitDivided(by: .minute())
				heartRate = statistics.mostRecentQuantity()?.doubleValue(for: heartRateUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .activeEnergyBurned):
				let energyUnit = HKUnit.kilocalorie()
				activeEnergy = statistics.sumQuantity()?.doubleValue(for: energyUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .distanceWalkingRunning):
				let distanceUnit = HKUnit.meter()
				distance = statistics.sumQuantity()?.doubleValue(for: distanceUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .runningPower):
				let powerUnit = HKUnit.watt()
				power = statistics.mostRecentQuantity()?.doubleValue(for: powerUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .runningSpeed):
				let speedUnit = HKUnit.meter().unitDivided(by: HKUnit.second())
				speed = statistics.mostRecentQuantity()?.doubleValue(for: speedUnit) ?? 0
				calculateLast10SpeedAverage(lastSpeedMeasurement: speed)
				autoPause(lastSpeedMeasurement: speed)


			case HKQuantityType.quantityType(forIdentifier: .runningStrideLength):
				let lengthUnit = HKUnit.meter()
				strideLength = statistics.averageQuantity()?.doubleValue(for: lengthUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .runningVerticalOscillation):
				let verticalOscillationhUnit = HKUnit.meter()
				verticalOscillation = statistics.sumQuantity()?.doubleValue(for: verticalOscillationhUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .runningGroundContactTime):
				let contactTimeUnit = HKUnit.secondUnit(with: .milli)
				groundContactTime = statistics.averageQuantity()?.doubleValue(for: contactTimeUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .vo2Max):
				let vo2MaxUnit = HKUnit.literUnit(with: .milli).unitDivided(by: HKUnit.gramUnit(with: .kilo).unitDivided(by: HKUnit.minute()))
				vo2Max = statistics.averageQuantity()?.doubleValue(for: vo2MaxUnit) ?? 0

			case HKQuantityType.quantityType(forIdentifier: .stepCount):
				stepCount = statistics.mostRecentQuantity()?.doubleValue(for: HKUnit.count()) ?? 0

			default:
				return
		}
	}
}

//	MARK: - HKWorkoutSessionDelegate
//	HealthKit calls the delegate methods on an anonymous serial background queue,
//	so the methods need to be nonisolated explicitly.
//
extension WorkoutManager: HKWorkoutSessionDelegate {
	nonisolated func workoutSession(_ workoutSession: HKWorkoutSession,
									didChangeTo toState: HKWorkoutSessionState,
									from fromState: HKWorkoutSessionState,
									date: Date) {
		Logger.shared.log("Session state changed from \(fromState.rawValue) to \(toState.rawValue)")

		//	Yield the new state change to the async stream synchronously.
		//	asynStreamTuple is a constant, so it's nonisolated
		let sessionSateChange = SessionSateChange(newState: toState, date: date)
		asynStreamTuple.continuation.yield(sessionSateChange)
	}

	nonisolated func workoutSession(_ workoutSession: HKWorkoutSession,
									didFailWithError error: Error) {
		Logger.shared.log("\(#function): \(error)")
	}


	//	HealthKit calls this method when it determines that
	//	the mirrored workout session is invalid
	nonisolated func workoutSession(_ workoutSession: HKWorkoutSession,
									didDisconnectFromRemoteDeviceWithError error: Error?) {
			Logger.shared.log("\(#function): \(error)")
	}


	//	In iOS, the sample app can go into the background and become suspended.
	//	When suspended, HealthKit gathers the data coming from the remote session.
	//	When the app resumes, HealthKit sends an array containing all the data objects
	//	it has accumulated to this delegate method.
	//	The data objects in the array appear in the order that the local system received them.
	//
	//	On watchOS, the workout session keeps the app running even if it is
	//	in the background; however, the system can temporarily suspend the app
	//	for example, if the app uses an excessive amount of CPU in the background.
	//	While suspended, HealthKit caches the incoming data objects and delivers
	//	an array of data objects when the app resumes, just like in the iOS app.
	nonisolated func workoutSession(_ workoutSession: HKWorkoutSession,
									didReceiveDataFromRemoteWorkoutSession data: [Data]) {
		Logger.shared.log("\(#function): \(data.debugDescription)")
		Task { @MainActor in
			do {
					for anElement in data {
						try handleReceivedData(anElement)
					}
			} catch {
				Logger.shared.log("Failed to handle received data: \(error))")
			}
		}
	}
}

//	MARK: - A structure for synchronizing the elapsed time.
//
struct WorkoutElapsedTime: Codable {
	var timeInterval: TimeInterval
	var date: Date
}

//	MARK: - Convenient workout state
//
extension HKWorkoutSessionState {
	var isActive: Bool {
		self != .notStarted && self != .ended
	}
}
