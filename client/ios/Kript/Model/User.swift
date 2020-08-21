//
//  User.swift
//  Kript
//
//  Created by Liam Stevenson on 8/2/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import Foundation

struct User {
    var id: String
    var refreshToken: Kript_Api_RefreshToken
    var accessToken: Kript_Api_AccessToken?
    var publicKey: Data
    var privateKey: Data
    var userObject: Kript_Api_User
    
    func decrypt(datum: Kript_Api_Datum) -> Kript_Api_Secret? {
        // decrypt the symmetric key used for the datum by this user
        guard let encryptedKey = datum.accessors[self.id]?.dataKey.data else { return nil }
        guard let key = self.userObject.public.dataEncryptionAlgorithm.encrypter?.decrypt(data: encryptedKey, privateKey: self.privateKey) else { return nil }
        
        // decrypt the secret data using the symmetric key
        guard let secretData = datum.dataEncryptionAlgorithm.encrypter?.decrypt(data: datum.data.data, key: key, iv: datum.dataIv) else {
            return nil
        }
        guard let secret = try? Kript_Api_Secret(serializedData: secretData) else {
            return nil
        }
        
        return secret
    }
    
    func encryptWithExistingKey(datum: Datum) -> (secret: Kript_Api_ESecret, iv: Data)? {
        guard let encryptedKey = datum.datum?.accessors[self.id]?.dataKey else { return nil }
        
        guard let keyEncrypter = self.userObject.public.dataEncryptionAlgorithm.encrypter else { return nil }
        guard let encryptionAlgorithm = datum.datum?.dataEncryptionAlgorithm else { return nil }
        guard let dataEncrypter = encryptionAlgorithm.encrypter else { return nil }
        
        let dataIv = dataEncrypter.generateIv()
        
        guard let key = keyEncrypter.decrypt(data: encryptedKey.data, privateKey: self.privateKey) else { return nil }
        guard let encryptedSecret = encrypt(secret: datum.secret, encrypter: dataEncrypter, key: key, iv: dataIv) else { return nil }
        
        return (secret: encryptedSecret, iv: dataIv)
    }
    
    func encryptWithNewKey(datum: Datum, encryptionAlgorithm: Kript_Api_SEncryptionAlgorithm = .aes256Cbc) -> (secret: Kript_Api_ESecret, encryptedKey: Kript_Api_EBytes, iv: Data, encryptionAlgorithm: Kript_Api_SEncryptionAlgorithm)? {
        guard let keyEncrypter = self.userObject.public.dataEncryptionAlgorithm.encrypter else { return nil }
        guard let dataEncrypter = encryptionAlgorithm.encrypter else { return nil }
        
        let key = dataEncrypter.generateKey()
        let dataIv = dataEncrypter.generateIv()
        
        guard let encryptedKeyData = keyEncrypter.encrypt(data: key, publicKey: self.publicKey) else { return nil }
        var encryptedKey = Kript_Api_EBytes()
        encryptedKey.data = encryptedKeyData
        
        guard let encryptedSecret = encrypt(secret: datum.secret, encrypter: dataEncrypter, key: key, iv: dataIv) else { return nil }
        
        return (secret: encryptedSecret, encryptedKey: encryptedKey, iv: dataIv, encryptionAlgorithm)
    }
    
    private func encrypt(secret: Kript_Api_Secret, encrypter: SEncrypyer, key: Data, iv: Data) -> Kript_Api_ESecret? {
        guard let encryptedSecretData = try? encrypter.encrypt(data: secret.serializedData(), key: key, iv: iv) else { return nil }
        var encryptedSecret = Kript_Api_ESecret()
        encryptedSecret.data = encryptedSecretData
        return encryptedSecret
    }
}
