//
//  ContentView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 04/01/2024.
//

import SwiftUI

struct ContentView: View {
	@EnvironmentObject var workoutManager: WorkoutManager
	@Environment(\.isLuminanceReduced) var isLuminanceReduced
	@State private var selection: Tab = .metrics
	@State private var isSheetActive = false

	private let appVersion = Bundle.main.infoDictionary?["CFBundleShortVersionString"] as? String ?? "N/A"
	private let buildNumber = Bundle.main.infoDictionary?["CFBundleVersion"] as? String ?? "N/A"

	private enum Tab {
		case controls, metrics
	}

	var body: some View {
		TabView(selection: $selection) {
			ControlsView().tag(Tab.controls)
			MetricsView().tag(Tab.metrics)
		}
		.navigationTitle("Run")
		.navigationBarBackButtonHidden(true)
		.tabViewStyle(PageTabViewStyle(indexDisplayMode: isLuminanceReduced ? .never : .automatic))
		.onChange(of: isLuminanceReduced) {
			displayMetricsView()
		}
		.onChange(of: workoutManager.sessionState) { _, newValue in
			if newValue == .stopped {
				isSheetActive = true
			} else if newValue == .running || newValue == .paused {
				displayMetricsView()
			}
		}
		.onAppear {
			AppSettings.shared.appVersion = appVersion + "." + buildNumber
			selection = .controls
			// MARK: - get profile
			if AppSettings.shared.connectedToRupl {
				Profile.getProfile()
				DeviceInfo.shared.sendDeviceInformation(createNew: false)
				TaskManager.shared.getTask { result in }
			}
		}
		.sheet(isPresented: $isSheetActive) {
			if workoutManager.sessionState != .ended {
				workoutManager.saveWorkout()
			}
		} content: {
			SummaryView()
		}
	}

	private func displayMetricsView() {
		withAnimation {
			selection = .metrics
		}
	}
}
