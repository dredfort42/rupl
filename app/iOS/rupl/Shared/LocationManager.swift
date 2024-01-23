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
	var autoPauseIndicator: Int = 0

	override init() {
		super.init()

		// Set up CLLocationManager
		locationManager.delegate = self
		locationManager.desiredAccuracy = kCLLocationAccuracyBestForNavigation
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

		checkAutoPause()
	}

	func locationManager(_ manager: CLLocationManager, didFailWithError error: Error) {
		Logger.shared.log("Location manager error: \(error.localizedDescription)")
	}
}

//	MARK: - Auto pause based on CLLocationManager
//
extension LocationManager {
	func checkAutoPause() {
		DispatchQueue.main.async {
			if self.speed < self.parameters.paceForAutoPause {
				self.autoPauseIndicator += 1
			} else if self.autoPauseIndicator > 0 {
				self.autoPauseIndicator -= 1
			}

			if self.speed > self.parameters.paceForAutoResume {
				self.autoPauseIndicator = 0
				self.autoPauseState = false
			} else if self.autoPauseIndicator == 5 {
				self.autoPauseState = true
			}
		}
	}
}
