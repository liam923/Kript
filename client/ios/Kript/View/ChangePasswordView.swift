//
//  ChangePasswordView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/21/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct ChangePasswordView: View {
    let manager: Manager
    @Binding var user: User?
    @Binding var presenting: Bool
    
    @State private var oldPassword = ""
    @State private var newPassword = ""
    @State private var confirmPassword = ""
    @State private var errorMessage: ErrorMessage?
    
    var body: some View {
        Form {
            SecureField("Old Password", text: self.$oldPassword)
            SecureField("New Password", text: self.$newPassword)
            SecureField("Confirm New Password", text: self.$confirmPassword)
            Button(action: {
                if self.newPassword != self.confirmPassword {
                    self.errorMessage = ErrorMessage(title: "Passwords don't match",
                                                     message: "Your new passwords don't match.")
                } else {
                    self.manager.updatePassword(oldPassword: self.oldPassword, newPassword: self.newPassword, forUser: self.$user) { user in
                        if let user = user {
                            self.user?.userObject = user
                            if let user = self.user {
                                self.manager.saveUserDataToKeychain(user: user)
                            }
                            self.presenting = false
                        } else {
                            self.errorMessage = ErrorMessage(title: "An Error Occurred", message: "Please try again.")
                        }
                    }
                }
            }, label: {
                Text("Update Password")
            })
        }
        .navigationBarTitle("Change Password")
        .alert(item: self.$errorMessage) { error in
            error.makeAlert()
        }
    }
}
