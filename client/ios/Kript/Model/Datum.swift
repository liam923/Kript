//
//  Datum.swift
//  Kript
//
//  Created by Liam Stevenson on 8/1/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation
import SwiftUI

struct Datum: Identifiable {
    let id: UUID
    
    var datum: Kript_Api_Datum?
    var secret: Kript_Api_Secret
    var title: String
    
    init?(fromCD datum: CDDatum, user: User) {
        guard let datum = try? datum.getProto(),
            let secret = user.decrypt(datum: datum) else {
            return nil
        }
        
        self.id = UUID()
        self.datum = datum
        self.secret = secret
        self.title = datum.title
    }
    
    init(datum: Kript_Api_Datum? = nil, secret: Kript_Api_Secret = Kript_Api_Secret(), id: UUID? = nil) {
        self.id = id ?? UUID()
        self.datum = datum
        self.secret = secret
        self.title = datum?.title ?? ""
    }
}
