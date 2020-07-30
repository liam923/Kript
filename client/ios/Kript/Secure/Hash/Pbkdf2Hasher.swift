//
//  Pbkdf2Hasher.swift
//  Kript
//
//  Created by Liam Stevenson on 7/30/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation
import CommonCrypto

struct Pbkdf2Hasher: Hasher {
    let keyCount: Int
    let rounds: UInt32
    
    func hash(password: Data, salt: Data) -> Data? {
        return pbkdf2(hash: CCPBKDFAlgorithm(kCCPRFHmacAlgSHA512),
                      password: password,
                      salt: salt,
                      keyCount: self.keyCount,
                      rounds: self.rounds)
    }
    
    private func pbkdf2(hash: CCPBKDFAlgorithm, password: Data, salt: Data, keyCount: Int, rounds: UInt32) -> Data? {
        var derivedKey   = [UInt8](repeating: 0, count:keyCount)
        let passwordStr = String(data: password, encoding: .utf8)
        let saltBytes = [UInt8](salt)
        

        let derivationStatus = saltBytes.withUnsafeBufferPointer { saltBytes in
            derivedKey.withUnsafeMutableBufferPointer { derivedKey in
                CCKeyDerivationPBKDF(
                CCPBKDFAlgorithm(kCCPBKDF2),
                passwordStr, password.count,
                saltBytes.baseAddress, salt.count,
                CCPseudoRandomAlgorithm(hash),
                rounds,
                derivedKey.baseAddress,
                derivedKey.count)
            }
        }

        if (derivationStatus != 0) {
            return nil;
        }

        return Data(derivedKey)
    }
    
    func generateSalt() -> Data? {
        return Data((0..<32).map { _ in UInt8.random(in: UInt8.min...UInt8.max) })
    }
}
