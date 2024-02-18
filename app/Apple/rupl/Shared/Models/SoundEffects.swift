//
//  SoundEffects.swift
//  rupl
//
//  Created by Dmitry Novikov on 19/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import AVFoundation
import os
#if os(watchOS)
import WatchKit
#endif

struct SoundEffects {
	private var startSound: AVAudioPlayer?
	private var stopSound: AVAudioPlayer?
	private var segmentSound: AVAudioPlayer?
	private var alarmSound: AVAudioPlayer?
	private var runFaster: AVAudioPlayer?
	private var runSlower: AVAudioPlayer?

	private let workoutStartSound: URL? = Bundle.main.url(
		forResource: "WorkoutStartSound",
		withExtension: "aif"
	)

	private let workoutStopSound: URL? = Bundle.main.url(
		forResource: "WorkoutStopSound",
		withExtension: "aif"
	)

	private let workoutSegmentSound: URL? = Bundle.main.url(
		forResource: "WorkoutSegmentSound",
		withExtension: "aif"
	)

	private let workoutAlarmSound: URL? = Bundle.main.url(
		forResource: "WorkoutAlarmSound",
		withExtension: "aif"
	)

	private let workoutRunFasterSound: URL? = Bundle.main.url(
		forResource: "WorkoutRunFasterSound",
		withExtension: "aif"
	)

	private let workoutRunSlowerSound: URL? = Bundle.main.url(
		forResource: "WorkoutRunSlowerSound",
		withExtension: "aif"
	)

	static let shared = SoundEffects()

	init() {
		do {
			try AVAudioSession.sharedInstance().setCategory(
				AVAudioSession.Category.playback,
				options: AVAudioSession.CategoryOptions.mixWithOthers
			)

			try AVAudioSession.sharedInstance().setActive(true)

			if workoutStartSound != nil {
				startSound = try AVAudioPlayer(contentsOf: workoutStartSound!)
			}
			if workoutStopSound != nil {
				stopSound = try AVAudioPlayer(contentsOf: workoutStopSound!)
			}
			if workoutSegmentSound != nil {
				segmentSound = try AVAudioPlayer(contentsOf: workoutSegmentSound!)
			}
			if workoutAlarmSound != nil {
				alarmSound = try AVAudioPlayer(contentsOf: workoutAlarmSound!)
			}
			if workoutRunFasterSound != nil {
				runFaster = try AVAudioPlayer(contentsOf: workoutRunFasterSound!)
			}
			if workoutRunSlowerSound != nil {
				runSlower = try AVAudioPlayer(contentsOf: workoutRunSlowerSound!)
			}
		} catch {
			Logger.shared.log("Error preparing sounds: \(error.localizedDescription)")
		}
	}

	func playStartSound() {
#if targetEnvironment(simulator)
		print("))*(( Start sound")
#else
		startSound?.play()
#endif
		Vibration.vibrate(type: .notification)
	}

	func playStopSound() {
#if targetEnvironment(simulator)
		print("))*(( Stop sound")
#else
		stopSound?.play()
#endif
		Vibration.vibrate(type: .notification)
	}

	func playSegmentSound() {
#if targetEnvironment(simulator)
		print("))*(( Segment sound")
#else
		segmentSound?.play()
#endif
		Vibration.vibrate(type: .info)
	}

	func playAlarmSound() {
#if targetEnvironment(simulator)
		print("))*(( Alarm sound")
#else
		alarmSound?.play()
#endif
		Vibration.vibrate(type: .alarm)
	}

	func playRunFasterSound() {
#if targetEnvironment(simulator)
		print("))*(( Run faster sound")
#else
		runFaster?.play()
#endif
		Vibration.vibrate(type: .up)
	}

	func playRunSlowerSound() {
#if targetEnvironment(simulator)
		print("))*(( Run slower sound")
#else
		runSlower?.play()
#endif
		Vibration.vibrate(type: .down)
	}

}

class Vibration {
	enum VibrationType {
		case notification
		case info
		case alarm
		case click
		case up
		case down
	}

	static func vibrate(type: VibrationType) {
#if targetEnvironment(simulator)
		print("Zz. Vibration \(type)")
#elseif os(watchOS)
		var hapticType: WKHapticType
		switch type {
			case .notification:
				hapticType = .notification
			case .info:
				hapticType = .success
			case .alarm:
				hapticType = .underwaterDepthCriticalPrompt
			case .click:
				hapticType = .click
			case .up:
				hapticType = .directionUp
			case .down:
				hapticType = .directionDown
		}
		WKInterfaceDevice.current().play(hapticType)
#endif
	}
}

