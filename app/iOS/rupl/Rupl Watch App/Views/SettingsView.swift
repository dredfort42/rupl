//
//  SettingsView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 28/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct SettingsView: View {
//	@Environment(\.dismiss) var dismiss

	@AppStorage(AppSettings.useAutoPauseKey) var useAutoPauseIsOn = AppSettings.shared.useAutoPause
	@AppStorage(AppSettings.connectedToRuplKey) var isConnectedToRupl = AppSettings.shared.connectedToRupl


	@State private var deviceAuthorization = false
	@State private var polling: Bool = false
	@State private var userCode: String = ""
	@State private var verificationUri: String = ""

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
		if !deviceAuthorization {
			VStack {
				Form {
					Section {
						Toggle("Auto pause", isOn: $useAutoPauseIsOn)
							.toggleStyle(SwitchToggleStyle(tint: .ruplBlue))
					} footer: {
						Text("Automatically pauses workout when you have paused your activity")
					}

					Section {
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
								sendRequest()
							} label: {
								HStack {
									Image(systemName: "link")
										.padding()
									Text("Connect")
								}
							}
						}
					} header: {
						Text("Connection to rupl.org")
					} footer: {
						if isConnectedToRupl {
							Text("Disconnect from rupl.org and stop uploading running results")
						} else {
							Text("Connect to rupl.org and get training tasks to your watch")
						}
					}
				}
			}
		} else if userCode.isEmpty || verificationUri.isEmpty {
			LoadingIndicatorView()
		} else {
			OAuthInstructionView(url: verificationUri, code: userCode)
				.onAppear() {
					polling = true
					pollingResponse()
				}
				.onDisappear() {
					polling = false
				}
		}
	}

	func sendRequest() {
		deviceAuthorization = true
		OAuth2.sendRequest { result in
			userCode = OAuth2.userCode
			verificationUri = OAuth2.verificationUri
		}
	}

	func pollingResponse() {
		DispatchQueue.global().async {
			var counter: Int = 0

			while self.polling {
				print("[\(counter)] polling")
				counter += 1
				sleep(3)
			}
		}
	}
}
//
//#Preview {
//	SettingsView()
//}
