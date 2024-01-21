//
//  ContentView.swift
//  rupl
//
//  Created by Dmitry Novikov on 04/01/2024.
//

import SwiftUI

struct ContentView: View {
    var body: some View {
        VStack {
            Image(systemName: "pencil.and.outline")
				.imageScale(.large)
            Text("rupl")
				.font(.system(size: 50, weight: .medium, design: .rounded))
				.foregroundStyle(.ruplBlue)

        }
        .padding()
    }
}

//#Preview {
//    ContentView()
//}
