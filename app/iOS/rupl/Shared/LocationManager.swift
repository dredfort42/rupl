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
	var last3SpeedMeasurements: [CLLocationSpeed] = []

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

		last3SpeedMeasurements.append(speed)
		while last3SpeedMeasurements.count > 3 {
			last3SpeedMeasurements.remove(at: 0)
		}

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
		let checkPause = last3SpeedMeasurements.filter { (measurment: CLLocationSpeed) -> Bool in
			measurment < parameters.paceForAutoPause
		}

//		let checkResume = last3SpeedMeasurements.filter { (measurment: CLLocationSpeed) -> Bool in
//			measurment > parameters.paceForAutoResume
//		}

		if checkPause.count == 3 {
			autoPauseState = true
		} else if speed > parameters.paceForAutoResume {
			autoPauseState = false
		}
	}
}
