//
//  LoadingIndicatorView.swift
//  rupl Watch App
//
//  Created by Dmitry Novikov on 30/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct LoadingIndicatorView: View {
	var wheelSize: CGFloat = 100.0
	var wheelStartPosition: Double = 0.0
	var wheelAnimation: Bool = true

	@State private var wheelRotating: Double = 0.0
	var circleColor: Color = .ruplBlue
	var indicatorColor: Color = .ruplBlue

	init() {
		circleColor = getColor()
		indicatorColor = getColor()
	}

	var body: some View {

		ZStack {
			Circle()
				.stroke(circleColor.opacity(0.25), lineWidth: 2)
			Circle()
				.trim(from: /*@START_MENU_TOKEN@*/0.0/*@END_MENU_TOKEN@*/, to: 0.33)
				.stroke(indicatorColor, lineWidth: 2)
				.rotationEffect(.degrees(wheelAnimation ? wheelRotating : wheelStartPosition))
		}
		.frame(width: wheelSize, height: wheelSize , alignment: .center)

		.onAppear {
			if wheelAnimation {
				withAnimation(
					.linear(duration: 1)
					.speed(0.25)
					.repeatForever(autoreverses: false)
				) {
					wheelRotating = 360.0
				}
			}
		}
	}

	private func getColor() -> Color {
		return ([.ruplBlue, .ruplRed, .ruplYellow, .ruplGreen, .ruplGray].randomElement() ?? .ruplBlue)
	}
}

//#Preview {
//	LoadingIndicatorView()
//}
