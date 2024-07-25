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
		let apiUrl = URL(string: "\(AppSettings.shared.profileURL)")!
		var request = URLRequest(url: apiUrl)

		request.httpMethod = "GET"

		let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
			if let error = error {
				Logger.shared.log("Error: \(error)")
				return
			}

			if let data = data {
				do {

					if let json = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any] {
						if let email = json["email"] as? String,
						   let firstName = json["first_name"] as? String,
						   let lastName = json["last_name"] as? String,
						   let dateOfBirth = json["date_of_birth"] as? String,
						   let gender = json["gender"] as? String {
							let dateFormatter = DateFormatter()

							dateFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssZ"
							dateFormatter.locale = Locale(identifier: "en_US_POSIX")
							dateFormatter.timeZone = TimeZone(secondsFromGMT: 0)

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
