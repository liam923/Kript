//
//  ErrorMessage.swift
//  Kript
//
//  Created by Liam Stevenson on 8/21/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct ErrorMessage: Identifiable {
    let id = UUID()
    let title: String
    let message: String?
    
    func makeAlert() -> Alert {
        if let messageText = message {
            return Alert(title: Text(title), message: Text(messageText), dismissButton: .none)
        } else {
            return Alert(title: Text(title), message: nil, dismissButton: .none)
        }
    }
}
