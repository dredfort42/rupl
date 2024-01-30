//
//  ActivityRingsView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 06/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import Foundation
import HealthKit
import SwiftUI

struct ActivityRingsView: WKInterfaceObjectRepresentable {
	@Environment(\.calendar) private var calendar
	let healthStore: HKHealthStore

	func makeWKInterfaceObject(context: Context) -> some WKInterfaceObject {
		let activityRingsObject = WKInterfaceActivityRing()
		var components = calendar.dateComponents([.era, .year, .month, .day], from: Date())
		components.calendar = calendar

		let predicate = HKQuery.predicateForActivitySummary(with: components)
		let query = HKActivitySummaryQuery(predicate: predicate) { query, summaries, error in
			DispatchQueue.main.async {
				activityRingsObject.setActivitySummary(summaries?.first, animated: true)
			}
		}
		healthStore.execute(query)
		return activityRingsObject
	}

	func updateWKInterfaceObject(_ wkInterfaceObject: WKInterfaceObjectType, context: Context) {
	}
}
