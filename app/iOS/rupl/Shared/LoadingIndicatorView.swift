//
//  LoadingIndicatorView.swift
//  rupl
//
//  Created by Dmitry Novikov on 21/01/2024.
//  Copyright © 2024 dredfort.42. All rights reserved.
//

import SwiftUI

struct LoadingIndicatorView: View {

	var body: some View {
			ProgressView()
				.progressViewStyle(CircularProgressViewStyle(tint: .ruplBlue))
				.scaleEffect(2, anchor: .center)
	}
}

//#Preview {
//    LoadingIndicatorView()
//}
