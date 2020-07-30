//
//  ContentView.swift
//  Kript
//
//  Created by Liam Stevenson on 7/28/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct ContentView: View {
    let manager: Manager
    @State var user: User?
    
    var body: some View {
        if user == nil {
            return AnyView(WelcomeView(manager: manager, user: self.$user))
        } else {
            return AnyView(DatumListView(manager: manager, user: self.$user))
        }
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView(manager: MockManager(), user: User(id: "", refreshToken: Kript_Api_RefreshToken(), accessToken: nil, publicKey: Data(), privateKey: Data(), dataEncryptionAlgorithm: .unknownAEncryptionAlgorithm))
    }
}
