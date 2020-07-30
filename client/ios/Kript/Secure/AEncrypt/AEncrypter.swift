//
//  AEncrypter.swift
//  Kript
//
//  Created by Liam Stevenson on 7/31/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation

protocol AEncrypyer {
    func encrypt(data: Data, publicKey: Data) -> Data?
    
    func decrypt(data: Data, privateKey: Data) -> Data?
    
    func generateKeyPair() -> (publicKey: Data, privateKey: Data)?
}

extension Kript_Api_AEncryptionAlgorithm {
    var encrypter: AEncrypyer? {
        switch self {
        case .rsa:
            return RsaEncrypter(padding: .OAEP)
        default:
            return nil
        }
    }
}
