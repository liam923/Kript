//
//  Hasher.swift
//  Kript
//
//  Created by Liam Stevenson on 7/30/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation

protocol Hasher {
    func hash(password: Data, salt: Data) -> Data?
    
    func generateSalt() -> Data?
}

extension Hasher {
    func hash(password: String, salt: Data, encoding: String.Encoding = .utf8) -> Data? {
        if let password = password.data(using: encoding) {
            return self.hash(password: password, salt: salt)
        } else {
            return nil
        }
    }
}

extension Kript_Api_HashAlgorithm {
    var hasher: Hasher? {
        switch self {
        case .scrypt:
            return ScryptHasher(length: 32, N: 16384, r: 8, p: 1)
        case .pbkdf2Sha512:
            return Pbkdf2Hasher(keyCount: 32, rounds: 20000)
        case .argon2:
            return Argon2Hasher(mode: .argon2id, iterations: 75, memory: 1024, parallelism: 1, hashLength: 32)
        default:
            return nil
        }
    }
}
