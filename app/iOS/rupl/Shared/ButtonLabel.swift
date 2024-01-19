//
//  ButtonLabel.swift
//  rupl
//
//  Created by Dmitry Novikov on 05/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

//	A menu button label that aligns its icon and title
struct ButtonLabel: View {
	struct WatchMenuLabelStyle: LabelStyle {
		func makeBody(configuration: Configuration) -> some View {
			HStack {
				configuration.icon
					.frame(width: 30)
				configuration.title
				Spacer()
			}
		}
	}

	struct VerticalIconTitleLabelStyle: LabelStyle {
		func makeBody(configuration: Configuration) -> some View {
			VStack {
				Spacer()
				configuration.icon
					.frame(width: 30)
				configuration.title
					.frame(minWidth: 60)
				Spacer()
			}
		}
	}

	let title: String
	let systemImage: String

	var body: some View {
		#if os(watchOS)
			Label(title, systemImage: systemImage)
				.labelStyle(WatchMenuLabelStyle())
		#else
			Label(title, systemImage: systemImage)
				.labelStyle(VerticalIconTitleLabelStyle())
		#endif
	}
}
