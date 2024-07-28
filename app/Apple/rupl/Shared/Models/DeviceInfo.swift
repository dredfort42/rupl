//
//  DeviceInfo.swift
//  rupl
//
//  Created by Dmitry Novikov on 28/03/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os

class DeviceInfo {
	struct Device {
		var model: String = "N/A"
		var name: String = "N/A"
		var system: String = "N/A"
		var version: String = "N/A"
		var identifier: String = "N/A"
	}

	var device = Device()

	static let shared = DeviceInfo()

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
		print("App Version: \(AppSettings.shared.appVersion)")
	}

	func sendDeviceInformation() {
		guard let apiUrl = URL(string: "\(AppSettings.shared.deviceInfoURL)?client_id=\(AppSettings.shared.clientID)") else {
			print("Invalid URL")
			return
		}

		var request = URLRequest(url: apiUrl)
		request.setValue("Bearer \(AppSettings.shared.deviceAccessToken)", forHTTPHeaderField: "Authorization")

		var parameters = [String: Any]()
		var jsonData = Data()

		parameters["device_uuid"] = device.identifier
		parameters["device_model"] = device.model
		parameters["device_name"] = device.name
		parameters["system_name"] = device.system
		parameters["system_version"] = device.version
		parameters["app_version"] = AppSettings.shared.appVersion

		do {
			jsonData = try JSONSerialization.data(withJSONObject: parameters, options: .prettyPrinted)
		} catch {
			Logger.shared.log("Error converting to JSON: \(error.localizedDescription)")
		}

		request.addValue("application/json", forHTTPHeaderField: "Content-Type")
		request.httpBody = jsonData

		request.httpMethod = "POST"

		let task = URLSession.shared.dataTask(with: request) { data, response, error in
			guard let data = data else {
				Logger.shared.log("Error: \(error)")
				return
			}
#if DEBUG
			print(String(data: data, encoding: .utf8)!)
#endif
		}

		task.resume()
	}

	func deleteDeviceInfo() {
		guard let apiUrl = URL(string: "\(AppSettings.shared.deviceInfoURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)") else {
			print("Invalid URL")
			return
		}
		var request = URLRequest(url: apiUrl)
		var parameters = [String: Any]()
		var jsonData = Data()

		parameters["device_model"] = device.model
		parameters["device_name"] = device.name
		parameters["system_name"] = device.system
		parameters["system_version"] = device.version
		parameters["device_id"] = device.identifier
		parameters["app_version"] = AppSettings.shared.appVersion

		do {
			jsonData = try JSONSerialization.data(withJSONObject: parameters, options: .prettyPrinted)
		} catch {
			Logger.shared.log("Error converting to JSON: \(error.localizedDescription)")
		}

		request.addValue("application/json", forHTTPHeaderField: "Content-Type")
		request.httpMethod = "DELETE"
		request.httpBody = jsonData

		let task = URLSession.shared.dataTask(with: request) { data, response, error in
//			guard let data = data else {
//				Logger.shared.log("Error: \(error)")
//				return
//			}
#if DEBUG
//			print(String(data: data, encoding: .utf8)!)
#endif
		}

		task.resume()
	}

}

