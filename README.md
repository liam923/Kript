
# Kript

Kript is a password manager designed to be both secure, by using a sophisticated encryption scheme, and transparent, by being open-source. All data is secured behind the user's password before it leaves the device, meaning that man in the middle attacks and database breaches will not reveal any sensitive data to the intruders.

Kript is built using modern tools and architectures. The backend is written in Go, operates as both a gRPC and REST  API, and runs within Docker containers on Google Cloud Platform with a micro-service architecture. This project also includes an iOS frontend written in Swift and encourages contributions of other frontends, like Android, Web, and a Chrome plugin. Additionally, a CI pipeline using Google Cloud Build has been setup for the project, providing GitHub status checks.

## Why trust Kript?

Passwords are important. They act as our identification throughout the internet, and so anyone with our passwords can steal our identity online. Thus, trusting a password manager is important. So why trust Kript?

The beauty of Kript is that you don't have to trust it. There are two components of Kript: the client and the server. Your passwords never leave your device, so, there is nothing nefarious that a bad actor who has access to the server can do. And for the client, you can build and run it from source, as all code is open-source and located in this repository. (You can even run your own instance of the backend if you'd like!)

## Encryption Schema

This section describes the protocol used for encrypting "datum"s. In Kript, a "datum" is a piece of data that is meant to be kept secure, like a password, social security number, or private note. These datums are always encrypted and decrypted client side, so understanding the encryption schema is important for properly interacting with the APIs, especially to avoid corrupting data.

This encryption schema was created with the "no-knowledge principle", which can be described as follows:
- The server and anyone with access to communications between the server and client has zero knowledge of any of the user's secure data. This is achieved by encrypting all secure data with a value only the user knows (their password) before it leaves their device. Thus, no one can steal the user's data unless they get or guess the user's password. This includes me (who has access to the servers and databases), any nefarious person who manages to gain access, or someone who manages to listen in on communications via some sort of man-in-the-middle attack.

Kript uses a mixture of symmetric and asymmetric encryption algorithms. In [symmetric-key encryption](https://en.wikipedia.org/wiki/Symmetric-key_algorithm), the same key is used to encrypt and decrypt a piece of data. In [asymmetric-key encryption](https://en.wikipedia.org/wiki/Public-key_cryptography), there is a public key and a private key that are related to each other, and data is encrypted with the public key and decrypted with the private key. Whenever a symmetric or asymmetric encryption algorithm is used, the client can choose which specific once to use. The list of supported encryption algorithms is found in the protocol definitions at [proto/kript/api/universal.proto](proto/kript/api/universal.proto).

Each of the following sections describe a different component of the encryption schema. Rather than describing how the component is used, these sections describe what the component is and how it uses other components. At the end, there is a diagram demonstrating how they all connect together for clarification, along with some notes about certain design choices. Note that throughout this discussion, a value being stored in an encrypted state on the server also implies that it was transported under the same encryption, as no encryption/decryption happens server-side.

### User password

The user's password is the top level key for all access. All decryption paths start using the user's password. The password is never sent to and subsequently stored on the server. Rather, a salted, hashed password is used for authentication so that the password never leaves the user's device. The password salt and the hash algorithm are publicly accessible by calling `kript.api.AccountService/GetUser`.

### User public and private keys

Each user is associated with a public key and a private key, and an asymmetric algorithm that those keys are associated with. These are generated client-side when an account is created. The public key is stored unencrypted on the server, but the private key is symmetrically encrypted using the user's password (and sent to the server encrypted). The algorithm that the public and private keys are used for, the public key, the encrypted private key, and the symmetric algorithm used to encrypt the private key, can be retrieved by calling `kript.api.AccountService/GetUser`. The encrypted private key and the algorithm used to encrypt the private key are only included if the user being retrieved is the currently authenticated user, whereas the rest of the values are publicly available.

### Datum key

Each datum (an individual piece of secure data, like a password) has a key associated with it, along with a symmetric encryption algorithm associated with the key. This key is stored multiple times: once for each user who has access to the datum. (Kript allows for users to share datums with each other). For each user who has access to the datum, the key is encrypted with that user's public key (the one mentioned in the previous section) according the asymmetric algorithm associated with it. Thus, the key can be decrypted with the user's private key.

### Datum data

The core part of each datum is the actual data that contains the secure information. This data is encrypted with the datum's key and its associated symmetric encryption algorithm.

### Diagram

[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcblxuc3ViZ3JhcGggTGVnZW5kXG5cdEEgLS0gQSBpcyB1c2VkIHRvIGhhc2ggb3IgZW5jcnlwdCBCIC0tPiBCXG5cdEMgLS4gQyBpcyB1c2VkIHRvIGRlY3J5cHQgRCAuLT4gRFxuZW5kXG5cbnN1YmdyYXBoIFVzZXIgQSdzIE1pbmRcblx0cGFzc3dvcmRBKHBhc3N3b3JkKVxuZW5kXG5zdWJncmFwaCBVc2VyIEIncyBNaW5kXG5cdHBhc3N3b3JkQihwYXNzd29yZClcbmVuZFxuXG5zdWJncmFwaCBEYXRhYmFzZVxuXHRzdWJncmFwaCBVc2VyIEFcblx0XHRoUGFzc3dvcmRBKHBhc3N3b3JkKVxuXHRcdHB1YmxpY0tleUEocHVibGljIGtleSlcblx0XHRwcml2YXRlS2V5QShwcml2YXRlIGtleSlcblx0ZW5kXG5cdHN1YmdyYXBoIFVzZXIgQlxuXHRcdGhQYXNzd29yZEIocGFzc3dvcmQpXG5cdFx0cHVibGljS2V5QihwdWJsaWMga2V5KVxuXHRcdHByaXZhdGVLZXlCKHByaXZhdGUga2V5KVxuXHRlbmRcblx0c3ViZ3JhcGggRGF0dW0gT2JqZWN0XG5cdFx0a2V5QShrZXkpXG5cdFx0a2V5QihrZXkpXG5cdFx0ZGF0YShkYXRhKVxuXHRlbmRcbmVuZFxuXG5wYXNzd29yZEEtLT5oUGFzc3dvcmRBXG5wYXNzd29yZEEtLT5wcml2YXRlS2V5QVxucGFzc3dvcmRBLS4tPnByaXZhdGVLZXlBXG5wdWJsaWNLZXlBLS0-a2V5QVxucHJpdmF0ZUtleUEtLi0-a2V5QVxua2V5QS0tPmRhdGFcbmtleUEtLi0-ZGF0YVxuXG5wYXNzd29yZEItLT5oUGFzc3dvcmRCXG5wYXNzd29yZEItLT5wcml2YXRlS2V5QlxucGFzc3dvcmRCLS4tPnByaXZhdGVLZXlCXG5wdWJsaWNLZXlCLS0-a2V5QlxucHJpdmF0ZUtleUItLi0-a2V5Qlxua2V5Qi0tPmRhdGFcbmtleUItLi0-ZGF0YSIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoiZ3JhcGggTFJcblxuc3ViZ3JhcGggTGVnZW5kXG5cdEEgLS0gQSBpcyB1c2VkIHRvIGhhc2ggb3IgZW5jcnlwdCBCIC0tPiBCXG5cdEMgLS4gQyBpcyB1c2VkIHRvIGRlY3J5cHQgRCAuLT4gRFxuZW5kXG5cbnN1YmdyYXBoIFVzZXIgQSdzIE1pbmRcblx0cGFzc3dvcmRBKHBhc3N3b3JkKVxuZW5kXG5zdWJncmFwaCBVc2VyIEIncyBNaW5kXG5cdHBhc3N3b3JkQihwYXNzd29yZClcbmVuZFxuXG5zdWJncmFwaCBEYXRhYmFzZVxuXHRzdWJncmFwaCBVc2VyIEFcblx0XHRoUGFzc3dvcmRBKHBhc3N3b3JkKVxuXHRcdHB1YmxpY0tleUEocHVibGljIGtleSlcblx0XHRwcml2YXRlS2V5QShwcml2YXRlIGtleSlcblx0ZW5kXG5cdHN1YmdyYXBoIFVzZXIgQlxuXHRcdGhQYXNzd29yZEIocGFzc3dvcmQpXG5cdFx0cHVibGljS2V5QihwdWJsaWMga2V5KVxuXHRcdHByaXZhdGVLZXlCKHByaXZhdGUga2V5KVxuXHRlbmRcblx0c3ViZ3JhcGggRGF0dW0gT2JqZWN0XG5cdFx0a2V5QShrZXkpXG5cdFx0a2V5QihrZXkpXG5cdFx0ZGF0YShkYXRhKVxuXHRlbmRcbmVuZFxuXG5wYXNzd29yZEEtLT5oUGFzc3dvcmRBXG5wYXNzd29yZEEtLT5wcml2YXRlS2V5QVxucGFzc3dvcmRBLS4tPnByaXZhdGVLZXlBXG5wdWJsaWNLZXlBLS0-a2V5QVxucHJpdmF0ZUtleUEtLi0-a2V5QVxua2V5QS0tPmRhdGFcbmtleUEtLi0-ZGF0YVxuXG5wYXNzd29yZEItLT5oUGFzc3dvcmRCXG5wYXNzd29yZEItLT5wcml2YXRlS2V5QlxucGFzc3dvcmRCLS4tPnByaXZhdGVLZXlCXG5wdWJsaWNLZXlCLS0-a2V5QlxucHJpdmF0ZUtleUItLi0-a2V5Qlxua2V5Qi0tPmRhdGFcbmtleUItLi0-ZGF0YSIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)

### Design Choices

- The user's password salt and the hash algorithm used to hash their password are publicly available because these are needed by the client to hash the password before sending it to the server to authenticate the user.
- There are a number of advantages to having a public and private key associated with each user rather than using their password to encrypt the datum keys:
  - If the user changes their password, the only value that needs to be re-encrypted is their private key. If instead the user's password was used to encrypt datum keys, each datum key would need to be re-encrypted.
  - Having a public key and private key associated with the user allows for asymmetric encryption to be used for the datum keys, whereas using the user's password would require symmetric encryption. The advantage of symmetric encryption is that anyone can encrypt the value using the public key, which is an essential property to enable sharing. If user A wants to share a datum with user B, A simply needs to (on device) decrypt the datum key with their private key and then encrypt it with B's public key. Then, the B-encrypted datum key can be securely uploaded to the server, granting read access to B and B only.
  - Once the client gains authorization from the server, it no longer needs to store the user's password. Rather, it can store the user's private key in order to encrypt and decrypt passwords. This makes it slightly harder for malware on a client to steal the user's master password.

### Disadvantages

- All encryption happening client-side means that the server is unable to validate encrypted values that are sent to it. Thus, someone creating a client could corrupt data by mis-encrypting it, which would corrupt the data as the server would not know that the data is malformed. Although it would be nice to be able to validate the data, in my research, I found no ways of doing so while retaining the no-knowledge principle.
- Since all data is encrypted using the user's password, and the server has no way of knowing the user's password, there is no way for a user to recover their data. This is an intrinsic issue of the no-knowledge principle, and thus a necessary disadvantage to maintain security.

## APIs

The backend is available as either a gRPC or REST API. The REST API service is located at https://api.kript.us (on port `443`), and the gRPC API service is available at https://grpc.kript.us (also on port `443`). The folder [proto](proto) contains the proto definitions for the gRPC API. The REST API is created using the gRPC REST API gateway, so each gRPC function also defines an REST endpoint within the proto files. Additionally, there is a generated OpenAPI definition at [server/docs/api/kript.swagger.json](server/docs/api/kript.swagger.json).

## Cloud Architecture

The below graph provides an overview of the architecture:

[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBYKChDbGllbnQpKVxuICAgIHN1YmdyYXBoIEdvb2dsZSBDbG91ZCBQbGF0Zm9ybVxuXHQgICAgQVtnUlBDIEFQSSBTZXJ2aWNlXVxuXHQgICAgQltSRVNUIEFQSSBTZXJ2aWNlXVxuXHQgICAgQ1tBY2NvdW50IFNlcnZpY2VdXG5cdCAgICBEW0RhdGEgU2VydmljZV1cblx0ICAgIEUoQWNjb3VudCBEYXRhYmFzZSlcblx0ICAgIEYoRGF0YSBEYXRhYmFzZSlcbiAgICBlbmRcbiAgICBcbiAgICBYIC0uIGdycGMua3JpcHQudXMgLi0-IEFcbiAgICBYIC0tIGFwaS5rcmlwdC51cyAtLT4gQlxuICAgIEEgLS4tPiBDXG4gICAgQSAtLi0-IERcbiAgICBCIC0uLT4gQ1xuICAgIEIgLS4tPiBEXG4gICAgQyAtLi0-IEVcbiAgICBEIC0uLT4gRiIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBYKChDbGllbnQpKVxuICAgIHN1YmdyYXBoIEdvb2dsZSBDbG91ZCBQbGF0Zm9ybVxuXHQgICAgQVtnUlBDIEFQSSBTZXJ2aWNlXVxuXHQgICAgQltSRVNUIEFQSSBTZXJ2aWNlXVxuXHQgICAgQ1tBY2NvdW50IFNlcnZpY2VdXG5cdCAgICBEW0RhdGEgU2VydmljZV1cblx0ICAgIEUoQWNjb3VudCBEYXRhYmFzZSlcblx0ICAgIEYoRGF0YSBEYXRhYmFzZSlcbiAgICBlbmRcbiAgICBcbiAgICBYIC0uIGdycGMua3JpcHQudXMgLi0-IEFcbiAgICBYIC0tIGFwaS5rcmlwdC51cyAtLT4gQlxuICAgIEEgLS4tPiBDXG4gICAgQSAtLi0-IERcbiAgICBCIC0uLT4gQ1xuICAgIEIgLS4tPiBEXG4gICAgQyAtLi0-IEVcbiAgICBEIC0uLT4gRiIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)

Legend:
 - Dotted lines represent communication over gRPC.
 - Normal lines represent communication over https.
 - Each rectangular box is a Cloud Run service.
 - Each rounded rectangular box is a database.

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