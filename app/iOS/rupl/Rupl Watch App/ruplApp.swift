//
//  ruplApp.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 04/01/2024.
//

import SwiftUI

@main
struct Rupl_Watch_AppApp: App {
	private let workoutManager = WorkoutManager.shared
    
	@SceneBuilder var body: some Scene {
		WindowGroup {
			ContentView()
				.environmentObject(workoutManager)
		}
	}
}
