//
//  SettingsView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 28/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct SettingsView: View {
	@AppStorage(AppSettings.useAutoPauseKey) var useAutoPause = AppSettings.shared.useAutoPause
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
						if isConnectedToRupl {
							Button {
								isConnectedToRupl = !isConnectedToRupl
								resetDeviceAccess()
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
						if isConnectedToRupl {
							Text("● Connected to rupl.org")
								.foregroundColor(.ruplGreen)
						} else {
							Text("Connection to rupl.org")
						}
					} footer: {
						if isConnectedToRupl {
							Text("Disconnect from rupl.org and stop uploading running results")
						} else {
							Text("Connect to rupl.org and get training tasks to your watch")
						}
					}

					Section {
						if isConnectedToRupl && AppSettings.shared.userDateOfBirth != nil {
							Text("\(AppSettings.shared.userFirstName) \(AppSettings.shared.userLastName)\nAge: \(AppSettings.shared.getUserAge(dateOfBirth: AppSettings.shared.userDateOfBirth!)) years")
						}
					} header: {
						if isConnectedToRupl && AppSettings.shared.userDateOfBirth != nil {
							Text("Profile")
						}
					} footer: {
						Text("Version: " + AppSettings.shared.appVersion)
							.foregroundColor(.ruplBlue)
					}
				}
			}
		} else if userCode.isEmpty || verificationUri.isEmpty {
			LoadingIndicatorView()
		} else {
			OAuthInstructionView(url: verificationUri, code: userCode)
				.onAppear() {
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
		polling = true

		DispatchQueue.global().async {
			while self.polling {
				if OAuth2.expiresIn <= Date() {
					self.sendRequest()
				}

				OAuth2.getDeviceAccessToken { result in
					AppSettings.shared.deviceAccessToken = OAuth2.accessToken
					AppSettings.shared.deviceAccessTokenType = OAuth2.tokenType
					AppSettings.shared.deviceAccessTokenExpiresIn = OAuth2.expiresIn.timeIntervalSince1970

					if !OAuth2.accessToken.isEmpty && !OAuth2.tokenType.isEmpty && OAuth2.expiresIn > Date() {
						self.polling = false
						self.deviceAuthorization = false
					}
				}

				sleep(OAuth2.interval)
			}

			OAuth2.accessToken = ""
			OAuth2.tokenType = ""
			OAuth2.expiresIn = Date()
		}
	}

	func resetDeviceAccess() {
		DeviceInfo.shared.deleteDeviceInfo()
		OAuth2.deleteDeviceAccess { result in
			AppSettings.shared.deviceAccessToken = ""
			AppSettings.shared.deviceAccessTokenType = ""
			AppSettings.shared.deviceAccessTokenExpiresIn = 0
			AppSettings.shared.userEmail = ""
			AppSettings.shared.userFirstName = ""
			AppSettings.shared.userLastName = ""
			AppSettings.shared.userDateOfBirth = nil
			AppSettings.shared.userGender = ""
		}
	}
}

//#Preview {
//	SettingsView()
//}
