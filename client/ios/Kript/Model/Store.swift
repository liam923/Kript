//
//  Store.swift
//  Kript
//
//  Created by Liam Stevenson on 8/1/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import SwiftUI

struct Store {
    var datums: [Datum]
    
    init(datums: [Datum]) {
        self.datums = datums
    }
}
