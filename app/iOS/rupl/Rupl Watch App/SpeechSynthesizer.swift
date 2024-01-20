//
//  SpeechSynthesizer.swift
//  rupl
//
//  Created by Dmitry Novikov on 20/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import WatchKit
import AVFoundation

class SpeechSynthesizer: WKInterfaceController {



	static func convertTextToSpeech() {
		let speechSynthesizer = AVSpeechSynthesizer()
		// Specify the text you want to convert to speech
		let textToSpeak = "Hello, this is a sample text."

		// Create an AVSpeechUtterance with the text
		let speechUtterance = AVSpeechUtterance(string: textToSpeak)

		// Configure speech settings if needed
		speechUtterance.rate = AVSpeechUtteranceDefaultSpeechRate
		speechUtterance.voice = AVSpeechSynthesisVoice(language: "en-US")

		// Use the speech synthesizer to speak the text
		speechSynthesizer.speak(speechUtterance)
	}

}

