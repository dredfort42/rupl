//
//  ruplApp.swift
//  rupl
//
//  Created by Dmitry Novikov on 04/01/2024.
//

import SwiftUI

@main
struct ruplApp: App {
	private let workoutManager = WorkoutManager.shared

	var body: some Scene {
		WindowGroup {
			if UIDevice.current.userInterfaceIdiom == .phone {
				StartView()
					.environmentObject(workoutManager)
			} else {
				WorkoutListView()
			}
		}
	}
}
