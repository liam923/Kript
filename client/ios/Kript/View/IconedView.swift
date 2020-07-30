//
//  IconedView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/14/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct IconedView<Content: View>: View {
    @State var icon: String?
    let content: () -> Content
    
    var body: some View {
        HStack {
            if icon != nil { Image(systemName: icon ?? "").foregroundColor(.gray) }
            content()
        }
    }
}
