//
//  AEncrypterTests.swift
//  KriptTests
//
//  Created by Liam Stevenson on 7/31/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import XCTest
@testable import Kript

class AEncrypterTests: XCTestCase {
    func testRsa() throws {
        let encrypter = RsaEncrypter(padding: .OAEP)
        
        let data = "this is a super secure password".data(using: .utf8)!
        let (publicKey, privateKey) = encrypter.generateKeyPair()!
        
        let encryptedData = encrypter.encrypt(data: data, publicKey: publicKey)
        XCTAssertNotNil(encryptedData)
        let decryptedData = encrypter.decrypt(data: encryptedData!, privateKey: privateKey)
        XCTAssertEqual(decryptedData?.base64EncodedString(), data.base64EncodedString())
    }
}
