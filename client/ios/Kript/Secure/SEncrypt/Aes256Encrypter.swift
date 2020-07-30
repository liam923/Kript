//
//  AesEncrypter.swift
//  Kript
//
//  Created by Liam Stevenson on 7/31/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation
import CryptoSwift

struct Aes256Encrypter: SEncrypyer {
    enum AesMode {
        case gcm
        case cbc(padding: Padding)
    }
    
    let mode: AesMode
    
    func encrypt(data: Data, key: Data, iv: Data) -> Data? {
        if let encrypedData = try? makeAes(key: key, iv: iv).encrypt([UInt8](data)) {
            return Data(encrypedData)
        } else {
            return nil
        }
    }
    
    func decrypt(data: Data, key: Data, iv: Data) -> Data? {
        if let decryptedData = try? makeAes(key: key, iv: iv).decrypt([UInt8](data)) {
            return Data(decryptedData)
        } else {
            return nil
        }
    }
    
    private func makeAes(key: Data, iv: Data) throws -> AES {
        let adjustedKey = adjust(data: key, toSize: 32)
        let adjustedIv = adjust(data: iv, toSize: 16)
        
        var gcm: GCM?
        var cryptoSwiftBlockMode: BlockMode
        var padding = Padding.noPadding
        switch mode {
        case .cbc(let cbcPadding):
            cryptoSwiftBlockMode = CBC(iv: adjustedIv)
            padding = cbcPadding
        case .gcm:
            gcm = GCM(iv: adjustedIv, mode: .combined)
            cryptoSwiftBlockMode = gcm!
        }
        
        return try AES(key: adjustedKey, blockMode: cryptoSwiftBlockMode, padding: padding)
    }
    
    private func adjust(data: Data, toSize size: Int) -> [UInt8] {
        if data.count < size {
            return data + [UInt8](repeating: 0, count: size - data.count)
        } else {
            return [UInt8](data[0..<size])
        }
    }
}
