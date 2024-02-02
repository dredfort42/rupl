//
//  OAuthInstructionView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 02/02/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct OAuthInstructionView: View {
	@State private var opacityAnimation: Double = 0.0

	let url: String
	let code: String

	var body: some View {
		VStack {
			Text("Using a browser on another device, visit:")
				.font(.footnote)
				.fontWeight(.light)
				.foregroundColor(.ruplGray)

			Text(url)
				.padding(.vertical)
				.font(.caption2)

			Divider()
				.background(getColor())
				.frame(height: 1)
				.opacity(opacityAnimation)

			Text("And enter the code:")
				.padding(.top)
				.font(.footnote)
				.fontWeight(.light)
				.foregroundColor(.ruplGray)

			Text(code)
				.font(.title2)
				.foregroundColor(.ruplBlue)
		}
		.onAppear {
			withAnimation(
				.linear(duration: 1)
				.speed(0.5)
				.repeatForever(autoreverses: true)
			) {opacityAnimation = 1.0}
		}
	}

	private func getColor() -> Color {
		return ([.ruplBlue, .ruplRed, .ruplYellow, .ruplGreen, .ruplGray].randomElement() ?? .ruplBlue)
	}
}

#Preview {
	OAuthInstructionView(url: "https://rupl.org/device", code: "WDJB-MJHT")
}
