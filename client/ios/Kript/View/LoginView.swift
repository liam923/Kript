//
//  LoginView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/14/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct LoginView: View {
    let manager: Manager
    @Binding var user: User?
    
    @State private var username = ""
    @State private var password = ""
    @State private var failedLoginMessage: ErrorMessage?
    
    @State private var showingVerification = false
    @State private var verificationToken = Kript_Api_VerificationToken()
    @State private var verificationOptions = [String: Kript_Api_TwoFactor]()
    
    var body: some View {
        NavigationView {
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
                            self.failedLoginMessage = ErrorMessage(title: "Bad Username",
                                                                   message: "That username doesn't seem to exist.")
                        case .badPassword:
                            self.failedLoginMessage = ErrorMessage(title: "Wrong Password",
                                                                   message: "That password was incorrect.")
                        case .twoFactor(verificationToken: let token, options: let options):
                            self.verificationToken = token
                            self.verificationOptions = options
                            self.showingVerification = true
                        default:
                            self.failedLoginMessage = ErrorMessage(title: "Error",
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
                message.makeAlert()
            }
            .sheet(isPresented: self.$showingVerification) {
                VerificationCodeView(manager: self.manager, user: self.$user, password: self.$password, verificationToken: self.$verificationToken, options: self.$verificationOptions, presenting: self.$showingVerification, failedLoginMessage: self.$failedLoginMessage)
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
