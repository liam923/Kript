//
//  LoginView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/14/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct LoginView: View {
    private struct Message: Identifiable {
        let id = UUID()
        let title: String
        let message: String?
    }
    
    let manager: Manager
    @Binding var user: User?
    
    @State private var username = ""
    @State private var password = ""
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
            Button(action: {
                self.manager.login(username: self.username, password: self.password) { response in
                    switch response {
                    case .complete(user: let user):
                        self.manager.saveUserDataToKeychain(user: user)
                        self.user = user
                    case .badUsername:
                        self.failedLoginMessage = Message(title: "Bad Username",
                                                          message: "That username doesn't seem to exist.")
                    case .badPassword:
                        self.failedLoginMessage = Message(title: "Wrong Password",
                                                          message: "That password was incorrect.")
                    case .twoFactor(_, _):
                        self.failedLoginMessage = Message(title: "Two Factor Required",
                                                          message: "This account requires two factor authentication, which is not yet supported on this app.")
                    default:
                        self.failedLoginMessage = Message(title: "Error",
                                                          message: "An error occurred. Please try again.")
                    }
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

struct LoginView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var user: User?
        
        var body: some View {
            LoginView(manager: MockManager(), user: self.$user)
        }
    }
    
    static var previews: some View {
        Wrapper()
    }
}
