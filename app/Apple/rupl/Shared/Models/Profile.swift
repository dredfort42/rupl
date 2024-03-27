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
	var email: String = ""
	var firstName: String = ""
	var lastName: String = ""
	var dateOfBirth: String = ""
	var gender: String = ""

	static let shared = Profile()

	func getProfile() {
		let apiUrl = URL(string: "\(AppSettings.shared.profileURL)?client_id=\(AppSettings.shared.clientID)&access_token=\(AppSettings.shared.deviceAccessToken)")!
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

							self.email = email
							self.firstName = firstName
							self.lastName = lastName
							self.dateOfBirth = dateOfBirth
							self.gender = gender
						}
					}
				} catch {
					Logger.shared.log("Error parsing JSON: \(error)")
				}
			}
//			completion("OK")

			print(self.email)
			print(self.firstName)
			print(self.lastName)
			print(self.dateOfBirth)
			print(self.gender)
		}
		task.resume()
	}
}
