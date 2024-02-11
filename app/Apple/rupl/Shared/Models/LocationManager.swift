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

class LocationManager: NSObject, CLLocationManagerDelegate {
	private let permissibleHorizontalAccuracy: Double = AppSettings.shared.permissibleHorizontalAccuracy
	private let paceForAutoPause: Double = AppSettings.shared.paceForAutoPause
	private let paceForAutoResume: Double = AppSettings.shared.paceForAutoResume
	private let locationManager = CLLocationManager()
	private var autoPauseIndicator: Int = 0

	var speed: CLLocationSpeed = 0
	var accuracy: CLLocationAccuracy = 1000
	var filteredLocations: [CLLocation] = []
	var autoPauseState: Bool = true

	override init() {
		super.init()

		locationManager.delegate = self
		locationManager.desiredAccuracy = kCLLocationAccuracyBestForNavigation
		locationManager.allowsBackgroundLocationUpdates = true
		locationManager.requestWhenInUseAuthorization()
		start()
	}

	func start() {
		locationManager.startUpdatingLocation()
	}

	func stop() {
		locationManager.stopUpdatingLocation()
	}

	// MARK: - CLLocationManagerDelegate
	//
	func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {

		accuracy = locations.max(by: { $0.horizontalAccuracy < $1.horizontalAccuracy })?.horizontalAccuracy ?? 1000

		// Filter the raw data.
		filteredLocations = locations.filter { (location: CLLocation) -> Bool in
			location.horizontalAccuracy <= permissibleHorizontalAccuracy
		}

		guard !filteredLocations.isEmpty else {
			if autoPauseIndicator > 0 {
				autoPauseIndicator -= 1
			}
			return
		}
		
		guard let location = filteredLocations.last else { return }

		// Access the speed property from the location object
		speed = location.speed

		checkAutoPause()
	}

	func locationManager(_ manager: CLLocationManager, didFailWithError error: Error) {
		Logger.shared.log("Location manager error: \(error.localizedDescription)")
	}
}

//	MARK: - Auto pause based on CoreLocation
//
extension LocationManager {
	func checkAutoPause() {
		DispatchQueue.main.async {
			if self.speed < self.paceForAutoPause {
				self.autoPauseIndicator += 1
			} else if self.autoPauseIndicator > 0 {
				self.autoPauseIndicator -= 1
			}

			if self.speed > self.paceForAutoResume {
				self.autoPauseIndicator = 0
				self.autoPauseState = false
			} else if self.autoPauseIndicator == 5 {
				self.autoPauseState = true
			}
		}
	}
}
