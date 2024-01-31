//
//  SettingsView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 28/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct SettingsView: View {
	@Environment(\.dismiss) var dismiss

	@AppStorage(AppSettings.useAutoPauseKey) var useAutoPauseIsOn = AppSettings.shared.useAutoPause
	@AppStorage(AppSettings.connectedToRuplKey) var isConnectedToRupl = !AppSettings.shared.connectedToRupl
	//	@State private var somethingElseIsOn = false
	//	@State private var anotherThingIsOn = false

	//	@State private var howMuch = 0.0
	//	@State private var amount = 0.0
	//	@State private var quantity = 0.0

	//	let pulsZones = [
	//		AppSettings.shared.pz1NotInZone,
	//		AppSettings.shared.pz2Easy,
	//		AppSettings.shared.pz3FatBurning,
	//		AppSettings.shared.pz4Aerobic,
	//		AppSettings.shared.pz5Anaerobic
	//	]

	var body: some View {
		
		VStack {


			Form {
				Section {
					Toggle("Auto pause", isOn: $useAutoPauseIsOn)
						.toggleStyle(SwitchToggleStyle(tint: .ruplBlue))
				} footer: {
					Text("Automatically pauses workout when you have paused your activity")
				}



				Section {
					//											Toggle("Something Else", isOn: $somethingElseIsOn)
					//

					//
					//											Slider(value: $amount, in: 0...100)
					if isConnectedToRupl {
						Button {
							isConnectedToRupl = !isConnectedToRupl
						} label: {
							HStack {
								Image(systemName: "xmark")
									.padding()
								Text("Disconnect")
							}
						}.foregroundColor(.ruplRed)
					} else {
						Button {
							isConnectedToRupl = !isConnectedToRupl
							let connect = OAuth2()
						} label: {
							HStack {
								Image(systemName: "link")
									.padding()
								Text("Connect")
							}
						}
						//							.foregroundColor(.ruplBlue)
					}
					//
				} header: {
					Text("Connection to rupl.org")
				} footer: {
					if isConnectedToRupl {
						Text("Disconnect from rupl.org and stop uploading running results")
					} else {
						Text("Connect to rupl.org and get training tasks to your watch")
					}
				}
				//
				//					Section {
				//						Toggle("Another Thing", isOn: $anotherThingIsOn)
				//
				//						Text("Quantity: \(quantity)")
				//
				//						Slider(value: $quantity, in: 0...100)

				//					} header: {
				//						Text("Even More Settings")
				//					}
				//
				//					Section {
				//						Image(systemName: "photo")
				//							.resizable()
				//							.scaledToFit()
				//
				//						Image(systemName: "photo")
				//							.resizable()
				//							.scaledToFit()
				//					} header: {
				//						Text("Photos")
				//					}

				Button {
					dismiss()
				} label: {
					Text("Done")
				}
			}


		}
	}

}

//#Preview {
//	SettingsView()
//}
