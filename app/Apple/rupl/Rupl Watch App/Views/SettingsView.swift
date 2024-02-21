//
//  SettingsView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 28/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct SettingsView: View {
	private let appVersion = Bundle.main.infoDictionary?["CFBundleShortVersionString"] as? String ?? "N/A"
	private let buildNumber = Bundle.main.infoDictionary?["CFBundleVersion"] as? String ?? "N/A"
	private let maxUserBirthYear: Int = AppSettings.shared.getCurrentYear() - 14
	private let minUserBirthYear: Int = AppSettings.shared.getCurrentYear() - 100

	@AppStorage(AppSettings.useAutoPauseKey) var useAutoPause = AppSettings.shared.useAutoPause
	@AppStorage(AppSettings.userYearOfBirthKey) var userYearOfBirth = AppSettings.shared.userYearOfBirth
	@AppStorage(AppSettings.criticalHeartRateKey) var criticalHeartRate = AppSettings.shared.criticalHeartRate
	@AppStorage(AppSettings.connectedToRuplKey) var isConnectedToRupl = AppSettings.shared.connectedToRupl

	@State private var deviceAuthorization = false
	@State private var polling: Bool = false
	@State private var userCode: String = ""
	@State private var verificationUri: String = ""

	var body: some View {
		if !deviceAuthorization {
			VStack {
				Form {
					Section(footer: Text("Automatically pauses workout when you have paused your activity")) {
						Toggle("Auto pause", isOn: $useAutoPause)
							.toggleStyle(SwitchToggleStyle(tint: .ruplBlue))
					}

					Section {
						Stepper(value: $criticalHeartRate,
								in: 180...210,
								step: 1) {
							Text("\(criticalHeartRate)")
								.font(.title2)
						}.focusable(false)
					} header: {
						Text("Critical heart rate")
					} footer: {
						Text("If heart rate exceeds this value, an alarm will sound")
					}

					Section {
						Stepper(value: $userYearOfBirth,
								in: minUserBirthYear...maxUserBirthYear,
								step: 1) {
							Text(String(userYearOfBirth))
								.font(.title2)
						}.focusable(false)
					} header: {
						Text("Year of birth")
					} footer: {
						Text("Year of birth is required to calculate individual heart rate intervals")
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
