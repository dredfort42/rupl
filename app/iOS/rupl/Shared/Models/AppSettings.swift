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

	static let useAutoPauseKey = "useAutoPause"
	static let connectedToRuplKey = "connectedToRupl"
	private let clientIDKey = "clientID"
	private let permissibleHorizontalAccuracyKey = "permissibleHorizontalAccuracy"
	private let paceForAutoPauseKey = "paceForAutoPause"
	private let paceForAutoResumeKey = "paceForAutoResume"
	private let timeForShowLastSegmentViewKey = "timeForShowLastSegmentView"
	private let pz1NotInZoneKey = "pz1NotInZone"
	private let pz2EasyKey = "pz2Easy"
	private let pz3FatBurningKey = "pz3FatBurning"
	private let pz4AerobicKey = "pz4Aerobic"
	private let pz5AnaerobicKey = "pz5Anaerobic"

	var clientID: String {
		get {
			return UserDefaults.standard.string(forKey: clientIDKey) ?? ""
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

	private init() {
		if clientID == "" {
			UserDefaults.standard.set(UUID().uuidString, forKey: clientIDKey)
		}

//		if permissibleHorizontalAccuracy == 0.0 {
			permissibleHorizontalAccuracy = 16.0
//		}

		if paceForAutoPause == 0.0 || paceForAutoResume == 0.0 {
			useAutoPause = true
			paceForAutoPause = 1.85
			paceForAutoResume = 2.25
		}

//		if timeForShowLastSegmentView == 0 {
			timeForShowLastSegmentView = 20
//		}

		if pz1NotInZone == 0 || pz2Easy == 0 || pz3FatBurning == 0 || pz4Aerobic == 0 || pz5Anaerobic == 0 {
			pz1NotInZone = 126
			pz2Easy = 137
			pz3FatBurning = 147
			pz4Aerobic = 158
			pz5Anaerobic = 168
		}
	}
}
