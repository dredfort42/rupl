//
//  DeviceInfo.swift
//  rupl
//
//  Created by Dmitry Novikov on 28/03/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation

class DeviceInfo {
	struct Device {
		var model: String = "N/A"
		var name: String = "N/A"
		var system: String = "N/A"
		var version: String = "N/A"
		var identifier: String = "N/A"
	}

	var device = Device()

	init() {
#if os(watchOS)
		device = collectWatchInfo()
#endif
	}

	func printDeviceInfo() {
		print("Device Model: \(device.model)")
		print("Device Name: \(device.name)")
		print("System Name: \(device.system)")
		print("System Version: \(device.version)")
		print("Device Identifier: \(device.identifier)")
	}
}
