//
//  LoadingIndicatorView.swift
//  rupl
//
//  Created by Dmitry Novikov on 21/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct LoadingIndicatorView: View {
	var wheelSize: CGFloat = 100.0
	var wheelStartPosition: Double = 0.0
	var wheelAnimation: Bool = true

	@State private var wheelRotating: Double = 0.0
	@State private var wheelOpacity: Double = 0.2

	var body: some View {
		let shadowShift: CGFloat = wheelSize / 10
		let shadowRadius: CGFloat = shadowShift / 10
		ZStack {
			Circle()
				.shadow(color: .ruplBlue.opacity(0.8), radius: shadowRadius, x: 0, y: -shadowShift)
				.shadow(color: .ruplYellow.opacity(0.8), radius: shadowRadius, x: 0, y: shadowShift)
				.shadow(color: .ruplRed.opacity(0.8), radius: shadowRadius, x: shadowShift, y: 0)
				.shadow(color: .ruplGreen.opacity(0.8), radius: shadowRadius, x: -shadowShift, y: 0)
				.rotationEffect(.degrees(wheelAnimation ? wheelRotating : wheelStartPosition))
				.opacity(wheelOpacity)
			Circle()
				.fill(.ruplBackground)
				.shadow(color: .ruplBackground, radius: 2)
		}
		.frame(width: wheelSize, height: wheelSize)
		.padding(10)
		.foregroundStyle(
			.shadow(.inner(color: .ruplBackground, radius: 2))
		)
		.clipShape(Circle())

		.onAppear {
			if wheelAnimation {
				withAnimation(
					.linear(duration: 1)
					.speed(0.1)
					.repeatForever(autoreverses: false)
				) {
					wheelRotating = 360.0
				}
				withAnimation(
					.easeInOut(duration: 0.8)
					.speed(0.2)
					.repeatForever(autoreverses: true)
				) {
					wheelOpacity = 1.0
				}
			}
		}
	}
}

//struct SpinningWheelView_Previews: PreviewProvider {
//	static var previews: some View {
//		LoadingIndicatorView()
//	}
//}
