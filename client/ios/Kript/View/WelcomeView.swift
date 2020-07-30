//
//  WelcomeView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/14/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct WelcomeView: View {
    let manager: Manager
    @Binding var user: User?
    
    var body: some View {
        NavigationView {
            VStack {
                NavigationLink(destination: LoginView(manager: manager, user: self.$user)) {
                    Text("Login")
                }
                NavigationLink(destination: CreateAccountView(manager: manager, user: self.$user)) {
                    Text("Create Account")
                }
            }
        }
    }
}
