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

	static var userCode: String = ""
	static var verificationUri: String = ""

//	init() {
//		sendRequest()
//	}

	static func sendRequest(completion: @escaping (String) -> Void) {
		let apiUrl = URL(string: "https://rupl.org/api/v1/auth/device_authorization?client_id=\(AppSettings.shared.clientID)")!
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
//						   let deviceCode = json["device_code"] as? String,
//						   let verificationUriComplete = json["verification_uri_complete"] as? String,
//						   let expiresIn = json["expires_in"] as? Int,
//						   let interval = json["interval"] as? Int,
						   let verificationUri = json["verification_uri"] as? String {

							self.userCode = userCode
							self.verificationUri = verificationUri

//							print("Device Code: \(deviceCode)")
//							print("User Code: \(userCode)")
//							print("Verification URI: \(verificationUri)")
//							print("Verification URI Complete: \(verificationUriComplete)")
//							print("Expires In: \(expiresIn)")
//							print("Interval: \(interval)")
							
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
}
