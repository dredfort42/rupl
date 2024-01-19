//
//  LocationManager.swift
//  rupl
//
//  Created by Dmitry Novikov on 11/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import os
import CoreLocation
import HealthKit
import SwiftUI

class LocationManager: NSObject, CLLocationManagerDelegate {
	private let parameters = WorkoutParameters()
	let locationManager = CLLocationManager()
	var autoPauseState: Bool = false
	var speed: CLLocationSpeed = 0
	var filteredLocations: [CLLocation] = []

	override init() {
		super.init()

		// Set up CLLocationManager
		locationManager.delegate = self
		locationManager.desiredAccuracy = kCLLocationAccuracyBest
		locationManager.requestWhenInUseAuthorization()
		locationManager.startUpdatingLocation()
	}

	// MARK: - CLLocationManagerDelegate
	//
	func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {

		// Filter the raw data.
		filteredLocations = locations.filter { (location: CLLocation) -> Bool in
			location.horizontalAccuracy <= parameters.minHorizontalAccuracy
		}

		guard !filteredLocations.isEmpty else { return }
		guard let location = filteredLocations.last else { return }

		// Access the speed property from the location object
		speed = location.speed
//		print("Current Speed: \(speed) meters per second")

		autoPause(speed: speed)
	}

	func locationManager(_ manager: CLLocationManager, didFailWithError error: Error) {
		Logger.shared.log("Location manager error: \(error.localizedDescription)")
	}
}

//	MARK: - Auto pause based on CLLocationManager
//
extension LocationManager {
	func autoPause(speed: Double) {
		if speed < parameters.paceForAutoPause {
			autoPauseState = true
		} else if speed > parameters.paceForAutoResume {
			autoPauseState = false
		}
	}
}
