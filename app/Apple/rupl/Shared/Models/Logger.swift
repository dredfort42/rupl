//
//  Logger.swift
//  rupl
//
//  Created by Dmitry Novikov on 04/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os

extension Logger {
	private static var subsystem = Bundle.main.bundleIdentifier!
	#if os(watchOS)
		static let shared = Logger(subsystem: subsystem, category: "rupl for watch")
	#else
		static let shared = Logger(subsystem: subsystem, category: "rupl for phone")
	#endif
}
