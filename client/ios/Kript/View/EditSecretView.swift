//
//  EditDatumView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/12/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation
import SwiftUI

struct EditSecretView: View {
    @State var secret: Kript_Api_Secret = Kript_Api_Secret()
    @Binding var presenting: Bool
    @State private var loading: Bool = false
    let completion: (Kript_Api_Secret) -> ()
    var deleteCallback: (() -> ())?
    @State var selectedTypeIndex: Int
    
    private let types: [(String, Kript_Api_Secret.OneOf_Secret)] = [
        ("Password", .password(Kript_Api_Secret.Password())),
        ("Credit Card", .creditCard(Kript_Api_Secret.CreditCard())),
        ("Note", .note(Kript_Api_Secret.Note())),
        ("Code", .code(Kript_Api_Secret.Code())),
    ]
    
    init(secret: Kript_Api_Secret = Kript_Api_Secret(),
         presenting: Binding<Bool>,
         completion: @escaping (Kript_Api_Secret) -> (),
         deleteCallback: (() -> ())? = nil) {
        self._secret = State(initialValue: secret)
        self._presenting = presenting
        self.completion = completion
        self.deleteCallback = deleteCallback
        
        switch secret.secret {
        case .password(_):
            self._selectedTypeIndex = State(initialValue: 0)
        case .creditCard(_):
            self._selectedTypeIndex = State(initialValue: 1)
        case .note(_):
            self._selectedTypeIndex = State(initialValue: 2)
        case .code(_):
            self._selectedTypeIndex = State(initialValue: 3)
        default:
            self._selectedTypeIndex = State(initialValue: 0)
        }
    }
    
    var body: some View {
        NavigationView {
            Form {
                Section {
                    Picker("Type", selection: $selectedTypeIndex) {
                        ForEach(0..<types.count) { i in
                            Text(self.types[i].0)
                        }
                    }
                }
                Section {
                    Group<AnyView> {
                        switch types[selectedTypeIndex].1 {
                        case .password(_):
                            return AnyView(PasswordView(password: self.$secret.password))
                        case .creditCard(_):
                            return AnyView(CreditVardView(creditCard: self.$secret.creditCard))
                        case .note(_):
                            return AnyView(NoteView(note: self.$secret.note))
                        case .code(_):
                            return AnyView(CodeView(code: self.$secret.code))
                        }
                    }
                }
                if self.deleteCallback != nil {
                    HStack {
                        Spacer()
                        Button(action: {
                            self.loading = true
                            self.deleteCallback?()
                        }, label: {
                            Text("Delete").foregroundColor(.red)
                        })
                        Spacer()
                    }
                }
            }
            .listStyle(GroupedListStyle())
            .navigationBarTitle("", displayMode: .inline)
            .navigationViewStyle(DefaultNavigationViewStyle())
            .navigationBarItems(
                leading: Button(action: {
                    self.presenting = false
                }, label: {
                    Text("Cancel").buttonStyle(DefaultButtonStyle())
                }),
                trailing: Button(action: {
                    self.loading = true
                    self.completion(self.secret)
                }, label: {
                    Text("Done").buttonStyle(PlainButtonStyle())
                })
            )
        }
    }
}


struct EditDatumView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var presenting: Bool = true
        
        var body: some View {
            NavigationView {
                EditSecretView(presenting: self.$presenting, completion: { (_) in })
            }
        }
    }
    
    static var previews: some View {
        Wrapper()
    }
}

fileprivate struct PasswordView: View {
    @Binding var password: Kript_Api_Secret.Password
    @State var hidePassword = true
    
    var body: some View {
        Group {
            IconedView(icon: "globe") {
                TextField("Website URL", text: self.$password.url)
                    .keyboardType(.URL)
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
            }
            IconedView(icon: "person") {
                TextField("Username", text: self.$password.username)
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
            }
            IconedView(icon: "lock") {
                HStack {
                    SecureTextField(label: "Password", text: self.$password.password, hidden: self.$hidePassword)
                    Spacer()
                    Button(action: {
                        self.password.password = self.generatePassword()
                        self.hidePassword = false
                    }, label: {
                        Image(systemName: "goforward.plus")
                    }).buttonStyle(BorderlessButtonStyle())
                }
            }
        }.navigationBarTitle("", displayMode: .inline)
    }
    
    private func generatePassword() -> String {
        let letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
        return String((0..<16).map{ _ in letters.randomElement()! })
    }
}

struct EditPasswordView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var secret: Kript_Api_Secret.Password = try! Kript_Api_Secret.Password(jsonUTF8Data: JSONEncoder().encode([
            "url": "google.com",
            "username": "liam923",
            "password": "secure123"
        ]))
        
        var body: some View {
            NavigationView {
                Form {
                    Section {
                        Group<AnyView> {
                            AnyView(PasswordView(password: $secret))
                        }
                    }
                }
            }
        }
    }
    
    static var previews: some View {
        Wrapper().listStyle(GroupedListStyle())
    }
}

fileprivate struct CreditVardView: View {
    private class DateWrapper: ObservableObject {
        var creditCard: Binding<Kript_Api_Secret.CreditCard>
        var date = Date() {
            didSet {
                self.creditCard.wrappedValue.expirationMonth = UInt32(Calendar.current.component(.month, from: date))
                self.creditCard.wrappedValue.expirationYear = UInt32(Calendar.current.component(.year, from: date))
            }
        }
        
        init(creditCard: Binding<Kript_Api_Secret.CreditCard>) {
            self.creditCard = creditCard
        }
    }
    
    @Binding var creditCard: Kript_Api_Secret.CreditCard
    @ObservedObject private var expirationDateWrapper: DateWrapper
    
    init(creditCard: Binding<Kript_Api_Secret.CreditCard>) {
        self._creditCard = creditCard
        self.expirationDateWrapper = DateWrapper(creditCard: creditCard)
    }
    
    var body: some View {
        Group {
            IconedView(icon: "doc.plaintext") {
                TextField("Description", text: self.$creditCard.description_p)
            }
            IconedView(icon: "person") {
                TextField("Name on Card", text: self.$creditCard.name)
            }
            IconedView(icon: "creditcard") {
                TextField("Number", text: self.$creditCard.number)
            }
            IconedView(icon: "calendar") {
                DatePicker(selection: self.$expirationDateWrapper.date,
                           displayedComponents: .date) {
                            EmptyView()
                }.datePickerStyle(WheelDatePickerStyle())
            }
        }.navigationBarTitle("", displayMode: .inline)
    }
}

struct EditCreditCardView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var secret: Kript_Api_Secret.CreditCard = try! Kript_Api_Secret.CreditCard(jsonUTF8Data: JSONEncoder().encode([
            "number": "1234567890123456",
            "name": "Bob Smith",
            "description": "Visa",
        ]))
        
        var body: some View {
            NavigationView {
                Form {
                    Section {
                        Group<AnyView> {
                            AnyView(CreditVardView(creditCard: $secret))
                        }
                    }
                }
            }
        }
    }
    
    static var previews: some View {
        Wrapper().listStyle(GroupedListStyle())
    }
}

fileprivate struct NoteView: View {
    @Binding var note: Kript_Api_Secret.Note
    
    var body: some View {
        Group {
            IconedView(icon: "doc.text") {
                TextField("Text", text: self.$note.text).lineLimit(nil)
            }
        }.navigationBarTitle("", displayMode: .inline)
    }
}

struct EditNoteView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var secret: Kript_Api_Secret.Note = try! Kript_Api_Secret.Note(jsonUTF8Data: JSONEncoder().encode([
            "text": "hello world"
        ]))
        
        var body: some View {
            NavigationView {
                Form {
                    Section {
                        Group<AnyView> {
                            AnyView(NoteView(note: $secret))
                        }
                    }
                }
            }
        }
    }
    
    static var previews: some View {
        Wrapper().listStyle(GroupedListStyle())
    }
}

fileprivate struct CodeView: View {
    @Binding var code: Kript_Api_Secret.Code
    @State private var hidden = true
    
    var body: some View {
        Group {
            IconedView(icon: "doc.plaintext") {
                TextField("Description", text: self.$code.description_p)
            }
            IconedView(icon: "lock") {
                SecureTextField(label: "Code",
                                text: self.$code.code,
                                hidden: self.$hidden)
            }
        }.navigationBarTitle("", displayMode: .inline)
    }
}

struct EditCodeView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var secret: Kript_Api_Secret.Code = try! Kript_Api_Secret.Code(jsonUTF8Data: JSONEncoder().encode([
            "code": "12345"
        ]))
        
        var body: some View {
            NavigationView {
                Form {
                    Section {
                        Group<AnyView> {
                            AnyView(CodeView(code: $secret))
                        }
                    }
                }
            }
        }
    }
    
    static var previews: some View {
        Wrapper().listStyle(GroupedListStyle())
    }
}

fileprivate struct SecureTextField: View {
    @State var label: String
    @Binding var text: String
    @Binding var hidden: Bool
    
    var body: some View {
        HStack {
            if self.hidden {
                SecureField(self.label, text: self.$text)
            } else {
                TextField(self.label, text: self.$text)
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
            }
            Button(action: {
                self.hidden.toggle()
            }) {
                if self.hidden {
                    Image(systemName: "eye.slash").foregroundColor(.blue)
                } else {
                    Image(systemName: "eye").foregroundColor(.blue)
                }
            }.buttonStyle(BorderlessButtonStyle())
        }
    }
}

