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
	private let permissibleHorizontalAccuracyForAutoPause: Double = AppSettings.shared.permissibleHorizontalAccuracy * 0.7
	private let paceForAutoPause: Double = AppSettings.shared.paceForAutoPause
	private let paceForAutoResume: Double = AppSettings.shared.paceForAutoResume
	private let locationManager = CLLocationManager()
	private var accuracy: CLLocationAccuracy = 100

	var isAvailable: Bool = false
	var speed: CLLocationSpeed = 0
	var filteredLocations: [CLLocation] = []
	var autoPauseState: Bool?

	static let shared = LocationManager()

	private func requestAuthorization() {
		if !isAvailable {
#if DEBUG
			print("LocationManager.requestAuthorization()")
#endif
			locationManager.delegate = self
			locationManager.desiredAccuracy = kCLLocationAccuracyBestForNavigation
			locationManager.allowsBackgroundLocationUpdates = true

			switch locationManager.authorizationStatus {
				case .authorizedWhenInUse, .authorizedAlways:
					isAvailable = true
				default:
					locationManager.requestWhenInUseAuthorization()
			}
		}
	}

	func start() {
		requestAuthorization()
		locationManager.startUpdatingLocation()
	}

	func stop() {
		locationManager.stopUpdatingLocation()
	}

	// MARK: - CLLocationManagerDelegate
	//
	func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {

		accuracy = locations.max(by: { $0.horizontalAccuracy < $1.horizontalAccuracy })?.horizontalAccuracy ?? 100
		filteredLocations = locations.filter { (location: CLLocation) -> Bool in
			location.horizontalAccuracy <= permissibleHorizontalAccuracy
		}

		guard let location = filteredLocations.last else {
			return
		}

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
			if self.accuracy > self.permissibleHorizontalAccuracyForAutoPause {
				self.autoPauseState = nil
			} else {
				if self.speed < self.paceForAutoPause {
					self.autoPauseState = true
				} else if self.speed > self.paceForAutoResume {
					self.autoPauseState = false
				}
			}
		}
	}
}
