//
//  SlideButtonView.swift
//  rupl
//
//  Created by Dmitry Novikov on 20/04/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

public struct SlideButton<Label: View>: View {
	@Environment(\.isEnabled) private var isEnabled
	@Environment(\.layoutDirection) private var layoutDirection

	@GestureState private var offset: CGFloat
	@State private var swipeState: SwipeState = .start

	private let title: Label
	private let callback: () async -> Void
	private let styling: SlideButtonStyle

	public init(styling: SlideButtonStyle = .default, action: @escaping () async -> Void, @ViewBuilder label: () -> Label) {
		self.title = label()
		self.callback = action
		self.styling = styling

		self._offset = .init(initialValue: styling.indicatorSpacing)
	}

	public var body: some View {
		GeometryReader { reading in
			let calculatedOffset: CGFloat = swipeState == .swiping ? offset : (swipeState == .start ? styling.indicatorSpacing : (reading.size.width - styling.indicatorSize + styling.indicatorSpacing))

			ZStack(alignment: .leading) {
				styling.color
					.saturation(0.5)
					.opacity(0.5)

				ZStack {

					title
						.multilineTextAlignment(.center)
						.foregroundColor(styling.color)
						.frame(maxWidth: max(0, reading.size.width - 2 * styling.indicatorSpacing), alignment: Alignment(horizontal: .center, vertical: .center))
						.padding(.trailing, styling.indicatorSpacing)
						.padding(.leading, styling.indicatorSize)
				}
				.opacity(1 - progress(from: styling.indicatorSpacing, to: reading.size.width - styling.indicatorSize + styling.indicatorSpacing, current: calculatedOffset))
				.animation(.interactiveSpring(), value: calculatedOffset)
				.mask {
					Rectangle()
						.overlay(alignment: .leading) {
							Color.ruplRed
								.frame(width: calculatedOffset + (0.5 * styling.indicatorSize - styling.indicatorSpacing))
								.frame(maxWidth: .infinity, alignment: .leading)
								.animation(.interactiveSpring(), value: swipeState)
								.blendMode(.destinationOut)
						}
				}
				.frame(maxWidth: .infinity, alignment: .trailing)

				Circle()
					.frame(width: styling.indicatorSize - 2 * styling.indicatorSpacing, height: styling.indicatorSize - 2 * styling.indicatorSpacing)
					.foregroundColor(isEnabled ? styling.color : .gray)
					.overlay(content: {
						ZStack {
							ProgressView().progressViewStyle(.circular)
								.tint(.white)
								.opacity(swipeState == .end ? 1 : 0)
							Image(systemName: styling.indicatorSystemName)
								.foregroundColor(.white)
								.font(.system(size: max(0.4 * styling.indicatorSize, 0.5 * styling.indicatorSize - 2 * styling.indicatorSpacing), weight: .semibold))
								.opacity(swipeState == .end ? 0 : 1)
						}
					})
					.offset(x: calculatedOffset)
					.animation(.interactiveSpring(), value: swipeState)
					.gesture(
						DragGesture()
							.updating($offset) { value, state, transaction in
								guard swipeState != .end else { return }

								if swipeState == .start {
									DispatchQueue.main.async {
										swipeState = .swiping
#if os(iOS)
										UIImpactFeedbackGenerator(style: .light).prepare()
#endif
									}
								}

								let val = value.translation.width

								state = clampValue(value: val, min: styling.indicatorSpacing, max: reading.size.width - styling.indicatorSize + styling.indicatorSpacing)
							}
							.onEnded { value in
								guard swipeState == .swiping else { return }
								swipeState = .end

								if value.predictedEndTranslation.width > reading.size.width
									|| value.translation.width > reading.size.width - styling.indicatorSize - 2 * styling.indicatorSpacing {
									Task {
#if os(iOS)
										UIImpactFeedbackGenerator(style: .light).impactOccurred()
#endif

										await callback()
										swipeState = .start
									}

								} else {
									swipeState = .start
#if os(iOS)
									UIImpactFeedbackGenerator(style: .light).impactOccurred()
#endif
								}
							}
					)
			}
			.mask({ Capsule() })
		}
		.frame(height: styling.indicatorSize)
	}

	private func clampValue(value: CGFloat, min minValue: CGFloat, max maxValue: CGFloat) -> CGFloat {
		return max(minValue, min(maxValue, value))
	}

	private func progress(from start: Double, to end: Double, current: Double) -> Double {
		let clampedCurrent = max(min(current, end), start)
		return (clampedCurrent - start) / (end - start)
	}

	private enum SwipeState {
		case start, swiping, end
	}
}

public extension SlideButton where Label == Text {
	init(_ titleKey: LocalizedStringKey, styling: SlideButtonStyle = .default, action: @escaping () async -> Void) {
		self.init(styling: styling, action: action, label: { Text(titleKey) })
	}
}

public struct SlideButtonStyle {
	var indicatorSize: CGFloat = 50
	var indicatorSpacing: CGFloat = 5
	var color: Color = .ruplBlue
	var indicatorSystemName: String = "chevron.right"

	public static let `default`: Self = .init()
}

//#if DEBUG
//struct SlideButton_Previews: PreviewProvider {
//	struct ContentView: View {
//		var body: some View {
//			ScrollView {
//				VStack(spacing: 25) {
//					SlideButton("Default slider", action: sliderCallback)
//					SlideButton("Decline slider", styling: .init(color: .ruplRed, indicatorSystemName: "xmark"), action: sliderCallback)
//					SlideButton("Spacing 3", styling: .init(indicatorSpacing: 3), action: sliderCallback)
//					SlideButton("Big 80", styling: .init(indicatorSize: 80), action: sliderCallback)
//					SlideButton("disabled green", styling: .init(color: .green), action: sliderCallback)
//						.disabled(true)
//					SlideButton("disabled", action: sliderCallback)
//						.disabled(true)
//				}.padding(.horizontal)
//			}
//		}
//
//		private func sliderCallback() async {
//			try? await Task.sleep(for: .seconds(2))
//		}
//	}
//
//	static var previews: some View {
//		ContentView()
//	}
//}
//#endif
