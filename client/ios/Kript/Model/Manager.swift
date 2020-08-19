//
//  Manager.swift
//  Kript
//
//  Created by Liam Stevenson on 8/2/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

import CoreData
import Foundation
import GRPC
import JWTDecode
import KeychainSwift
import NIO
import SwiftUI

class Manager {
    private static let idKey = "id"
    private static let refreshTokenKey = "refreshToken"
    private static let accessTokenKey = "accessToken"
    private static let publicKeyKey = "publicKey"
    private static let privateKeyKey = "privateKey"
    private static let dataEncryptionAlgorithmKey = "dataEncryptionAlgorithm"
    private static let keychainAccessLevel = KeychainSwiftAccessOptions.accessibleWhenUnlockedThisDeviceOnly
    
    private static let passwordHashAlgorithm = Kript_Api_HashAlgorithm.argon2
    private static let pivateKeyGenerationAlgorithm = Kript_Api_HashAlgorithm.pbkdf2Sha512
    private static let dataEncryptionAlgorithm = Kript_Api_AEncryptionAlgorithm.rsa
    private static let privateKeyEncryptionAlgorithm = Kript_Api_SEncryptionAlgorithm.aes256Cbc
    
    private static let cdDatumName = "CDDatum"
    
    let keychain: KeychainSwift
    let accountClient: Kript_Api_AccountServiceClient
    let dataClient: Kript_Api_DataServiceClient
    
    init(keychain: KeychainSwift, accountClient: Kript_Api_AccountServiceClient, dataClient: Kript_Api_DataServiceClient) {
        self.keychain = keychain
        self.accountClient = accountClient
        self.dataClient = dataClient
    }
    
    // MARK: User Management
    
    enum LoginResponse {
        case complete(user: User)
        case twoFactor(verificationToken: Kript_Api_VerificationToken, options: [String: Kript_Api_TwoFactor])
        case badUsername
        case badPassword
        case otherError // connection error, etc
    }
    func login(username: String, password: String, completion: @escaping (LoginResponse) -> ()) {
        var request = Kript_Api_GetUserRequest()
        request.username = username
        let response = self.accountClient.getUser(request).response
        response.whenFailure { error in
            if let status = error as? GRPCStatus, status.code == .notFound {
                completion(.badUsername)
            } else {
                completion(.otherError)
            }
        }
        response.whenSuccess { response in
            let publicUser = response.user.public
            
            guard let passwordHash = publicUser.passwordHashAlgorithm.hasher?.hash(password: password, salt: publicUser.passwordSalt) else {
                completion(.otherError)
                return
            }
            
            var request = Kript_Api_LoginUserRequest()
            request.userID = publicUser.id
            request.password.data = passwordHash
            
            let response = self.accountClient.loginUser(request).response
            response.whenFailure { error in
                if let status = error as? GRPCStatus, status.code == .unauthenticated {
                    completion(.badPassword)
                } else {
                    completion(.otherError)
                }
            }
            response.whenSuccess { response in
                switch response.responseType {
                case .twoFactor(let twoFactorInfo):
                    completion(.twoFactor(verificationToken: twoFactorInfo.verificationToken, options: twoFactorInfo.options))
                case .response(let loginResponse):
                    if let user = self.processLoginSuccess(response: loginResponse, password: password) {
                        completion(.complete(user: user))
                    } else {
                        completion(.otherError)
                    }
                default:
                    completion(.otherError)
                }
            }
        }
    }
    
    enum CreateResponse {
        case complete(user: User)
        case usernameTaken
        case otherError // connection error, etc
    }
    func createAccount(username: String, password: String, completion: @escaping (CreateResponse) -> ()) {
        let passwordHashAlgorithm = Self.passwordHashAlgorithm
        let pivateKeyGenerationAlgorithm = Self.pivateKeyGenerationAlgorithm
        let dataEncryptionAlgorithm = Self.dataEncryptionAlgorithm
        let privateKeyEncryptionAlgorithm = Self.privateKeyEncryptionAlgorithm
        
        guard let passwordHasher = passwordHashAlgorithm.hasher,
            let pivateKeyGenerator = pivateKeyGenerationAlgorithm.hasher,
            let dataEncrypter = dataEncryptionAlgorithm.encrypter,
            let privateKeyEncrypter = privateKeyEncryptionAlgorithm.encrypter else {
            completion(.otherError)
            return
        }
        guard let (publicKey, privateKey) = dataEncrypter.generateKeyPair() else {
            completion(.otherError)
            return
        }
        
        guard let salt = passwordHasher.generateSalt(),
            let hashedPasswordData = passwordHasher.hash(password: password, salt: salt) else {
            completion(.otherError)
            return
        }
        
        guard let privateKeyKeySalt = pivateKeyGenerator.generateSalt(),
            let privateKeyKey = pivateKeyGenerator.hash(password: password, salt: privateKeyKeySalt) else {
            completion(.otherError)
            return
        }
        
        let privateKeyIv = privateKeyEncrypter.generateIv()
        guard let privateKeyEncryptedData = privateKeyEncrypter.encrypt(data: privateKey, key: privateKeyKey, iv: privateKeyIv) else {
            completion(.otherError)
            return
        }
        
        var request = Kript_Api_CreateAccountRequest()
        request.username = username
        request.password.data = hashedPasswordData
        request.salt = salt
        request.passwordHashAlgorithm = passwordHashAlgorithm
        request.publicKey = publicKey
        request.privateKey.data = privateKeyEncryptedData
        request.dataEncryptionAlgorithm = dataEncryptionAlgorithm
        request.privateKeyEncryptionAlgorithm = privateKeyEncryptionAlgorithm
        request.privateKeyIv = privateKeyIv
        request.privateKeyKeySalt = privateKeyKeySalt
        request.privateKeyKeyHashAlgorithm = pivateKeyGenerationAlgorithm
        let response = self.accountClient.createAccount(request).response
        response.whenFailure { error in
            if let status = error as? GRPCStatus, status.code == .alreadyExists {
                completion(.usernameTaken)
            } else {
                completion(.otherError)
            }
        }
        response.whenSuccess { response in
            if let user = self.processLoginSuccess(response: response.response, password: password) {
                completion(.complete(user: user))
            } else {
                completion(.otherError)
            }
        }
    }
    
    private func processLoginSuccess(response: Kript_Api_SuccessfulLoginMessage, password: String) -> User? {
        guard let privateKeyKey = response.user.private.privateKeyKeyHashAlgorithm.hasher?.hash(password: password, salt: response.user.private.privateKeyKeySalt) else { return nil }
        
        let encryptedPrivateKey = response.user.private.privateKey
        guard let privateKey = response.user.private.privateKeyEncryptionAlgorithm.encrypter?.decrypt(data: encryptedPrivateKey.data, key: privateKeyKey, iv: response.user.private.privateKeyIv) else { return nil }
        
        return User(id: response.user.public.id,
                    refreshToken: response.refreshToken,
                    accessToken: response.accessToken,
                    publicKey: response.user.public.publicKey,
                    privateKey: privateKey,
                    dataEncryptionAlgorithm: response.user.public.dataEncryptionAlgorithm)
    }
    
    func loadUser() -> User? {
        guard let id = keychain.get(Self.idKey) else {
            return nil
        }
        guard let refreshTokenData = keychain.getData(Self.refreshTokenKey),
            let refreshToken = try? Kript_Api_RefreshToken(serializedData: refreshTokenData) else {
            return nil
        }
        var accessToken: Kript_Api_AccessToken?
        if let accessTokenData = keychain.getData(Self.accessTokenKey) {
            accessToken = try? Kript_Api_AccessToken(serializedData: accessTokenData)
        }
        guard let publicKey = keychain.getData(Self.publicKeyKey) else {
            return nil
        }
        guard let privateKey = keychain.getData(Self.privateKeyKey) else {
            return nil
        }
        
        guard let dataEncryptionAlgorithm = Kript_Api_AEncryptionAlgorithm(rawValue: Int(keychain.get(Self.dataEncryptionAlgorithmKey) ?? "not int") ?? -1) else {
            return nil
        }
        
        return User(id: id,
                    refreshToken: refreshToken,
                    accessToken: accessToken,
                    publicKey: publicKey,
                    privateKey: privateKey,
                    dataEncryptionAlgorithm: dataEncryptionAlgorithm)
    }
    
    func saveUserDataToKeychain(user: User) {
        self.keychain.set(user.id, forKey: Self.idKey, withAccess: Self.keychainAccessLevel)
        if let refreshToken = try? user.refreshToken.serializedData() {
            self.keychain.set(refreshToken, forKey: Self.refreshTokenKey, withAccess: Self.keychainAccessLevel)
        } else {
            self.keychain.delete(Self.refreshTokenKey)
        }
        if let accessToken = try? user.accessToken?.serializedData() {
            self.keychain.set(accessToken, forKey: Self.accessTokenKey, withAccess: Self.keychainAccessLevel)
        } else {
            self.keychain.delete(Self.accessTokenKey)
        }
        self.keychain.set(user.publicKey, forKey: Self.publicKeyKey, withAccess: Self.keychainAccessLevel)
        self.keychain.set(user.privateKey, forKey: Self.privateKeyKey, withAccess: Self.keychainAccessLevel)
        self.keychain.set(String(user.dataEncryptionAlgorithm.rawValue), forKey: Self.dataEncryptionAlgorithmKey, withAccess: Self.keychainAccessLevel)
    }
    
    private func updateAccessToken(user: User?, completion: @escaping (User?, Bool) -> ()) {
        if let user = user {
            if let expiration = try? decode(jwt: user.accessToken?.jwt.token ?? "").expiresAt, expiration > Date() + 60 {
                completion(user, false)
            } else if let expiration = try? decode(jwt: user.refreshToken.jwt.token).expiresAt, expiration > Date() + 60 {
                var refreshAuthRequest = Kript_Api_RefreshAuthRequest()
                refreshAuthRequest.refreshToken = user.refreshToken
                accountClient.refreshAuth(refreshAuthRequest).response.whenComplete { result in
                    do {
                        let response = try result.get()
                        var newUser = user
                        newUser.accessToken = response.accessToken
                        completion(newUser, false)
                    } catch {
                        if let status = error as? GRPCStatus, status.code == .unauthenticated {
                            completion(nil, true)
                        } else {
                            completion(nil, false)
                        }
                    }
                }
            } else {
                completion(nil, true)
            }
        } else {
            completion(nil, false)
        }
    }
    
    func logoutUser() {
        // delete saved datums
        let context = self.persistentContainer.viewContext
        let fetchRequest: NSFetchRequest<NSFetchRequestResult> = NSFetchRequest(entityName: Self.cdDatumName)
        let deleteRequest = NSBatchDeleteRequest(fetchRequest: fetchRequest)
        _ = try? context.execute(deleteRequest)
        
        // remove user login information
        self.keychain.delete(Self.idKey)
        self.keychain.delete(Self.refreshTokenKey)
        self.keychain.delete(Self.accessTokenKey)
        self.keychain.delete(Self.publicKeyKey)
        self.keychain.delete(Self.privateKeyKey)
        self.keychain.delete(Self.dataEncryptionAlgorithmKey)
    }
    
    // MARK: Store Management
    
    func loadStoreFromCoreData(user: User) -> Store {
        let context = self.persistentContainer.viewContext
        let cdDatums = try? context.fetch(NSFetchRequest(entityName: Self.cdDatumName)) as? [CDDatum]
        let datums = cdDatums?.compactMap { Datum(fromCD: $0, user: user) }
        
        return Store(datums: datums ?? [])
    }
    
    @discardableResult func saveStoreToCoreData(store: Store) -> Bool {
        let context = self.persistentContainer.viewContext
        
        // delete all existing datums
        let fetchRequest: NSFetchRequest<NSFetchRequestResult> = NSFetchRequest(entityName: Self.cdDatumName)
        let deleteRequest = NSBatchDeleteRequest(fetchRequest: fetchRequest)
        _ = try? context.execute(deleteRequest)
        
        store.datums.compactMap({ $0.datum }).forEach { datum in
            let cdDatum = CDDatum(context: context)
            do {
                try cdDatum.setProto(fromDatum: datum)
            } catch {
                context.delete(cdDatum)
            }
        }
        
        return (try? context.save()) != nil
    }
    
    typealias CompletionHandler<Response> = (Response) -> ()
    
    func refresh(store: Binding<Store>, forUser user: Binding<User?>, completion: @escaping CompletionHandler<Bool>) {
        performOperation(request: { accessToken in
            var request = Kript_Api_GetDataRequest()
            request.accessToken = accessToken
            request.datumIds = []
            return request
        },
                         grpcCall: self.dataClient.getData,
                         user: user,
                         completionFailureResponse: false,
                         completion: completion) { response in
                            guard let user = user.wrappedValue else { return false }
                            var idMap = [String:UUID]()
                            for datum in store.wrappedValue.datums {
                                if let id = datum.datum?.id {
                                    idMap[id] = datum.id
                                }
                            }
                            let datums = response.datums.compactMap { datum -> Datum? in
                                guard let secret = user.decrypt(datum: datum) else {
                                    return nil
                                }
                                return Datum(datum: datum, secret: secret, id: idMap[datum.id])
                            }
                            store.wrappedValue = Store(datums: datums)
                            return true
        }
    }
    
    func add(datum: Datum, forUser user: Binding<User?>, completion: @escaping CompletionHandler<Kript_Api_Datum?>) {
        performOperation(request: { accessToken in
            if let (secret, key, iv, encryptionAlgorithm) = user.wrappedValue?.encryptWithNewKey(datum: datum) {
                var request = Kript_Api_CreateDatumRequest()
                request.accessToken = accessToken
                request.data = secret
                request.dataKey = key
                request.dataEncryptionAlgorithm = encryptionAlgorithm
                request.dataIv = iv
                return request
            } else {
                return nil
            }
        },
                         grpcCall: self.dataClient.createDatum,
                         user: user,
                         completionFailureResponse: nil,
                         completion: completion,
                         whenSuccess: { $0.datum })
    }
    
    func remove(datum: Datum, forUser user: Binding<User?>, completion: @escaping CompletionHandler<Kript_Api_Datum?>) {
        performOperation(request: { accessToken in
            var request = Kript_Api_DeleteDatumRequest()
            request.accessToken = accessToken
            guard let id = datum.datum?.id else { return nil }
            request.id = id
            return request
        },
                         grpcCall: dataClient.deleteDatum,
                         user: user,
                         completionFailureResponse: nil,
                         completion: completion,
                         whenSuccess: { $0.datum })
    }
    
    func update(datum: Datum, forUser user: Binding<User?>, completion: @escaping CompletionHandler<Kript_Api_Datum?>) {
        performOperation(request: { accessToken in
            if let id = datum.datum?.id, let (secret, iv) = user.wrappedValue?.encryptWithExistingKey(datum: datum) {
                var request = Kript_Api_UpdateDatumRequest()
                request.accessToken = accessToken
                request.id = id
                request.data = secret
                request.dataIv = iv
                return request
            } else {
                return nil
            }
        },
                         grpcCall: self.dataClient.updateDatum,
                         user: user,
                         completionFailureResponse: nil,
                         completion: completion,
                         whenSuccess: { $0.datum })
    }
    
    private func performOperation<Request, Response, CompletionResponse>(request: @escaping (Kript_Api_AccessToken) -> Request?,
                                                                         grpcCall: @escaping (Request, CallOptions?) -> UnaryCall<Request, Response>,
                                                                         user: Binding<User?>,
                                                                         completionFailureResponse: CompletionResponse,
                                                                         completion: @escaping CompletionHandler<CompletionResponse>,
                                                                         whenSuccess: @escaping (Response) -> (CompletionResponse)) {
        DispatchQueue.main.async {
            self.updateAccessToken(user: user.wrappedValue) { (newUser, badRefreshToken) in
                if let newUser = newUser {
                    user.wrappedValue = newUser
                    
                    guard let accessToken = user.wrappedValue?.accessToken else {
                        completion(completionFailureResponse)
                        return
                    }
                    if let request = request(accessToken) {
                        let call = grpcCall(request, nil).response
                        call.whenSuccess { response in
                            let completionResponse = whenSuccess(response)
                            completion(completionResponse)
                        }
                        call.whenFailure { _ in
                            completion(completionFailureResponse)
                        }
                    } else {
                        completion(completionFailureResponse)
                    }
                } else if badRefreshToken {
                    user.wrappedValue = nil
                    completion(completionFailureResponse)
                } else {
                    completion(completionFailureResponse)
                }
            }
        }
    }
    
    // MARK: Core Data
    
    lazy var persistentContainer: NSPersistentContainer = {
        /*
         The persistent container for the application. This implementation
         creates and returns a container, having loaded the store for the
         application to it. This property is optional since there are legitimate
         error conditions that could cause the creation of the store to fail.
         */
        let container = NSPersistentContainer(name: "Kript")
        container.loadPersistentStores(completionHandler: { (storeDescription, error) in
            if let error = error {
                print(error)
            }
        })
        return container
    }()
    
    func saveContext () {
        let context = persistentContainer.viewContext
        if context.hasChanges {
            do {
                try context.save()
            } catch {
                // Replace this implementation with code to handle the error appropriately.
                // fatalError() causes the application to generate a crash log and terminate. You should not use this function in a shipping application, although it may be useful during development.
                let nserror = error as NSError
                print(nserror.userInfo)
            }
        }
    }
}

class MockManager: Manager {
    init() {
        let group = MultiThreadedEventLoopGroup(numberOfThreads: 1)
        let accountChannel = ClientConnection
            .insecure(group: group)
            .connect(host: "localhost", port: 9232)
        let dataChannel = ClientConnection
            .insecure(group: group)
            .connect(host: "localhost", port: 9233)
        
        let accountClient = Kript_Api_AccountServiceClient(channel: accountChannel, defaultCallOptions: CallOptions(timeLimit: TimeLimit.timeout(.minutes(1))))
        let dataClient = Kript_Api_DataServiceClient(channel: dataChannel, defaultCallOptions: CallOptions(timeLimit: TimeLimit.timeout(.minutes(1))))
        
        super.init(keychain: KeychainSwift(), accountClient: accountClient, dataClient: dataClient)
    }
    
    override func login(username: String, password: String, completion: @escaping (Manager.LoginResponse) -> ()) {
        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            completion(.complete(user: User(id: "id", refreshToken: Kript_Api_RefreshToken(), accessToken: Kript_Api_AccessToken(), publicKey: Data(), privateKey: Data(), dataEncryptionAlgorithm: .rsa)))
        }
    }
    
    override func createAccount(username: String, password: String, completion: @escaping (Manager.CreateResponse) -> ()) {
        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            completion(.complete(user: User(id: "id", refreshToken: Kript_Api_RefreshToken(), accessToken: Kript_Api_AccessToken(), publicKey: Data(), privateKey: Data(), dataEncryptionAlgorithm: .rsa)))
        }
    }
    
    override func add(datum: Datum, forUser user: Binding<User?>, completion: @escaping Manager.CompletionHandler<Kript_Api_Datum?>) {
        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            completion(Kript_Api_Datum())
        }
    }
    
    override func remove(datum: Datum, forUser user: Binding<User?>, completion: @escaping Manager.CompletionHandler<Kript_Api_Datum?>) {
        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            completion(Kript_Api_Datum())
        }
    }
    
    override func update(datum: Datum, forUser user: Binding<User?>, completion: @escaping Manager.CompletionHandler<Kript_Api_Datum?>) {
        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            completion(Kript_Api_Datum())
        }
    }
    
    override func refresh(store: Binding<Store>, forUser user: Binding<User?>, completion: @escaping Manager.CompletionHandler<Bool>) {
        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            let datums = [Datum(), Datum()]
            store.wrappedValue = Store(datums: datums)
            completion(true)
        }
    }
}
