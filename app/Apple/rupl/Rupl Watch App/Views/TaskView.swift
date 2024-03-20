//
//  TaskView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 21/02/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct TaskView: View {

	//	private enum RunningTask: String {
	//		case no = "No task"
	//		case pz1 = "Easy"
	//		case pz2 = "Endurance"
	//		case pz3 = "Tempo"
	//		case pz4 = "Threshold"
	//		case pz5 = "Anaerobic"
	//	}



	@AppStorage(AppSettings.runningTaskHeartRateKey) var runningTaskHeartRate = AppSettings.shared.runningTaskHeartRate
	//	@State private var selectedTask = TaskManager.HeartRateZones.any

	var body: some View {
		VStack(alignment: .leading) {
			Text("Running task")
				.font(.headline)
				.foregroundColor(.ruplBlue)
				.padding(.bottom)

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
			}.frame(height: 80)

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
