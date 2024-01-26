//
//  SoundEffects.swift
//  rupl
//
//  Created by Dmitry Novikov on 19/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import AVFoundation
import os

struct SoundEffects {
	var startSound: AVAudioPlayer?
	var stopSound: AVAudioPlayer?
	var segmentSound: AVAudioPlayer?
	var alarmSound: AVAudioPlayer?
	var runFaster: AVAudioPlayer?
	var runSlower: AVAudioPlayer?

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
}
