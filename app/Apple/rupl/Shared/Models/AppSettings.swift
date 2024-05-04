//
//  AppSettings.swift
//  rupl
//
//  Created by Dmitry Novikov on 08/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation

//	MARK: - A class for declare app serrings constants
//
class AppSettings {
	static let shared = AppSettings()

	// User app settings
	static let useAutoPauseKey = "useAutoPause"
	static let isSoundNotificationOnKey = "isSoundNotificationOn"
	static let criticalHeartRateKey = "criticalHeartRate"
	static let connectedToRuplKey = "connectedToRupl"
	static let runningTaskHeartRateKey = "runningTaskHeartRate"

	// Connection data
	private let appVersionKey = "appVersion"
	private let clientIDKey = "clientID"
	private let deviceAuthURLKey = "deviceAuthURL"
	private let deviceTokenURLKey = "deviceTokenURL"
	private let profileURLKey = "profileURL"
	private let deviceInfoURLKey = "deviceInfoURL"
	private let taskURLKey = "taskURL"
	private let sessionURLKey = "sessionURL"
	private let deviceAccessTokenKey = "deviceAccessToken"
	private let deviceAccessTokenTypeKey = "deviceAccessTokenType"
	private let deviceAccessTokenExpiresInKey = "deviceAccessTokenExpiresIn"

	// User profile data
	private let userEmailKey = "userEmail"
	private let userFirstNameKey = "userFirstName"
	private let userLastNameKey = "userLastName"
	private let userDateOfBirthKey = "userDateOfBirth"
	private let userGenderKey = "userGender"

	// General app settings
	private let permissibleHorizontalAccuracyKey = "permissibleHorizontalAccuracy"
	private let paceForAutoPauseKey = "paceForAutoPause"
	private let paceForAutoResumeKey = "paceForAutoResume"
	private let accelerationForAutoPauseKey = "accelerationForAutoPause"
	private let accelerationForAutoResumeKey = "accelerationForAutoResume"
	private let timeForShowLastSegmentViewKey = "timeForShowLastSegmentView"
	private let pz1NotInZoneKey = "pz1NotInZone"
	private let pz2EasyKey = "pz2Easy"
	private let pz3FatBurningKey = "pz3FatBurning"
	private let pz4AerobicKey = "pz4Aerobic"
	private let pz5AnaerobicKey = "pz5Anaerobic"
	private let soundNotificationTimeOutKey = "soundNotificationTimeOut"
	private let viewNotificationTimeOutKey = "viewNotificationTimeOut"

	var appVersion: String {
		get {
			return UserDefaults.standard.string(forKey: appVersionKey) ?? ""
		}
		set {
			UserDefaults.standard.set(newValue, forKey: appVersionKey)
		}
	}

	var clientID: String {
		get {
			return UserDefaults.standard.string(forKey: clientIDKey) ?? ""
		}
	}

	var deviceAuthURL: String {
		get {
			return UserDefaults.standard.string(forKey: deviceAuthURLKey) ?? ""
		}
	}

	var deviceTokenURL: String {
		get {
			return UserDefaults.standard.string(forKey: deviceTokenURLKey) ?? ""
		}
	}

	var profileURL: String {
		get {
			return UserDefaults.standard.string(forKey: profileURLKey) ?? ""
		}
	}

	var deviceInfoURL: String {
		get {
			return UserDefaults.standard.string(forKey: deviceInfoURLKey) ?? ""
		}
	}

	var taskURL: String {
		get {
			return UserDefaults.standard.string(forKey: taskURLKey) ?? ""
		}
	}

	var sessionURL: String {
		get {
			return UserDefaults.standard.string(forKey: sessionURLKey) ?? ""
		}
	}

	var deviceAccessToken: String {
		get {
			return UserDefaults.standard.string(forKey: deviceAccessTokenKey) ?? ""
		}
		set {
			UserDefaults.standard.set(newValue, forKey: deviceAccessTokenKey)
		}
	}

	var deviceAccessTokenType: String {
		get {
			return UserDefaults.standard.string(forKey: deviceAccessTokenTypeKey) ?? ""
		}
		set {
			UserDefaults.standard.set(newValue, forKey: deviceAccessTokenTypeKey)
		}
	}

	var deviceAccessTokenExpiresIn: Double {
		get {
			return UserDefaults.standard.double(forKey: deviceAccessTokenExpiresInKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: deviceAccessTokenExpiresInKey)
		}
	}

	var useAutoPause: Bool {
		get {
			return UserDefaults.standard.bool(forKey: AppSettings.useAutoPauseKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: AppSettings.useAutoPauseKey)
		}
	}

	var isSoundNotificationOn: Bool {
		get {
			return UserDefaults.standard.bool(forKey: AppSettings.isSoundNotificationOnKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: AppSettings.isSoundNotificationOnKey)
		}
	}

	var userEmail: String {
		get {
			return UserDefaults.standard.string(forKey: userEmailKey) ?? ""
		}
		set {
			UserDefaults.standard.set(newValue, forKey: userEmailKey)
		}
	}

	var userFirstName: String {
		get {
			return UserDefaults.standard.string(forKey: userFirstNameKey) ?? ""
		}
		set {
			UserDefaults.standard.set(newValue, forKey: userFirstNameKey)
		}
	}

	var userLastName: String {
		get {
			return UserDefaults.standard.string(forKey: userLastNameKey) ?? ""
		}
		set {
			UserDefaults.standard.set(newValue, forKey: userLastNameKey)
		}
	}

	var userDateOfBirth: Date? {
		get {
			return UserDefaults.standard.object(forKey: userDateOfBirthKey) as? Date
		}
		set {
			UserDefaults.standard.set(newValue, forKey: userDateOfBirthKey)
		}
	}

	var userGender: String {
		get {
			return UserDefaults.standard.string(forKey: userGenderKey) ?? ""
		}
		set {
			UserDefaults.standard.set(newValue, forKey: userGenderKey)
		}
	}

	var criticalHeartRate: Int {
		get {
			return UserDefaults.standard.integer(forKey: AppSettings.criticalHeartRateKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: AppSettings.criticalHeartRateKey)
		}
	}

	var runningTaskHeartRate: String {
		get {
			return UserDefaults.standard.string(forKey: AppSettings.runningTaskHeartRateKey) ?? ""
		}
		set {
			UserDefaults.standard.set(newValue, forKey: AppSettings.runningTaskHeartRateKey)
		}
	}

	var connectedToRupl: Bool {
		get {
			return UserDefaults.standard.bool(forKey: AppSettings.connectedToRuplKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: AppSettings.connectedToRuplKey)
		}
	}

	var permissibleHorizontalAccuracy: Double {
		get {
			return UserDefaults.standard.double(forKey: permissibleHorizontalAccuracyKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: permissibleHorizontalAccuracyKey)
		}
	}

	var paceForAutoPause: Double {
		get {
			return UserDefaults.standard.double(forKey: paceForAutoPauseKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: paceForAutoPauseKey)
		}
	}

	var paceForAutoResume: Double {
		get {
			return UserDefaults.standard.double(forKey: paceForAutoResumeKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: paceForAutoResumeKey)
		}
	}

	var accelerationForAutoPause: Double {
		get {
			return UserDefaults.standard.double(forKey: accelerationForAutoPauseKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: accelerationForAutoPauseKey)
		}
	}
	
	var accelerationForAutoResume: Double {
		get {
			return UserDefaults.standard.double(forKey: accelerationForAutoResumeKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: accelerationForAutoResumeKey)
		}
	}

	var timeForShowLastSegmentView: Int {
		get {
			return UserDefaults.standard.integer(forKey: timeForShowLastSegmentViewKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: timeForShowLastSegmentViewKey)
		}
	}

	var pz1NotInZone: Int {
		get {
			return UserDefaults.standard.integer(forKey: pz1NotInZoneKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: pz1NotInZoneKey)
		}
	}

	var pz2Easy: Int {
		get {
			return UserDefaults.standard.integer(forKey: pz2EasyKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: pz2EasyKey)
		}
	}

	var pz3FatBurning: Int {
		get {
			return UserDefaults.standard.integer(forKey: pz3FatBurningKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: pz3FatBurningKey)
		}
	}

	var pz4Aerobic: Int {
		get {
			return UserDefaults.standard.integer(forKey: pz4AerobicKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: pz4AerobicKey)
		}
	}

	var pz5Anaerobic: Int {
		get {
			return UserDefaults.standard.integer(forKey: pz5AnaerobicKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: pz5AnaerobicKey)
		}
	}

	var soundNotificationTimeOut: Int {
		get {
			return UserDefaults.standard.integer(forKey: soundNotificationTimeOutKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: soundNotificationTimeOutKey)
		}
	}

	var viewNotificationTimeOut: Int {
		get {
			return UserDefaults.standard.integer(forKey: viewNotificationTimeOutKey)
		}
		set {
			UserDefaults.standard.set(newValue, forKey: viewNotificationTimeOutKey)
		}
	}

	private init() {
		let age: Int = userDateOfBirth != nil ? getUserAge(dateOfBirth: userDateOfBirth!) : 21
		let maximumHeartRate: Double = 208 - Double(age) * 0.7

		if clientID == "" {
			UserDefaults.standard.set(UUID().uuidString, forKey: clientIDKey)
		}

		UserDefaults.standard.set("https://rupl.org/api/v1/auth/device_authorization", forKey: deviceAuthURLKey)
		UserDefaults.standard.set("https://rupl.org/api/v1/auth/device_token", forKey: deviceTokenURLKey)
		UserDefaults.standard.set("https://rupl.org/api/v1/profile", forKey: profileURLKey)
		UserDefaults.standard.set("https://rupl.org/api/v1/profile/devices", forKey: deviceInfoURLKey)
		UserDefaults.standard.set("https://rupl.org/api/v1/training/task", forKey: taskURLKey)
		UserDefaults.standard.set("https://rupl.org/api/v1/training/session", forKey: sessionURLKey)

		if paceForAutoPause == 0.0 || paceForAutoResume == 0.0 {
			useAutoPause = true
			isSoundNotificationOn = true
		}

		if criticalHeartRate == 0 {
			criticalHeartRate = 190
		}

		pz1NotInZone = Int(maximumHeartRate * 0.6 + 0.5)
		pz2Easy = Int(maximumHeartRate * 0.7 + 0.5)
		pz3FatBurning = Int(maximumHeartRate * 0.8 + 0.5)
		pz4Aerobic = Int(maximumHeartRate * 0.9 + 0.5)
		pz5Anaerobic = Int(maximumHeartRate + 0.5)
		permissibleHorizontalAccuracy = 16.0
		paceForAutoPause = 1.75
		paceForAutoResume = 1.95
		accelerationForAutoPause = 0.3
		accelerationForAutoResume = 0.75
		timeForShowLastSegmentView = 20
		soundNotificationTimeOut = 10
		viewNotificationTimeOut = 20
	}

	func getUserAge(dateOfBirth: Date) -> Int {
		let calendar = Calendar.current
		let ageComponents = calendar.dateComponents([.year], from: dateOfBirth, to: Date())

		if let years = ageComponents.year {
			return years
		} else {
			return 21
		}
	}
}
