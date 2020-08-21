//
//  DatumListView.swift
//  Kript
//
//  Created by Liam Stevenson on 8/1/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import GRPC
import KeychainSwift
import NIO
import SwiftUI
import SwiftUIRefresh

struct DatumListView: View {
    let manager: Manager
    @Binding var user: User?
    @State var store: Store
    @State private var refreshing: Bool = false
    @State private var creating: Bool = false
    @State var datumIsActive: [Bool] = [Bool](repeating: false, count: 100)
    
    init(manager: Manager, user: Binding<User?>, store: Store? = nil) {
        self.manager = manager
        self._user = user
        if let store = store {
            self._store = State(initialValue: store)
        } else if let user = user.wrappedValue {
            self._store = State(initialValue: manager.loadStoreFromCoreData(user: user))
        } else {
            self._store = State(initialValue: Store(datums: []))
        }
    }
    
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Passwords")) {
                    ForEach(store.datums.identifiableIndices.filter({ i in
                        if case .password = self.store.datums[i.val].secret.secret {
                            return true
                        } else {
                            return false
                        }
                    })) { i in
                        self.makeDatumCell(i: i.val)
                    }
                }
                Section(header: Text("Credit Cards")) {
                    ForEach(store.datums.identifiableIndices.filter({ i in
                        if case .creditCard = self.store.datums[i.val].secret.secret {
                            return true
                        } else {
                            return false
                        }
                    })) { i in
                        self.makeDatumCell(i: i.val)
                    }
                }
                Section(header: Text("Notes")) {
                    ForEach(store.datums.identifiableIndices.filter({ i in
                        if case .note = self.store.datums[i.val].secret.secret {
                            return true
                        } else {
                            return false
                        }
                    })) { i in
                        self.makeDatumCell(i: i.val)
                    }
                }
                Section(header: Text("Codes")) {
                    ForEach(store.datums.identifiableIndices.filter({ i in
                        if case .code = self.store.datums[i.val].secret.secret {
                            return true
                        } else {
                            return false
                        }
                    })) { i in
                        self.makeDatumCell(i: i.val)
                    }
                }
            }
            .listStyle(GroupedListStyle())
            .pullToRefresh(isShowing: $refreshing) {
                self.refresh()
            }
            .navigationBarTitle("Your Data")
            .navigationBarItems(
                leading: NavigationLink(destination: SettingsView(manager: manager, user: self.$user), label: {
                    Text("Settings")
                }), trailing: Button(action: {
                    self.creating = true
                }, label: {
                    Image(systemName: "plus")
                })
            )
        }
        .sheet(isPresented: $creating) {
            EditSecretView(presenting: self.$creating, completion: { secret in
                self.manager.add(datum: Datum(secret: secret), forUser: self.$user) { kDatum in
                    if let kDatum = kDatum {
                        if self.datumIsActive.count < self.store.datums.count + 1 {
                            self.datumIsActive.append(false)
                        }
                        self.store.datums.append(Datum(datum: kDatum, secret: secret))
                        self.saveStore()
                    }
                    self.creating = false
                }
            })
        }
        .onAppear {
            self.refresh()
        }
    }
    
    private func makeDatumCell(i: Int) -> DatumCellView {
        return DatumCellView(datum: self.$store.datums[i],
                             manager: self.manager,
                             user: self.$user,
                             deleteCallback: {
                                self.manager.remove(datum: self.store.datums[i], forUser: self.$user) { kDatum in
                                    if let _ = kDatum {
                                        self.datumIsActive[i] = false
                                        self.store.datums.remove(at: i)
                                        self.saveStore()
                                        self.creating = false
                                    }
                                }
        },
                             saveCallback: self.saveStore,
                             isActive: self.$datumIsActive[i])
    }
    
    private func addDatumCompletion(secret: Kript_Api_Secret, completion: @escaping () -> ()) {
        self.manager.add(datum: Datum(secret: Kript_Api_Secret()), forUser: self.$user) { kDatum in
            if let kDatum = kDatum {
                self.store.datums.append(Datum(datum: kDatum, secret: secret))
            }
            completion()
        }
    }
    
    private func refresh() {
        self.manager.refresh(store: $store, forUser: $user) { _ in
            self.refreshing = false
            self.saveStore()
        }
    }
    
    private func sortDatums(a: Array<Datum>.IdentifiableInt, b: Array<Datum>.IdentifiableInt) -> Bool {
        let datumA = self.store.datums[a.val]
        let datumB = self.store.datums[b.val]
        
        switch datumA.secret.secret {
        case .password(let passwordA):
            switch datumB.secret.secret {
            case .password(let passwordB):
                return passwordA.url < passwordB.url
            default:
                return false
            }
        default:
            return false
        }
    }
    
    private func saveStore() {
        self.manager.saveStoreToCoreData(store: self.store)
    }
}

extension Array where Element: Identifiable {
    struct IdentifiableInt: Identifiable {
        var id: Element.ID
        var val: Int
    }
    
    var identifiableIndices: [IdentifiableInt] {
        return self.indices.map { IdentifiableInt(id: self[$0].id, val: $0) }
    }
}

fileprivate struct DatumCellView: View {
    @Binding var datum: Datum
    let manager: Manager
    @Binding var user: User?
    let deleteCallback: () -> ()
    let saveCallback: () -> ()
    @Binding var isActive: Bool
    
    var body: some View {
        NavigationLink(destination: DatumView(datum: $datum,
                                              manager: self.manager,
                                              user: $user,
                                              deleteCallback: deleteCallback,
                                              saveCallback: saveCallback),
                       isActive: $isActive) {
                        HStack {
                            Text(getPresentationText(forSecret: datum.secret))
                            Spacer()
                            Button(action: {
                                UIPasteboard.general.string = self.getCopyText(forSecret: self.datum.secret)
                            }) {
                                Image(systemName: "doc.on.doc").foregroundColor(.blue)
                            }.buttonStyle(BorderlessButtonStyle())
                        }
        }
    }
    
    func getPresentationText(forSecret secret: Kript_Api_Secret) -> String {
        switch secret.secret {
        case .password(let password):
            return password.url
        case .creditCard(let creditCard):
            return creditCard.description_p
        case .note(let note):
            return note.text
        case .code(let code):
            return code.description_p
        default:
            return ""
        }
    }
    
    func getCopyText(forSecret secret: Kript_Api_Secret) -> String {
        switch secret.secret {
        case .password(let password):
            return password.password
        case .creditCard(let creditCard):
            return creditCard.number
        case .note(let note):
            return note.text
        case .code(let code):
            return code.code
        default:
            return ""
        }
    }
}

struct DatumListView_Previews: PreviewProvider {
    private struct Wrapper: View {
        @State var user: User?
        
        var body: some View {
            var datums = [Datum(), Datum(), Datum(), Datum()]
            datums[0].secret.password = Kript_Api_Secret.Password()
            datums[0].secret.password.url = "google.com"
            datums[0].secret.password.password = "secret123"
            
            datums[1].secret.creditCard = Kript_Api_Secret.CreditCard()
            datums[1].secret.creditCard.description_p = "Visa"
            
            datums[2].secret.note = Kript_Api_Secret.Note()
            datums[2].secret.note.text = "birthday is 9/23/99"
            
            datums[3].secret.code = Kript_Api_Secret.Code()
            datums[3].secret.code.description_p = "Bank Number"
            let store = Store(datums: datums)
            return DatumListView(manager: MockManager(), user: $user, store: store)
        }
    }
    
    static var previews: some View {
        Wrapper()
    }
}
