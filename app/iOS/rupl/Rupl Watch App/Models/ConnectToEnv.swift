//
//  ConnectToEnv.swift
//  rupl
//
//  Created by Dmitry Novikov on 29/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os

class ConnectToEnv {
	let apiUrl = URL(string: "http://rupl.org/api/device_authorization")!
	let clientID = AppSettings.shared.clientID

	init() {
		print("Vendor ID: \(clientID)")
		sendRequest(requestData: ["client_id": clientID])
	}

	func sendRequest(requestData: [String: String]) {

		if let jsonData = try? JSONSerialization.data(withJSONObject: requestData) {

			// Create a URL request
			var request = URLRequest(url: apiUrl)
			request.httpMethod = "POST"
			request.setValue("application/x-www-form-urlencoded", forHTTPHeaderField: "Content-Type")
			request.httpBody = jsonData


			// Perform the network request
			let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
				if let error = error {
					Logger.shared.log("Error: \(error)")
					return
				}

				if let data = data {
					do {
						// Parse the received JSON data
						if let json = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any] {
							// Extract and handle the random 6 digits
							if let randomDigits = json["randomDigits"] as? Int {
								print("Received random 4 digits: \(randomDigits)")
							} else {
								Logger.shared.log("Error: Unable to retrieve random digits from response.")
							}
						}
					} catch {
						Logger.shared.log("Error parsing JSON: \(error)")
					}
				}
			}

			// Start the network task
			task.resume()
		} else {
			Logger.shared.log("Error converting request data to JSON")
		}
	}
}
