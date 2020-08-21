//
//  SettingsView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/20/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct SettingsView: View {
    let manager: Manager
    @Binding var user: User?
    
    @State var showingAddTwoFactor = false
    @State var showingChangePassword = false
    
    var body: some View {
        Form {
            Section {
                NavigationLink(destination: AddTwoFactorView(manager: self.manager, user: self.$user, presenting: self.$showingAddTwoFactor), isActive: self.$showingAddTwoFactor) {
                    Text("Add two factor authentication option")
                }
                NavigationLink(destination: ChangePasswordView(manager: self.manager, user: self.$user, presenting: self.$showingChangePassword), isActive: self.$showingChangePassword) {
                    Text("Change password")
                }
            }
            Section {
                Button(action: {
                    self.manager.logoutUser()
                    self.user = nil
                }, label: {
                    Text("Logout")
                })
            }
        }
        .navigationBarTitle("Settings")
    }
}
