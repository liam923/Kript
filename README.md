# Kript

## Table of Contents

- [Kript](#kript)
  * [What is Kript?](#what-is-kript-)
  * [Why use Kript?](#why-use-kript-)
  * [Why trust Kript?](#why-trust-kript-)
  * [Encryption Schema](#encryption-schema)
    + [User password](#user-password)
    + [User symmetric key](#user-symmetric-key)
    + [User public and private keys](#user-public-and-private-keys)
    + [Datum key](#datum-key)
    + [Datum data](#datum-data)
    + [Diagram](#diagram)
      - [Legend](#legend)
    + [Design Choices](#design-choices)
    + [Disadvantages](#disadvantages)
  * [API Guide](#api-guide)
    + [gRPC API](#grpc-api)
      - [How To](#how-to)
        * [Authentication](#authentication)
        * [Create an account](#create-an-account)
        * [Login a user](#login-a-user)
        * [Store a password](#store-a-password)
        * [Retrieve saved passwords](#retrieve-saved-passwords)
        * [Share a password](#share-a-password)
    + [REST API](#rest-api)
  * [Cloud Architecture](#cloud-architecture)
    + [Client](#client)
    + [gRPC API Service](#grpc-api-service)
    + [REST API Service](#rest-api-service)
    + [Account Service](#account-service)
    + [Data Service](#data-service)
    + [Account Database and Data Database](#account-database-and-data-database)
    + [Service Deployments](#service-deployments)
  * [How to Contribute](#how-to-contribute)
    + [Git](#git)
    + [Backend](#backend)
    + [Frontend](#frontend)
  * [Glossary](#glossary)

## What is Kript?

Kript is password manager that is secure, convenient, and transparent. It uses a sophistocated encryption scheme than ensures all data is encrypted behind your master password before it leaves your device, guaranteeing that your data is as secure as your master password. Due to your data being synced to the cloud, you can access your passwords from any device anywhere in the world. Additionally, Kript allows users to share passwords securely, allowing you to grant a family member or friend access to an account. And best of all, Kript is open-source, meaning that you can view and even compile the code yourself to ensure that no funny-business is happening with your passwords.

Kript is built using modern tools and architectures. The backend is written in Go, operates as both a gRPC and REST API, and runs within Docker containers on Google Cloud Platform with a micro-service architecture. Currently, the project includes a very basic iOS frontend written in Swift. Additionally, a CI pipeline using Google Cloud Build has been setup for the project, providing GitHub status checks.

Currently, the project includes a very basic iOS frontent written in Swift to serve as an example of how to correctly encrypt and decrypt data and interact with the API. However, from a user perspective, it is missing many features. Contributions of clients, like a better iOS app, an Android app, a website, or a Chrome plugin, are encouraged and would be greatly appreciated.

## Why use Kript?

There is a very large consensus among cybersecurity experts that there are two easy things someone can do to greatly improve their online security: enable multi-factor authentication and use a password manager. Password managers solve a simple problem which has been plaguing cybersecurity for a while, which is that good passwords are hard to remember. As a result, most people reuse passwords, and many of those are easy to guess ones. If humans had perfect memory, we could instead create complex, unique passwords for each account that we have.

Password managers are a convenient solution to this. A password manager is a piece of software that stores passwords for accounts so that they don't need to be memorized. Instead of needing to remember countless passwords, you only need to know one. This one password is called the master password and is used to login to the password manager. A good password manager will keep the passwords secure by encrypting the passwords, ideally with the master password.

Additionally, a good password manager is easy to use. When logging in to a website, it will automatically fill in the password. And when creating an account, it will automatically generate a long, secure password (like `Suphut-gyswu6-zisrew`) and save it for future use. Unfortunately, Kript does not yet have autofill functionality, as it only has a proof-of-concept frontend at the moment, and so it is not a viable option as a password manager as it currently stands.

## Why trust Kript?

The straightforward answer to this is that you shouldn't. As we all know, strangers on the internet shouldn't be trusted, and that includes me. Passwords are important, and having them stolen can be disastrous. I don't trust anyone with my passwords, and I don't expect anyone else to.

The beauty of Kript, however, is that it doesn't need to be trusted, and this is for two reasons. First, it is designed so that your password never leaves your device. Before any secure data is uploaded to the cloud, it is encrypted, and your master password is the key to decrypt the data. This means that even someone with full access to Kript's servers and databases could not steal your passwords without knowing your master password.

## Encryption Schema

This section describes the protocol used for encrypting "datum"s. In Kript, a "datum" is a piece of data that is meant to be kept secure, like a password, social security number, or private note. These datums are always encrypted and decrypted client side, so understanding the encryption schema is important for properly interacting with the APIs, especially to avoid corrupting data.

This encryption schema was created with the "zero knowledge" principle, which can be described as follows:
- The server and anyone with access to communications between the server and client has zero knowledge of any of the user's secure data. This is achieved by encrypting all secure data with a value only the user knows (their password) before it leaves their device. Thus, no one can steal the user's data unless they get or guess the user's password. This includes me (who has access to the servers and databases), any nefarious person who manages to gain access, or someone who manages to listen in on communications via some sort of man-in-the-middle attack.

Kript uses a mixture of symmetric and asymmetric encryption algorithms. In [symmetric-key encryption](https://en.wikipedia.org/wiki/Symmetric-key_algorithm), the same key is used to encrypt and decrypt a piece of data. In [asymmetric-key encryption](https://en.wikipedia.org/wiki/Public-key_cryptography), there is a public key and a private key that are related to each other, and data is encrypted with the public key and decrypted with the private key. Whenever a symmetric or asymmetric encryption algorithm is used, the client can choose which specific once to use. The list of supported encryption algorithms is found in the protocol definitions at [proto/kript/api/universal.proto](proto/kript/api/universal.proto).

Each of the following sections describe a different component of the encryption schema. Rather than describing how the component is used, these sections describe what the component is and how it uses other components. At the end, there is a diagram demonstrating how they all connect together for clarification, along with some notes about certain design choices. Note that throughout this discussion, a value being stored in an encrypted state on the server also implies that it was transported under the same encryption, as no encryption/decryption happens server-side.

### User password

The user's password is the top level key for all access. All decryption paths start using the user's password. The password is never sent to and subsequently stored on the server. Rather, a salted, hashed password is used for authentication so that the password never leaves the user's device. This hashed password is then hashed again with Bcrypt, which is stored in the database. This is necessary due to the initial hash being client-side. The password salt and the hash algorithm are publicly accessible by calling `kript.api.AccountService/GetUser`.

### User symmetric key

Each user has an associated symmetric key that is a hash of the user's password. This value is never sent to and subsequently stored on the server. It should be noted that this is different than the password hash that is used for authentication, which is stored on the server. It uses a different salt and, optionally, a different hash algorithm, which are available by calling `kript.api.AccountService/GetUser`. (These values will only be included if it is the currently authenticated user being retrieved.)

### User public and private keys

Each user is associated with a public key and a private key, and an asymmetric algorithm that those keys are associated with. These are generated client-side when an account is created. The public key is stored unencrypted on the server, but the private key is symmetrically encrypted using the user's symmetric key (from the previous section) client-side. The algorithm that the public and private keys are used for, the public key, the encrypted private key, the initialization vector for encrypting the private key, and the symmetric algorithm used to encrypt the private key can be retrieved by calling `kript.api.AccountService/GetUser`. The encrypted private key and the algorithm used to encrypt the private key are only included if the user being retrieved is the currently authenticated user, whereas the rest of the values are publicly available.

### Datum key

Each datum (an individual piece of secure data, like a password) has a key associated with it, along with a symmetric encryption algorithm associated with the key. This key is stored multiple times: once for each user who has access to the datum. (Kript allows for users to share datums with each other). For each user who has access to the datum, the key is encrypted with that user's public key (the one mentioned in the previous section) according the asymmetric algorithm associated with it. Thus, the key can be decrypted with the user's private key. All relevant values are retrieved by calling `kript.api.DataService/GetData`.

### Datum data

The core part of each datum is the actual data that contains the secure information. This data is encrypted with the datum's key and its associated symmetric encryption algorithm, as well as an associated initialization vector. These values are all retrieved by calling `kript.api.DataService/GetData`.

### Diagram

[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgXG5zdWJncmFwaCBVc2VyIEEncyBDbGllbnRcblx0cGFzc3dvcmRBKHBhc3N3b3JkKVxuXHRoa1Bhc3N3b3JkQShzeW1tZXRyaWMga2V5KVxuZW5kXG5zdWJncmFwaCBVc2VyIEIncyBDbGllbnRcblx0cGFzc3dvcmRCKHBhc3N3b3JkKVxuXHRoa1Bhc3N3b3JkQihzeW1tZXRyaWMga2V5KVxuZW5kXG4gIFxuc3ViZ3JhcGggRGF0YWJhc2Vcblx0c3ViZ3JhcGggVXNlciBBXG5cdFx0aFBhc3N3b3JkQShwYXNzd29yZClcblx0XHRwdWJsaWNLZXlBKHB1YmxpYyBrZXkpXG5cdFx0cHJpdmF0ZUtleUEocHJpdmF0ZSBrZXkpXG5cdGVuZFxuXHRzdWJncmFwaCBVc2VyIEJcblx0XHRoUGFzc3dvcmRCKHBhc3N3b3JkKVxuXHRcdHB1YmxpY0tleUIocHVibGljIGtleSlcblx0XHRwcml2YXRlS2V5Qihwcml2YXRlIGtleSlcblx0ZW5kXG5cdHN1YmdyYXBoIERhdHVtIE9iamVjdFxuXHRcdGtleUEoQSdzIGRhdHVtIGtleSlcblx0XHRrZXlCKEIncyBkYXR1bSBrZXkpXG5cdFx0ZGF0YShkYXRhKVxuXHRlbmRcbmVuZFxuICBcbnBhc3N3b3JkQS0tPmhQYXNzd29yZEFcbnBhc3N3b3JkQS0tPmhrUGFzc3dvcmRBXG5oa1Bhc3N3b3JkQS0tPnByaXZhdGVLZXlBXG5oa1Bhc3N3b3JkQS0uLT5wcml2YXRlS2V5QVxucHVibGljS2V5QS0tPmtleUFcbnByaXZhdGVLZXlBLS4tPmtleUFcbmtleUEtLT5kYXRhXG5rZXlBLS4tPmRhdGFcbiAgXG5wYXNzd29yZEItLT5oUGFzc3dvcmRCXG5wYXNzd29yZEItLT5oa1Bhc3N3b3JkQlxuaGtQYXNzd29yZEItLT5wcml2YXRlS2V5QlxuaGtQYXNzd29yZEItLi0-cHJpdmF0ZUtleUJcbnB1YmxpY0tleUItLT5rZXlCXG5wcml2YXRlS2V5Qi0uLT5rZXlCXG5rZXlCLS0-ZGF0YVxua2V5Qi0uLT5kYXRhIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQifSwidXBkYXRlRWRpdG9yIjpmYWxzZX0)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoiZ3JhcGggTFJcbiAgXG5zdWJncmFwaCBVc2VyIEEncyBDbGllbnRcblx0cGFzc3dvcmRBKHBhc3N3b3JkKVxuXHRoa1Bhc3N3b3JkQShzeW1tZXRyaWMga2V5KVxuZW5kXG5zdWJncmFwaCBVc2VyIEIncyBDbGllbnRcblx0cGFzc3dvcmRCKHBhc3N3b3JkKVxuXHRoa1Bhc3N3b3JkQihzeW1tZXRyaWMga2V5KVxuZW5kXG4gIFxuc3ViZ3JhcGggRGF0YWJhc2Vcblx0c3ViZ3JhcGggVXNlciBBXG5cdFx0aFBhc3N3b3JkQShwYXNzd29yZClcblx0XHRwdWJsaWNLZXlBKHB1YmxpYyBrZXkpXG5cdFx0cHJpdmF0ZUtleUEocHJpdmF0ZSBrZXkpXG5cdGVuZFxuXHRzdWJncmFwaCBVc2VyIEJcblx0XHRoUGFzc3dvcmRCKHBhc3N3b3JkKVxuXHRcdHB1YmxpY0tleUIocHVibGljIGtleSlcblx0XHRwcml2YXRlS2V5Qihwcml2YXRlIGtleSlcblx0ZW5kXG5cdHN1YmdyYXBoIERhdHVtIE9iamVjdFxuXHRcdGtleUEoQSdzIGRhdHVtIGtleSlcblx0XHRrZXlCKEIncyBkYXR1bSBrZXkpXG5cdFx0ZGF0YShkYXRhKVxuXHRlbmRcbmVuZFxuICBcbnBhc3N3b3JkQS0tPmhQYXNzd29yZEFcbnBhc3N3b3JkQS0tPmhrUGFzc3dvcmRBXG5oa1Bhc3N3b3JkQS0tPnByaXZhdGVLZXlBXG5oa1Bhc3N3b3JkQS0uLT5wcml2YXRlS2V5QVxucHVibGljS2V5QS0tPmtleUFcbnByaXZhdGVLZXlBLS4tPmtleUFcbmtleUEtLT5kYXRhXG5rZXlBLS4tPmRhdGFcbiAgXG5wYXNzd29yZEItLT5oUGFzc3dvcmRCXG5wYXNzd29yZEItLT5oa1Bhc3N3b3JkQlxuaGtQYXNzd29yZEItLT5wcml2YXRlS2V5QlxuaGtQYXNzd29yZEItLi0-cHJpdmF0ZUtleUJcbnB1YmxpY0tleUItLT5rZXlCXG5wcml2YXRlS2V5Qi0uLT5rZXlCXG5rZXlCLS0-ZGF0YVxua2V5Qi0uLT5kYXRhIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQifSwidXBkYXRlRWRpdG9yIjpmYWxzZX0)

#### Legend

[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcblx0QSAtLSBBIGlzIEMncyBwdWJsaWMga2V5IC0tPiBDXG5cdEIgLS4gQiBpcyBDJ3MgcHJpdmF0ZSBrZXkgLi0-IENcblx0RCAtLSBEIGlzIEUncyBzeW1tZXRyaWMga2V5IC0tPiBFXG5cdEQgLS4tPiBFXG5cdEYgLS0gRyBpcyBhIGhhc2ggb2YgRiAtLT4gR1xuIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQifSwidXBkYXRlRWRpdG9yIjpmYWxzZX0)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoiZ3JhcGggTFJcblx0QSAtLSBBIGlzIEMncyBwdWJsaWMga2V5IC0tPiBDXG5cdEIgLS4gQiBpcyBDJ3MgcHJpdmF0ZSBrZXkgLi0-IENcblx0RCAtLSBEIGlzIEUncyBzeW1tZXRyaWMga2V5IC0tPiBFXG5cdEQgLS4tPiBFXG5cdEYgLS0gRyBpcyBhIGhhc2ggb2YgRiAtLT4gR1xuIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQifSwidXBkYXRlRWRpdG9yIjpmYWxzZX0)

### Design Choices

- The user's password salt and the hash algorithm used to hash their password are publicly available because these are needed by the client to hash the password before sending it to the server to authenticate the user.

  - There are a number of advantages to having a public and private key associated with each user rather than using their password to encrypt the datum keys:

  - If the user changes their password, the only value that needs to be re-encrypted is their private key. If instead the user's password was used to encrypt datum keys, each datum key would need to be re-encrypted.

  - Having a public key and private key associated with the user allows for asymmetric encryption to be used for the datum keys, whereas using the user's password would require symmetric encryption. The advantage of symmetric encryption is that anyone can encrypt the value using the public key, which is an essential property to enable sharing. If user A wants to share a datum with user B, A simply needs to (on device) decrypt the datum key with their private key and then encrypt it with B's public key. Then, the B-encrypted datum key can be securely uploaded to the server, granting read access to B and B only.

  - Once the client gains authorization from the server, it no longer needs to store the user's password. Rather, it can store the user's private key in order to encrypt and decrypt passwords. This makes it slightly harder for malware on a client to steal the user's master password.

### Disadvantages

- All encryption happening client-side means that the server is unable to validate encrypted values that are sent to it. Thus, someone creating a client could corrupt data by mis-encrypting it, which would corrupt the data as the server would not know that the data is malformed. Although it would be nice to be able to validate the data, in my research, I found no ways of doing so while retaining the no-knowledge principle.

- All data is encrypted using the user's password, and the server has no way of knowing the user's password. Therefore, if the user forgets their password, there is no way for a user to recover their data. This is an intrinsic issue of the no-knowledge principle, and thus a necessary disadvantage to maintain security.

## API Guide

Note: The API is no longer available at `kript.us` due to cost constraints. To access it, you must instead deploy it yourself.

### gRPC API

The Kript API is primarily available via [gRPC](https://grpc.io) at `grpc.kript.us:443`, and the gRPC interface is described below. The folder [proto](proto) contains the proto definitions for the gRPC API. These proto files define two services, `AccountServer` and `DataService`, both of which are available at `grpc.kript.us:443`.

#### How To

Note that this section will make much more sense if the [Encryption Schema](#encryption-schema) section is read first.

##### Authentication

Kript follows the OAuth model of using refresh tokens and access tokens to authenticate users. The function `AccountService.LoginUser` can be used to retreive an access token and refresh token (see [Login a user](#login-a-user) for information on how to do so). Access tokens stay valid for 24 hours and are used to authenticate the user. Refresh tokens, however, stay valid for 100 years and can be used to retreive fresh access tokens by calling `AccountService.RefreshAuth`.

##### Create an account

 1. Choose a username and master password for the account
 2. Choose an [asymmetric encryption algorithm](proto/kript/api/universal.proto) to use for encrypting data.
 3. Choose a [symmetric encryption algorithm](proto/kript/api/universal.proto)  to use to encrypt the user's private key.
 4. Choose a [hash algorithm](proto/kript/api/universal.proto)  to use for hashing the user's password.
 5. Choose a [hash algorithm](proto/kript/api/universal.proto)  to use for generating the key to the user's private key.
 6. Generate a salt for the user's password and hash it via the algorithm chosen in step 4.
 7. Generate a salt for the key to the user's private key and generate it via the hash algorithm chosen in step 5.
 8. Generate a public and private key for the user that is appropriate for the algorithm chosen in step 2.
 9. Generate an initialization vector for encrypting the user's private key that is appropriate for the symmetric algorithm algorithm chosen in step 3.
 10. Encrypt the user's private key (generated in step 8) using the key generated in step 7 via the symmetric encryption algorithm chosen in step 3 and the initialization vector generated in step 9.
 11. Call `AccountService.CreateAccount` with the following `CreateAccountRequest`:
```
{
  "username": "...", // the user's username
  "password": {
    "data": {
      "data": ... // the hashed user's password, generated in step 6
    }
  },
  "salt": {
    "data": ... // the salt for the user's password hash, generated in step 6
  },
  "password_hash_algorithm": ..., // the hash algorithm chosen for hashing the user's password, chosen in step 4
  "public_key": {
    "data": ... // the user's public key, generated in step 8
  },
  "private_key": {
    "data": {
      "data": ... // the user's encrypted private key, generated in step 10
    }
  },
  "data_encryption_algorithm": ..., // the algorithm chosen for encrypting the user's data, chosen in step 2
  "private_key_encryption_algorithm": ..., // the algorithm chosen for encrypting the user's private key, chosen in step 3
  "private_key_iv": {
    "data": ... // the initialization vector used for encrypting the user's private key, generated in step 9
  },
  "private_key_key_salt": {
    "data": ... // the salt used for hashing the key to the user's private key, generated in step 7
  },
  "private_key_key_hash_algorithm": ... // the hashing algorithm for creating the key to the user's private key, chosen in step 5
}
```
12. If creation is successful, the returned message will include the user's refresh token and access token.

For an example, look at the `Manager.createAccount` function in [Manager.swift](client/ios/Kript/Model/Manager.swift).

##### Login a user

1. Get the user's username and their password.
2. Call `AccountService.GetUser` with the user's username to get the user's information.
3. Hash the user's password. The salt for this, along with the algorithm to use for performing the hash, are included in the data returned in step 2.
4. Call `AccountService.LoginUser` with the following `LoginUserRequest`:
```
{
  "username": ..., // the user's username
  "password": {
    "data": {
      "data": ... // the user's hashed password, generated in step 3
    }
  }
}
```
5. If login is successful and the user does not have two factor authentication set up, the response will include the user's refresh token and access token.

For an example, look at the `login` function in [Manager.swift](client/ios/Kript/Model/Manager.swift).

##### Store a password

1. Get the user's username and call `AccountService.GetUser` with it.
2. Create a `Secret` proto message (see [encrypt.proto](proto/kript/api/encrypt.proto)), which will be the encrypted data uploaded.
3. Encode the proto from step 2 into bytes.
4. Choose a [symmetric encryption algorithm](proto/kript/api/universal.proto) to use to encrypt the data.
5. Generate a key and initialization vector for encrypting the data key, appropriate for the algorithm chosen in step 4.
6. Encrypt the proto data from step 3 using the encryption algorithm and key from step 4.
7. Encrypt the data key from step 4 using the user's data encryption algorithm and their public key, which are included in the response from step 1.
8. Call `AccountService.CreateDatum` using the following `CreateDatumRequest`:
```
{
  "access_token": ..., // the user's access token
  "data": {
    "data": {
      "data": ... // the encrypted data, obtained in step 6
    }
  },
  "data_key": {
    "data": {
      "data": ... // the encrypted key for the data, generated in step 5
    }
  },
  "data_encryption_algorithm": ..., // the symmetric algorithm for encypting the data, chosen in step 4
  "data_iv": {
    "data": ... // the initialization vector for encrypting the data, generated in step 5
  }
}
```

For an example, look at the `Manager.add` function in [Manager.swift](client/ios/Kript/Model/Manager.swift).

##### Retrieve saved passwords

1. Get the username and password for the user.
2. Call `AccountService.GetUser` with the user's username to get the user's information.
3. Re-create the key to the user's private key  by hashing their password using the hashing algorithm and salt specified by the `user.private.private_key_key_hash_algorithm` and `user.private.private_key_key_salt` fields of the response from step 2.
4. Decrypt the user's private key using the key obtained in step 3 via the symmetric encryption algorithm specified by the `user.private.private_key_encryption_algorithm` field of the response from step 2, along with the initialization vector specified at `user.private.private_key_iv`.
5. Call `DataService.GetData` with the following `GetDataRequest`:
```
{
  "access_token": ... // the user's access token
}
```
6. The response will be a list of `Datum`s. For each, do the following to decrypt them:
    1. Get the datum key that is encrypted behind the user's private key, which is located at `accessors[user_id].data_key`, where `user_id` is the id of the user (which is located in the response from step 2).
    2. Decrypt the data key from step 6.1 using the user's private key from step 4, via the asymmetric encryption algorithm specified by the `user.public.data_encryption_algorithm` field of the response from step 2.
    3. Decrypt the data at field `data` using the key from step 6.2 and the initialization vector at field `data_iv` via the symmetric encryption algorithm at field `data_encryption_algorithm`.

For an example, look at the `Manager.refresh` function in [Manager.swift](client/ios/Kript/Model/Manager.swift).

##### Share a password

1. Obtain the information of the datum to be shared. This can be done by performing the steps from [Retrieve saved passwords](#retrieve-saved-passwords) up to 6.2.
2. Call `AccountService.GetUser` with the username of the user to share the datum with.
3. Encrypt the datum key (from step 6.2 of "Retrieve saved passwords") using the public key of the user to share the data with, via their data encryption algorithm. These values are both located in the response from step 2.
4. Choose the permissions to grant the user on the datum.
5. Call `DataService.ShareDatum` with the following `ShareDatumRequest`:
```
{
  "access_token": ..., // the user's access token
  "id": ..., // the id of the datum to share
  "target_id": "9ff6edb2-9d06-41dc-b173-8b0865be83f6", // the user id of the user to share the data with, from the reponse in step 2
  "data_key": ..., // the data key encrypted using the target user's public key, from step 3
  "permissions": // the permission to grant, chosen in step 5
}
```

### REST API

There is also a REST API available available at `https://api.kript.us` that was generated using [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway). Documentation is available for it at [server/docs/api/kript.swagger.json](server/docs/api/kript.swagger.json). Each REST endpoint corresponds to a gRPC function and vice-versa, so using it is very similar to using the gRPC API. The only difference is in how calls to the API are made.

## Cloud Architecture

The below graph provides an overview of the architecture:
[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBYKChDbGllbnQpKVxuICAgIHN1YmdyYXBoIEdvb2dsZSBDbG91ZCBQbGF0Zm9ybVxuXHQgICAgQVtnUlBDIEFQSSBTZXJ2aWNlXVxuXHQgICAgQltSRVNUIEFQSSBTZXJ2aWNlXVxuXHQgICAgQ1tBY2NvdW50IFNlcnZpY2VdXG5cdCAgICBEW0RhdGEgU2VydmljZV1cblx0ICAgIEVbKEFjY291bnQgRGF0YWJhc2UpXVxuXHQgICAgRlsoRGF0YSBEYXRhYmFzZSldXG4gICAgZW5kXG5cbiAgICBYIC0uIGdycGMua3JpcHQudXMgLi0-IEFcbiAgICBYIC0tIGFwaS5rcmlwdC51cyAtLT4gQlxuICAgIEEgLS4tPiBDXG4gICAgQSAtLi0-IERcbiAgICBCIC0uLT4gQ1xuICAgIEIgLS4tPiBEXG4gICAgQyAtLi0-IEVcbiAgICBEIC0uLT4gRiIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBYKChDbGllbnQpKVxuICAgIHN1YmdyYXBoIEdvb2dsZSBDbG91ZCBQbGF0Zm9ybVxuXHQgICAgQVtnUlBDIEFQSSBTZXJ2aWNlXVxuXHQgICAgQltSRVNUIEFQSSBTZXJ2aWNlXVxuXHQgICAgQ1tBY2NvdW50IFNlcnZpY2VdXG5cdCAgICBEW0RhdGEgU2VydmljZV1cblx0ICAgIEVbKEFjY291bnQgRGF0YWJhc2UpXVxuXHQgICAgRlsoRGF0YSBEYXRhYmFzZSldXG4gICAgZW5kXG5cbiAgICBYIC0uIGdycGMua3JpcHQudXMgLi0-IEFcbiAgICBYIC0tIGFwaS5rcmlwdC51cyAtLT4gQlxuICAgIEEgLS4tPiBDXG4gICAgQSAtLi0-IERcbiAgICBCIC0uLT4gQ1xuICAgIEIgLS4tPiBEXG4gICAgQyAtLi0-IEVcbiAgICBEIC0uLT4gRiIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)

Legend:
- Dotted lines represent communication over gRPC.
- Normal lines represent communication over https.

### Client

The client has the choice to interact with either the gRPC or REST API, as both are fully featured. While the gRPC connection has a performance advantage, https provides better compatibility.

### gRPC API Service

The gRPC API Service simply merges together the different micro-services (the account and data micro-services) to provide one single endpoint. It can be reached at grpc.kript.us.

### REST API Service

The REST API Service also merges together the different micro-services. However, it uses [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway/) to provide the API as a REST one. It can be reached at api.kript.us.

### Account Service

The account micro-service is responsible for implementing all functions defined in the [kript.api.AccountService service](proto/kript/api/account.proto). These are all functions related to account creation, account management, and authorization (logging users in).

### Data Service

The data micro-service is responsible for implementing all functions defined in the [kript.api.DataService service](proto/kript/api/data.proto). These are all functions related to creating, fetching, sharing, and deleting data. ("data" is referring to a collection of "datum"s. A "datum" in Kript is a piece of secure information, like a password, note, or code.)

### Account Database and Data Database

The Account Database and Data Database are Cloud Firestore instances that store information for the account and data micro-services, respectively. (This is actually a white lie. They reside in one Firestore instance, as a GCP project is only allowed to have one Firestore instance. However, their data reside in separate "collections" within the instance and are completely separate; neither micro-service ever interacts with the other's data, so they are effectively treated as two separate databases.)

### Service Deployments

Each of the four services (the gRPC Service, REST Service, Account Service, and Data Service) is containerized using Docker and deployed on Cloud Run, which allows for easy deployment, load-balancing, and autoscaling.

## How to Contribute

### Git

In order to contribute, please open a pull request! To learn about pull requests, look [here](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests). Pull requests should be target the development branch, not master.

Kript uses [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) to format commits in order to keep the history clean. Also, please squash commits on pull requests.

### Backend

The backend code is located in the [server] folder. To run tests, run `make test`. To build the services, run `make build`. This will create the directory `server/bin`, which will contain executables for each service. If the proto files are edited, run `make setup` afterwards to recompile with `protoc`. Note that the proto files located at `proto`, not `server/docs/proto`, should be edited.

### Frontend

The code for each client is located in the directory `client`. The code for the proof-of-concept iOS app is at `client/ios` and is structured as a typical CocoaPods project. If the proto files are edited, run `make setup` afterwards to recompile with `protoc`. Note that the proto files located at `proto`, not `client/docs/proto`, should be edited.

## Glossary

- asymmetric-key encryption: An encryption algorithm that uses one key to encrypt data and a different key to decrypt the same data. The key used to encrypt the data is usually referred to as the public key because it is usually public knowledge. The key used for decryption is usually called the private key. This is usually used when you want to encrypt data for someone else to decrypt. See [Wikipedia](https://en.wikipedia.org/wiki/Public-key_cryptography) for more info.
- master password: The password used to login to Kript and encrypt data.
- password manager: A program that stores secure data like passwords.
- private key: In asymmetric-key encryption, the key that is used to decrypt data. It is not publicly known, as that would defeat the purpose of encrypting the data.
- public key: In asymmetric-key encryption, the key that is used to encrypt data. It is usually publicly known since it does not give any ability to decrypt the data.
- symmetric-key encryption: An encryption algorithm that uses the same key to encrypt and decrypt data. This is usually used when you want to encrypt data that only you can decrypt. See [Wikipedia](https://en.wikipedia.org/wiki/Symmetric-key_algorithm) for more info.