//
//  HasherTests.swift
//  KriptTests
//
//  Created by Liam Stevenson on 7/30/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import XCTest
@testable import Kript

class HasherTests: XCTestCase {
    func testScrypt() throws {
        let password = "password"
        let salt = "salty salt".data(using: .utf8)!
        
        let hasher = ScryptHasher(length: 32, N: 16384, r: 8, p: 1)
        let hash = hasher.hash(password: password, salt: salt, encoding: .utf8)
        let hashBase64 = hash?.base64EncodedString()
        XCTAssertEqual(hashBase64, "zWSwcnevcbXTg+HottPf4tDTnpX6Ug8LhuCIy0Gt/Us=")
    }
    
    func testPbkdf2() throws {
        let password = "password"
        let salt = Data(base64Encoded: "5c9TSrl+adMyW+3jB7ZgJQ==")!
        
        let hasher = Pbkdf2Hasher(keyCount: 32, rounds: 20000)
        let hash = hasher.hash(password: password, salt: salt, encoding: .utf8)
        let hashBase64 = hash?.base64EncodedString()
        XCTAssertEqual(hashBase64, "V+KVaknWplhP6FRfoCeWE7aeMrEu8sxbAhYj22jxoyM=")
    }
    
    func testArgon2() throws {
        let password = "password"
        let salt = "salty salt".data(using: .utf8)!
        
        let hasher = Argon2Hasher(mode: .argon2id, iterations: 75, memory: 1024, parallelism: 1, hashLength: 32)
        let hash = hasher.hash(password: password, salt: salt, encoding: .utf8)
        let hashBase64 = hash?.base64EncodedString()
        XCTAssertEqual(hashBase64, "JGFyZ29uMmlkJHY9MTkkbT0xMDI0LHQ9NzUscD0xJGMyRnNkSGtnYzJGc2RBJGd3eFBhVjVPck9XN3oyd3ovcng3WTJmT1JxelR5dVh6Mkt0bkE2WkEwNDgA")
    }
}
