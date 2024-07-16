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
	static var refreshToken: String = ""
	static var tokenType: String = ""
	static var userCode: String = ""
	static var deviceCode: String = ""
	static var verificationUri: String = ""
	static var expiresIn: Date = Date.now
	static var interval: UInt32 = 0

	static func sendAuthorizeRequest(completion: @escaping (String) -> Void) {
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

	static func getDeviceTokens(completion: @escaping (String) -> Void) {
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
						   let refreshToken = json["refresh_token"] as? String,
						   let tokenType = json["token_type"] as? String,
						   let expiresIn = json["expires_in"] as? Int {

							self.accessToken = accessToken
							self.refreshToken = refreshToken
							self.tokenType = tokenType
							self.expiresIn = Date.now + TimeInterval(expiresIn - 30)
							print("A: ", accessToken)
							print("R: ", refreshToken)
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

	static func identifyDevice(completion: @escaping (String) -> Void) {
		let apiUrl = URL(string: "\(AppSettings.shared.deviceIdentifyURL)?client_id=\(AppSettings.shared.clientID)")!
		var request = URLRequest(url: apiUrl)

		request.setValue("Bearer \(AppSettings.shared.deviceAccessToken)", forHTTPHeaderField: "Authorization")
		request.httpMethod = "GET"

		let task = URLSession.shared.dataTask(with: request) { data, response, error in
			if let error = error {
				Logger.shared.log("Error: \(error)")
				return
			}

			if let httpResponse = response as? HTTPURLResponse {
				print("Status code: \(httpResponse.statusCode)")

				if httpResponse.statusCode != 200 {
					print("REFRESH TOKEN")
					// TODO: chack access token
				}
			}
		}

		task.resume()

	}

	static func deleteDevice(completion: @escaping (String) -> Void) {
// TODO: chack access token


		let apiUrl = URL(string: "\(AppSettings.shared.deviceDeleteURL)?client_id=\(AppSettings.shared.clientID)")!
		var request = URLRequest(url: apiUrl)

		request.setValue("Bearer \(AppSettings.shared.deviceAccessToken)", forHTTPHeaderField: "Authorization")
		request.httpMethod = "DELETE"

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

		AppSettings.shared.resetUser()
	}

	private func checkAccessToken(completion: @escaping (String) -> Void) {
		if AppSettings.shared.deviceAccessTokenExpiresIn < Date.now {
			let apiUrl = URL(string: "\(AppSettings.shared.deviceRefreshURL)?client_id=\(AppSettings.shared.clientID)&grant_type=refresh_token&refresh_token=\(AppSettings.shared.deviceRefreshToken)")!
			var request = URLRequest(url: apiUrl)

			request.httpMethod = "POST"

			let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
				if let error = error {
					Logger.shared.log("Error: \(error)")
					return
				}

				if let httpResponse = response as? HTTPURLResponse {
					if httpResponse.statusCode == 401 {
						AppSettings.shared.resetUser()
						completion("KO")
						return
					}
				}

				if let data = data {
					do {
						if let json = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any] {
							if let accessToken = json["access_token"] as? String,
							   let refreshToken = json["refresh_token"] as? String,
							   let tokenType = json["token_type"] as? String,
							   let expiresIn = json["expires_in"] as? Int {

								AppSettings.shared.deviceAccessToken = accessToken
								AppSettings.shared.deviceRefreshToken = refreshToken
								AppSettings.shared.deviceAccessTokenType = tokenType
								AppSettings.shared.deviceAccessTokenExpiresIn = Date.now + TimeInterval(expiresIn - 30)
								print("A: ", accessToken)
								print("R: ", refreshToken)
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
}
