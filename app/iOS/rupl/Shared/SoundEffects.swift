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

	private let workoutStartSound: URL? = Bundle.main.url(
		forResource: "WorkoutStartSound",
		withExtension: "aif"
	)

	private let workoutStopSound: URL? = Bundle.main.url(
		forResource: "WorkoutStopSound",
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
		} catch {
			print("Error preparing sounds: \(error.localizedDescription)")
		}
	}
}
