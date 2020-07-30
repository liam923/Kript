//
//  CreateAccountView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/14/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct CreateAccountView: View {
    private struct Message: Identifiable {
        let id = UUID()
        let title: String
        let message: String?
    }
    
    let manager: Manager
    @Binding var user: User?
    
    @State private var username = ""
    @State private var password = ""
    @State private var confirmedPassword = ""
    @State private var failedLoginMessage: Message?
    
    var body: some View {
        Form {
            IconedView(icon: "person") {
                TextField("Username", text: self.$username)
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
            }
            IconedView(icon: "lock") {
                SecureField("Password", text: self.$password)
            }
            IconedView(icon: "lock") {
                SecureField("Confirm Password", text: self.$confirmedPassword)
            }
            Button(action: {
                if self.password == self.confirmedPassword {
                    self.manager.createAccount(username: self.username, password: self.password) { response in
                        switch response {
                        case .complete(user: let user):
                            self.user = user
                        case .usernameTaken:
                            self.failedLoginMessage = Message(title: "Bad Username",
                                                              message: "That username has already been taken.")
                        default:
                            self.failedLoginMessage = Message(title: "Error",
                                                              message: "An error occurred. Please try again.")
                        }
                    }
                } else {
                    self.failedLoginMessage = Message(title: "Passwords Don't Match",
                    message: "Your password and confirmed password don't match. Please try again.")
                }
            }, label: {
                HStack {
                    Spacer()
                    Text("Login")
                    Spacer()
                }
            })
        }
        .alert(item: $failedLoginMessage) { message in
            if let messageText = message.message {
                return Alert(title: Text(message.title), message: Text(messageText), dismissButton: .none)
            } else {
                return Alert(title: Text(message.title), message: nil, dismissButton: .none)
            }
        }
    }
}
