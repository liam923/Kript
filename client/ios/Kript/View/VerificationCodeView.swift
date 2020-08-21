//
//  SendVerificationCodeView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/20/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct VerificationCodeView: View {
    let manager: Manager
    @Binding var user: User?
    @Binding var password: String
    @Binding var verificationToken: Kript_Api_VerificationToken
    @Binding var options: [String: Kript_Api_TwoFactor]
    @Binding var presenting: Bool
    @Binding var failedLoginMessage: ErrorMessage?
    @State var selectedOption: String?
    @State private var code = ""
    
    var body: some View {
        NavigationView {
            Form {
                Section {
                    Text("Choose a two-factor authentication option:")
                }
                Section(header: Text("Options")) {
                    ForEach(options.keys.sorted(), id: \.self) { id in
                        Button(action: {
                            self.manager.sendVerificationCode(token: self.verificationToken, optionId: id) { success in
                                if success {
                                    self.selectedOption = id
                                }
                            }
                        }, label: {
                            HStack {
                                Image(systemName: self.getIconName(type: self.options[id]?.type))
                                Text(self.options[id]?.destination ?? "")
                            }
                        })
                    }
                }
                if selectedOption != nil {
                    Section(header: Text("Enter Code")) {
                        Text("Enter the code sent to \(self.options[selectedOption ?? ""]?.destination ?? ""):")
                        TextField("Code", text: self.$code)
                        Button(action: {
                            self.manager.verifyUser(token: self.verificationToken, code: self.code, password: self.password) { response in
                                self.presenting = false
                                switch response {
                                case .complete(user: let user):
                                    self.user = user
                                case .badCode:
                                    self.failedLoginMessage = ErrorMessage(title: "Verification Failed",
                                                                           message: "Invalid verification code entered")
                                default:
                                    self.failedLoginMessage = ErrorMessage(title: "Verification Failed",
                                                                           message: "Please try again.")
                                }
                            }
                        }, label: {
                            Text("Enter Code")
                        })
                    }
                }
            }
            .navigationBarTitle("Verify Identity", displayMode: .inline)
            .navigationBarItems(leading:
                Button(action: {
                    self.presenting = false
                }, label: {
                    Text("Cancel")
                })
            )
        }
    }
    
    func getIconName(type: Kript_Api_TwoFactorType?) -> String {
        switch type {
        case .phoneCall:
            return "phone"
        case .phoneText:
            return "bubble.right"
        case .email:
            return "envelope"
        default:
            return ""
        }
    }
}
