//
//  DatumView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/1/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct DatumView: View {
    @Binding var datum: Datum
    let manager: Manager
    @Binding var user: User?
    @State var editing: Bool = false
    let deleteCallback: () -> ()
    let saveCallback: () -> ()
    
    var body: some View {
        var view: AnyView
        if case .password = datum.secret.secret {
            view = AnyView(PasswordView(password: $datum.secret.password))
        } else if case .creditCard = datum.secret.secret {
            view = AnyView(CreditVardView(creditCard: $datum.secret.creditCard))
        } else if case .note = datum.secret.secret {
            view = AnyView(NoteView(note: $datum.secret.note))
        } else if case .code = datum.secret.secret {
            view = AnyView(CodeView(code: $datum.secret.code))
        } else {
            view = AnyView(Text(""))
        }
        
        return view
            .navigationBarTitle("", displayMode: .inline)
            .navigationBarItems(trailing:
                Button(action: {
                    self.editing = true
                }) {
                    Image(systemName: "square.and.pencil")
            })
            .sheet(isPresented: $editing) {
                EditSecretView(secret: self.datum.secret, presenting: self.$editing, completion: { secret in
                    var newDatum = self.datum
                    newDatum.secret = secret
                    self.manager.update(datum: newDatum, forUser: self.$user) { kDatum in
                        if let kDatum = kDatum {
                            self.datum.secret = secret
                            self.datum.datum = kDatum
                            self.saveCallback()
                        }
                        self.editing = false
                    }
                }, deleteCallback: {
                    self.manager.remove(datum: self.datum, forUser: self.$user) { kDatum in
                        self.editing = false
                        self.deleteCallback()
                    }
                })
        }
    }
}

fileprivate struct PasswordView: View {
    @Binding var password: Kript_Api_Secret.Password
    
    var body: some View {
        List {
            CopiableText(text: $password.url, icon: "globe")
            CopiableText(text: $password.username, icon: "person")
            SecureText(text: $password.password, icon: "lock")
        }.navigationBarTitle("", displayMode: .inline)
    }
}

struct PasswordView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var secret: Kript_Api_Secret.Password = try! Kript_Api_Secret.Password(jsonUTF8Data: JSONEncoder().encode([
            "url": "google.com",
            "username": "liam923",
            "password": "secure123"
        ]))
        
        var body: some View {
            PasswordView(password: $secret)
        }
    }
    
    static var previews: some View {
        Wrapper().listStyle(GroupedListStyle())
    }
}

fileprivate struct CreditVardView: View {
    @Binding var creditCard: Kript_Api_Secret.CreditCard
    
    var body: some View {
        List {
            Section {
                CopiableText(text: $creditCard.description_p, icon: "doc.plaintext")
                CopiableText(text: $creditCard.name, icon: "person")
            }
            Section {
                SecureText(text: $creditCard.number, icon: "creditcard", exceptLast: 4)
                CopiableText(text: Binding<String>(get: {
                    "Exp: \(String(format: "%02d", self.creditCard.expirationMonth))/\(String(format: "%02d",self.creditCard.expirationYear))"
                }, set: { _ in }), icon: "calendar")
            }
        }
        .navigationBarTitle("", displayMode: .inline)
    }
}

struct CreditCardView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var secret: Kript_Api_Secret.CreditCard = try! Kript_Api_Secret.CreditCard(jsonUTF8Data: JSONEncoder().encode([
            "number": "1234567890123456",
            "name": "Bob Smith",
            "description": "Visa",
        ]))
        
        var body: some View {
            CreditVardView(creditCard: $secret)
        }
    }
    
    static var previews: some View {
        Wrapper().listStyle(GroupedListStyle())
    }
}

fileprivate struct NoteView: View {
    @Binding var note: Kript_Api_Secret.Note
    
    var body: some View {
        List {
            CopiableText(text: $note.text, icon: "doc.text")
        }.navigationBarTitle("", displayMode: .inline)
    }
}

struct NoteView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var secret: Kript_Api_Secret.Note = try! Kript_Api_Secret.Note(jsonUTF8Data: JSONEncoder().encode([
            "text": "hello world\nhey\nhi again"
        ]))
        
        var body: some View {
            NoteView(note: $secret)
        }
    }
    
    static var previews: some View {
        Wrapper().listStyle(GroupedListStyle())
    }
}

fileprivate struct CodeView: View {
    @Binding var code: Kript_Api_Secret.Code
    
    var body: some View {
        List {
            CopiableText(text: $code.description_p, icon: "doc.plaintext")
            SecureText(text: $code.code, icon: "lock")
        }.navigationBarTitle("", displayMode: .inline)
    }
}

struct CodeView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var secret: Kript_Api_Secret.Code = try! Kript_Api_Secret.Code(jsonUTF8Data: JSONEncoder().encode([
            "code": "12345"
        ]))
        
        var body: some View {
            CodeView(code: $secret)
        }
    }
    
    static var previews: some View {
        Wrapper().listStyle(GroupedListStyle())
    }
}

fileprivate struct CopiableText: View {
    @Binding var text: String
    @State var icon: String?
    let textWrapper: ((String) -> String) = { $0 }
    
    var body: some View {
        IconedView(icon: self.icon) {
            HStack {
                Text(self.textWrapper(self.text))
                Spacer()
                Button(action: {
                    UIPasteboard.general.string = self.text
                }) {
                    Image(systemName: "doc.on.doc").foregroundColor(.blue)
                }.buttonStyle(BorderlessButtonStyle())
            }
        }
    }
}

fileprivate struct SecureText: View {
    @Binding var text: String
    @State var icon: String?
    @State var hidden: Bool = true
    @State var exceptLast: Int = 0
    var textWrapper: ((String) -> String) = { $0 }
    
    var body: some View {
        IconedView(icon: self.icon) {
            HStack {
                Text(self.textWrapper(self.hidden ? self.obscure(text: self.text) : self.text))
                Spacer()
                Button(action: {
                    self.hidden.toggle()
                }) {
                    if self.hidden {
                        Image(systemName: "eye.slash").foregroundColor(.blue)
                    } else {
                        Image(systemName: "eye").foregroundColor(.blue)
                    }
                }.buttonStyle(BorderlessButtonStyle())
                Button(action: {
                    UIPasteboard.general.string = self.text
                }) {
                    Image(systemName: "doc.on.doc").foregroundColor(.blue)
                }.buttonStyle(BorderlessButtonStyle())
            }
        }
    }
    
    private func obscure(text: String) -> String {
        return String((0..<text.count).map { $0 < text.count - exceptLast ? "•" : text[text.index(text.startIndex, offsetBy: $0)] })
    }
}
