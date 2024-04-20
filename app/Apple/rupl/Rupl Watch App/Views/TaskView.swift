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
	@EnvironmentObject var workoutManager: WorkoutManager
	@Environment(\.dismiss) var dismiss

	var body: some View {

		if TaskManager.shared.task == nil {
			CreateTaskView()
		} else {
			ShowTaskView()
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

	@ViewBuilder
	private func ShowTaskView() -> some View {
		ScrollView {
			ShowHeader()
			VStack(alignment: .leading) {
				Text(TaskManager.shared.task?.description ?? "")
					.font(.title2)
					.padding()

				if let intervals = TaskManager.shared.task?.intervals {
					ForEach(intervals, id: \.self) { interval in
						TaskIntervalView(interval: interval)
							.padding()
					}
				} else {
					Text("No intervals available")
				}
			}
			SlideButton("Decline task", styling: .init(color: .ruplRed, indicatorSystemName: "xmark")) {
				TaskManager.shared.declineTask() { result in
					if result {
						dismiss()
					}
				}
			}
		}
		.onDisappear() {
		}
	}

	@ViewBuilder
	private func TaskIntervalView(interval: TaskManager.Interval) -> some View {
		Divider()
		VStack(alignment: .leading) {
			HStack {
				Text("\(interval.id)")
					.foregroundColor(.ruplGray)
					.padding(.trailing)
				Text(interval.description)
					.foregroundColor(.ruplBlue)
			}
			.font(.title3)
			.padding(.bottom)


			if interval.speed != 0 {
				Text("Speed")
					.foregroundColor(.ruplGray)
					.font(.footnote)
				Text(workoutManager.convertToMinutesPerKilometer(metersPerSecond: Double(interval.speed)) + " min/km")
					.padding(.bottom)
			} else {
				Text("Heart rate")
					.foregroundColor(.ruplGray)
					.font(.footnote)
				Text("\(TaskManager.shared.getHeartRateInterval(pz: "pz\(interval.pulse_zone)").minHeartRate) bpm - \(TaskManager.shared.getHeartRateInterval(pz: "pz\(interval.pulse_zone)").maxHeartRate) bpm")
					.padding(.bottom)
			}

			if interval.distance != 0 {
				Text("Distance")
					.foregroundColor(.ruplGray)
					.font(.footnote)
				Text("\(interval.distance / 1000) km")
					.padding(.bottom)
			} else {
				Text("Dutation")
					.foregroundColor(.ruplGray)
					.font(.footnote)
				Text("\(workoutManager.formatDuration(seconds: Double(interval.duration))) min")
					.padding(.bottom)
			}
		}
	}
}

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
