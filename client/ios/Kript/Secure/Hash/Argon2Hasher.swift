//
//  Argon2Hasher.swift
//  Kript
//
//  Created by Liam Stevenson on 7/30/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation
import CatCrypto

struct Argon2Hasher: Hasher {
    let mode: CatArgon2Mode
    let iterations: Int
    let memory: Int
    let parallelism: Int
    let hashLength: Int
    
    func hash(password: Data, salt: Data) -> Data? {
        guard let passwordStr = String(data: password, encoding: .utf8), let saltStr = String(data: salt, encoding: .utf8) else { return nil }
        
        let argonCrypto = CatArgon2Crypto()
        argonCrypto.context.mode = self.mode
        argonCrypto.context.iterations = self.iterations
        argonCrypto.context.memory = self.memory
        argonCrypto.context.parallelism = self.parallelism
        argonCrypto.context.hashLength = self.hashLength
        argonCrypto.context.salt = saltStr
        
        let hash = argonCrypto.hash(password: passwordStr)
        return Data(base64Encoded: hash.base64StringValue())
    }
    
    func generateSalt() -> Data? {
        let letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
        let salt = String((0..<16).map{ _ in letters.randomElement()! })
        return salt.data(using: .utf8)!
    }
}

