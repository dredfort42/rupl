//
//  DeviceInfo+watchOS.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 28/03/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import WatchKit

extension DeviceInfo {
	func collectWatchInfo() -> DeviceInfo.Device {
		let currentDevice = WKInterfaceDevice.current()

		return Device(model: currentDevice.model,
					  name: currentDevice.name,
					  system: currentDevice.systemName,
					  version: currentDevice.systemVersion,
					  identifier: currentDevice.identifierForVendor?.uuidString ?? "N/A")
	}
}
