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
			selection = .controls
		}
		.sheet(isPresented: $isSheetActive) {
			workoutManager.resetWorkout()
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
