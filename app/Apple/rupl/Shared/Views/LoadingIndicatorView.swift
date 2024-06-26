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
	var indicatorColor: Color = .ruplBlue

	@State private var wheelRotating: Double = 0.0
//	@State var circleColor: Color = .ruplBlue
//	@State var indicatorColor: Color = .ruplBlue
	
	var body: some View {
		
		ZStack {
//			Circle()
//				.stroke(circleColor.opacity(0.25), lineWidth: 2)
			Circle()
				.trim(from: /*@START_MENU_TOKEN@*/0.0/*@END_MENU_TOKEN@*/, to: 0.75)
				.stroke(indicatorColor, lineWidth: 4)
				.rotationEffect(.degrees(wheelRotating))
		}
		.frame(width: wheelSize, height: wheelSize , alignment: .center)
		
		.onAppear {
//			circleColor = getColor()
//			indicatorColor = getColor()

			withAnimation(
				.linear(duration: 1)
				.speed(0.25)
				.repeatForever(autoreverses: false)
			) {
				wheelRotating = 360.0
			}
		}
	}
	
//	private func getColor() -> Color {
//		return ([.ruplBlue, .ruplRed, .ruplYellow, .ruplGreen, .ruplGray].randomElement() ?? .ruplBlue)
//	}
}

#Preview {
	LoadingIndicatorView()
}
