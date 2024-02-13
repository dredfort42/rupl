//
//  TimerManager.swift
//  rupl
//
//  Created by Dmitry Novikov on 11/02/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation

class TimerManager {
	private var timers = [Timer]()

	func start(timeInterval: TimeInterval, repeats: Bool, action: @escaping () -> Void) {
		let timer = Timer.scheduledTimer(withTimeInterval: timeInterval, repeats: repeats) { _ in
			action()
		}
		timers.append(timer)
	}

	func stop() {
		for timer in timers {
			timer.invalidate()
		}
		timers.removeAll()
	}
}
