//
//  TaskView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 21/02/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct TaskView: View {
	@AppStorage(AppSettings.runningTaskHeartRateKey) var runningTaskHeartRate = AppSettings.shared.runningTaskHeartRate

	@Environment(\.dismiss) var dismiss

	@State private var isNewTaskAvailable = TaskManager.shared.isNewRunTaskAvailable

	var body: some View {

		if isNewTaskAvailable {
			ShowTaskView()
		} else {
			if TaskManager.shared.isRunTaskAccepted == true && !TaskManager.shared.isRunTaskStarted  {
				ShowTaskView()
			} else {
				CreateTaskView()
			}
		}
	}

	@ViewBuilder
	private func ShowHeader() -> some View {
		HStack(content: {
			Spacer()
			Text("Running task")
				.font(.headline)
				.foregroundColor(.ruplBlue)
		})
		.padding(.horizontal)
		.padding(.top, -25)
		.padding(.bottom)
	}

	@ViewBuilder
	private func ShowTaskView() -> some View {
		ScrollView {
			ShowHeader()
			Text(TaskManager.shared.task?.description ?? "")

			if let intervals = TaskManager.shared.task?.intervals {
				ForEach(intervals, id: \.self) { interval in
					TaskIntervalView(interval: interval)
						.padding()
				}
			} else {
				Text("No intervals available")
			}


			Button {
				TaskManager.shared.isRunTaskAccepted = false
				dismiss()
			} label: {
				Text("Decline")
			}
			.padding(.vertical)
		}
		.onDisappear() {
		}
		//		.onAppear() {
		//		}
	}

	struct TaskIntervalView: View {
		var interval: TaskManager.Interval

		var body: some View {
			let title: String = interval.description
			let intencety: String = interval.speed != 0 ? String(interval.speed) : String(interval.pulse_zone)
			let duration: String = interval.duration != 0 ? String(interval.distance) : String(interval.distance)

			Divider()
			Text(title)
				.foregroundStyle(.foreground)
			Text(intencety)
				.font(.system(.title2, design: .rounded).lowercaseSmallCaps())
			Text(duration)
				.font(.system(.title2, design: .rounded).lowercaseSmallCaps())
		}
	}

	@ViewBuilder
	private func CreateTaskView() -> some View {
		ScrollView {
			ShowHeader()
			Picker("Heart rate zone", selection: $runningTaskHeartRate) {
				Text("Easy").tag(TaskManager.HeartRateZones.pz1.rawValue)
					.foregroundColor(.ruplBlue)
				Text("Endurance").tag(TaskManager.HeartRateZones.pz2.rawValue)
					.foregroundColor(.ruplGreen)
				Text("Tempo").tag(TaskManager.HeartRateZones.pz3.rawValue)
					.foregroundColor(.ruplGreen)
				Text("Threshold").tag(TaskManager.HeartRateZones.pz4.rawValue)
					.foregroundColor(.ruplYellow)
				Text("Anaerobic").tag(TaskManager.HeartRateZones.pz5.rawValue)
					.foregroundColor(.ruplRed)
				Text("Any").tag(TaskManager.HeartRateZones.any.rawValue)
					.foregroundColor(.ruplGray)
			}
			.frame(height: 80)
			.pickerStyle(WheelPickerStyle())

			Text("\(TaskManager.shared.getHeartRateInterval(pz: runningTaskHeartRate).minHeartRate) bpm - \( TaskManager.shared.getHeartRateInterval(pz: runningTaskHeartRate).maxHeartRate) bpm")
				.font(.caption2)
				.foregroundColor(.ruplGray)
				.padding(.horizontal)

			Spacer()
		}
		.onDisappear() {
			TaskManager.shared.intervalHeartRateZone =  TaskManager.shared.getHeartRateInterval(pz: runningTaskHeartRate)
		}
	}
}

//#Preview {
//	TaskView()
//}

//Zone 1: Recovery/Easy (50-60% of MHR)
//Lower limit: 0.50 x MHR
//Upper limit: 0.60 x MHR
//Zone 2: Aerobic/Endurance (60-70% of MHR)
//Lower limit: 0.60 x MHR
//Upper limit: 0.70 x MHR
//Zone 3: Tempo (70-80% of MHR)
//Lower limit: 0.70 x MHR
//Upper limit: 0.80 x MHR
//Zone 4: Threshold (80-90% of MHR)
//Lower limit: 0.80 x MHR
//Upper limit: 0.90 x MHR
//Zone 5: Anaerobic/VO2 Max (90-100% of MHR)
//Lower limit: 0.90 x MHR
//Upper limit: 1.00 x MHR
