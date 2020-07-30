//
//  SEncrypter.swift
//  Kript
//
//  Created by Liam Stevenson on 7/31/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation

protocol SEncrypyer {
    func encrypt(data: Data, key: Data, iv: Data) -> Data?
    
    func decrypt(data: Data, key: Data, iv: Data) -> Data?
}

extension SEncrypyer {
    func generateKey() -> Data {
        return Data((0..<32).map { _ in UInt8.random(in: UInt8.min...UInt8.max) })
    }
    
    func generateIv() -> Data {
        return Data((0..<16).map { _ in UInt8.random(in: UInt8.min...UInt8.max) })
    }
}

extension Kript_Api_SEncryptionAlgorithm {
    var encrypter: SEncrypyer? {
        switch self {
        case .aes256Cbc:
            return Aes256Encrypter(mode: .cbc(padding: .pkcs7))
        case .aes256Gcm:
            return Aes256Encrypter(mode: .gcm)
        default:
            return nil
        }
    }
}
