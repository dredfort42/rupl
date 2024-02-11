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

	let appVersion = Bundle.main.infoDictionary?["CFBundleShortVersionString"] as? String ?? "N/A"
	let buildNumber = Bundle.main.infoDictionary?["CFBundleVersion"] as? String ?? "N/A"

	@State private var deviceAuthorization = false
	@State private var polling: Bool = false
	@State private var userCode: String = ""
	@State private var verificationUri: String = ""

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

					Section {
					} footer: {
						Text("Version: " + appVersion + "." + buildNumber)
							.foregroundColor(.ruplBlue)
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

//#Preview {
//	SettingsView()
//}
