//
//  Vibration.swift
//  rupl
//
//  Created by Dmitry Novikov on 19/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import WatchKit

class Vibration: WKInterfaceController {
	static func vibrate() {
		WKInterfaceDevice.current().play(.notification)
	}
}
