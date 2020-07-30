//
//  SEncrypterTests.swift
//  KriptTests
//
//  Created by Liam Stevenson on 7/31/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import XCTest
@testable import Kript

class SEncrypterTests: XCTestCase {
    func testAesCbc() throws {
        let data = "hello word".data(using: .utf8)!
        let key = "this is a key".data(using: .utf8)!
        let iv = "initialization vector".data(using: .utf8)!
        
        let encrypter = Aes256Encrypter(mode: .cbc(padding: .pkcs7))
        let encryptedData = encrypter.encrypt(data: data, key: key, iv: iv)
        XCTAssertNotNil(encryptedData)
        let decryptedData = encrypter.decrypt(data: encryptedData!, key: key, iv: iv)
        XCTAssertEqual(decryptedData, data)
    }
    
    func testAesGcm() throws {
        let data = "hello word".data(using: .utf8)!
        let key = "this is a key".data(using: .utf8)!
        let iv = "initialization vector".data(using: .utf8)!
        
        let encrypter = Aes256Encrypter(mode: .gcm)
        let encryptedData = encrypter.encrypt(data: data, key: key, iv: iv)
        XCTAssertNotNil(encryptedData)
        let decryptedData = encrypter.decrypt(data: encryptedData!, key: key, iv: iv)
        XCTAssertEqual(decryptedData, data)
    }
}

