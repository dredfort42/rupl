//
//  OAuth2.swift
//  rupl
//
//  Created by Dmitry Novikov on 29/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os


struct OAuth2 {
	static var accessToken: String = ""
	static var tokenType: String = ""
	static var userCode: String = ""
	static var deviceCode: String = ""
	static var verificationUri: String = ""
	static var expiresIn: Date = Date()
	static var interval: UInt32 = 0

	static func sendRequest(completion: @escaping (String) -> Void) {
		let apiUrl = URL(string: "\(AppSettings.shared.deviceAuthURL)?client_id=\(AppSettings.shared.clientID)")!
		var request = URLRequest(url: apiUrl)
		
		request.httpMethod = "POST"
		request.setValue("application/x-www-form-urlencoded", forHTTPHeaderField: "Content-Type")

		let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
			if let error = error {
				Logger.shared.log("Error: \(error)")
				return
			}

			if let data = data {
				do {
					if let json = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any] {
						if let userCode = json["user_code"] as? String,
						   let deviceCode = json["device_code"] as? String,
						   let verificationUri = json["verification_uri"] as? String,
						   let expiresIn = json["expires_in"] as? Int,
						   let interval = json["interval"] as? Int {

							self.userCode = userCode
							self.deviceCode = deviceCode
							self.verificationUri = verificationUri
							self.expiresIn = Date() + TimeInterval(expiresIn)
							self.interval = UInt32(interval)
						}
					}
				} catch {
					Logger.shared.log("Error parsing JSON: \(error)")
				}
			}
			completion("OK")
		}
		task.resume()
	}

	static func getDeviceAccessToken(completion: @escaping (String) -> Void) {
		let apiUrl = URL(string: "\(AppSettings.shared.deviceTokenURL)?grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Adevice_code&device_code=\(deviceCode)&client_id=\(AppSettings.shared.clientID)")!
		var request = URLRequest(url: apiUrl)

		request.httpMethod = "POST"
		request.setValue("application/x-www-form-urlencoded", forHTTPHeaderField: "Content-Type")

		let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
			if let error = error {
				Logger.shared.log("Error: \(error)")
				return
			}

			if let data = data {
				do {
					if let json = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any] {
						if let accessToken = json["access_token"] as? String,
						   let tokenType = json["token_type"] as? String,
						   let expiresIn = json["expires_in"] as? Int {

							self.accessToken = accessToken
							self.tokenType = tokenType
							self.expiresIn = Date() + TimeInterval(expiresIn)
						}
					}
				} catch {
					Logger.shared.log("Error parsing JSON: \(error)")
				}
			}
			completion("OK")
		}
		task.resume()
	}

	static func deleteDeviceAccess(completion: @escaping (String) -> Void) {
		let apiUrl = URL(string: "\(AppSettings.shared.deviceAuthURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)")!
		var request = URLRequest(url: apiUrl)

		request.httpMethod = "DELETE"

		let task = URLSession.shared.dataTask(with: request) { data, response, error in
			guard let data = data else {
				print(String(describing: error))
				return
			}
			print(String(data: data, encoding: .utf8)!)
		}

		task.resume()
	}
}
