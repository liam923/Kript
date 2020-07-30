//
//  ScryptHasher.swift
//  Kript
//
//  Created by Liam Stevenson on 7/30/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation
import Scrypt

struct ScryptHasher: Hasher {
    let length: Int
    let N: UInt64
    let r: UInt32
    let p: UInt32
    
    func hash(password: Data, salt: Data) -> Data? {
        let passwordBytes = [UInt8](password)
        let saltBytes = [UInt8](salt)
        if let bytes = try? scrypt(password: passwordBytes,
                                   salt: saltBytes,
                                   length: self.length,
                                   N: self.N,
                                   r: self.r,
                                   p: self.p) {
            return Data(bytes)
        } else {
            return nil
        }
    }
    
    func generateSalt() -> Data? {
        return Data((0..<32).map { _ in UInt8.random(in: UInt8.min...UInt8.max) })
    }
}
