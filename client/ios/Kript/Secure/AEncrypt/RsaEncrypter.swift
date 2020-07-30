//
//  RsaEncrypter.swift
//  Kript
//
//  Created by Liam Stevenson on 7/31/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation
import SwiftyRSA

struct RsaEncrypter: AEncrypyer {
    let padding: SecPadding
    
    func encrypt(data: Data, publicKey: Data) -> Data? {
        guard let pemString = String(data: publicKey, encoding: .utf8) else { return nil }
        
        do {
            let key = try PublicKey(pemEncoded: pemString)
            let encryptedData = try ClearMessage(data: data).encrypted(with: key, padding: self.padding)
            return encryptedData.data
        } catch {
            return nil
        }
    }
    
    func decrypt(data: Data, privateKey: Data) -> Data? {
        guard let pemString = String(data: privateKey, encoding: .utf8) else { return nil }
        
        do {
            let key = try PrivateKey(pemEncoded: pemString)
            let decryptedData = try EncryptedMessage(data: data).decrypted(with: key, padding: self.padding)
            return decryptedData.data
        } catch {
            return nil
        }
    }
    
    func generateKeyPair() -> (publicKey: Data, privateKey: Data)? {
        do {
            let (privateKey, publicKey) = try SwiftyRSA.generateRSAKeyPair(sizeInBits: 4096)
            if let publicData = try publicKey.pemString().data(using: .utf8),
                let privateData = try privateKey.pemString().data(using: .utf8) {
                return (publicData, privateData)
            } else {
                return nil
            }
        } catch {
            return nil
        }
    }
}
