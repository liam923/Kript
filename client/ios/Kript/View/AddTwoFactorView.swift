//
//  AddTwoFactorView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/20/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct AddTwoFactorView: View {
    private static let types: [(String, String, Kript_Api_TwoFactorType)] = [
        ("Email", "Address", .email),
        ("Text", "Phone Number", .phoneText),
        ("Call", "Phone Number", .phoneCall),
    ]
    
    let manager: Manager
    @Binding var user: User?
    @Binding var presenting: Bool
    
    @State private var selectedTypeIndex = 0
    @State private var verificationToken: Kript_Api_VerificationToken?
    @State private var destination = ""
    @State private var code = ""
    @State private var errorMessage: ErrorMessage?
    
    var body: some View {
        Form {
            Section {
                Picker("Type", selection: $selectedTypeIndex) {
                    ForEach(0..<Self.types.count) { i in
                        Text(Self.types[i].0)
                    }
                }
                TextField(Self.types[selectedTypeIndex].1, text: self.$destination)
                Button(action: {
                    var option = Kript_Api_TwoFactor()
                    option.destination = self.destination
                    option.type = Self.types[self.selectedTypeIndex].2
                    self.manager.add(twoFactorOption: option, forUser: self.$user) { token in
                        if let token = token {
                            self.verificationToken = token
                        } else {
                            self.errorMessage = ErrorMessage(title: "An Error Occurred",
                                                             message: "Please try again.")
                        }
                    }
                }, label: {
                    Text("Send verification code")
                })
            }
            if verificationToken != nil {
                Section {
                    Text("A verification code has been sent")
                    TextField("Enter verification code", text: self.$code)
                    Button(action: {
                        self.manager.verifyTwoFactorOption(token: self.verificationToken ?? Kript_Api_VerificationToken(), code: self.code) { response in
                            switch response {
                            case .complete:
                                self.presenting = false
                            case .badCode:
                                self.errorMessage = ErrorMessage(title: "Invalid Code",
                                                                 message: "The verification code you entered was incorrect.")
                            case .otherError:
                                self.errorMessage = ErrorMessage(title: "An Error Occurred",
                                message: "Please try again.")
                            }
                        }
                    }, label: {
                        Text("Verify")
                    })
                }
            }
        }
        .navigationBarTitle("Add Two Factor Authentication")
        .alert(item: self.$errorMessage) { error in
            error.makeAlert()
        }
    }
}
