//
//  Profile.swift
//  rupl
//
//  Created by Dmitry Novikov on 27/03/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os

class Profile {
	static func getProfile() {
		OAuth2.checkAccessToken() { result in
			if result == "KO" {
				return
			}
		}

		let apiUrl = URL(string: "\(AppSettings.shared.profileURL)?client_id=\(AppSettings.shared.clientID)")!
		var request = URLRequest(url: apiUrl)

		request.setValue("Bearer \(AppSettings.shared.deviceAccessToken)", forHTTPHeaderField: "Authorization")
		request.httpMethod = "GET"

		let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
			if let error = error {
				Logger.shared.log("Error: \(error)")
				return
			}

			if let data = data {
				do {

					print(String(data: data, encoding: .utf8)!)
					
					if let json = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any] {
						if let email = json["email"] as? String,
						   let firstName = json["first_name"] as? String,
						   let lastName = json["last_name"] as? String,
						   let dateOfBirth = json["date_of_birth"] as? String,
						   let gender = json["gender"] as? String {
							let dateFormatter = DateFormatter()

							dateFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssZ"
							dateFormatter.timeZone = TimeZone(abbreviation: "UTC")

							AppSettings.shared.userEmail = email
							AppSettings.shared.userFirstName = firstName
							AppSettings.shared.userLastName = lastName
							AppSettings.shared.userDateOfBirth = dateFormatter.date(from: dateOfBirth)
							AppSettings.shared.userGender = gender
						}
					}
				} catch {
					Logger.shared.log("Error parsing JSON: \(error)")
				}
			}
		}
		task.resume()
	}
}
