//
//  MeasurementsConverters.swift
//  rupl
//
//  Created by Dmitry Novikov on 23/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation

//	MARK: - Workout measurements converters
//
extension WorkoutManager {
	//	Convert speed from meters per second to minutes per kilometer
	func convertToMinutesPerKilometer(metersPerSecond: Double) -> Double {
		return metersPerSecond > 0 ? (1 / (metersPerSecond * (60 / 1000))) : 0
	}

	func convertToMinutesPerKilometer(metersPerSecond: Double) -> String {
		if metersPerSecond == 0 {
			return "00:00"
		}

		let secondsPerKilometer: Double = 1 / (metersPerSecond * (1 / 1000))

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
