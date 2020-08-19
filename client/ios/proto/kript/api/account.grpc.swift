//
// DO NOT EDIT.
//
// Generated by the protocol buffer compiler.
// Source: kript/api/account.proto
//

//
// Copyright 2018, gRPC Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
import Foundation
import GRPC
import NIO
import NIOHTTP1
import SwiftProtobuf


/// Usage: instantiate Kript_Api_AccountServiceClient, then call methods of this protocol to make API calls.
internal protocol Kript_Api_AccountServiceClientProtocol: GRPCClient {
  func loginUser(
    _ request: Kript_Api_LoginUserRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_LoginUserRequest, Kript_Api_LoginUserResponse>

  func sendVerification(
    _ request: Kript_Api_SendVerificationRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_SendVerificationRequest, Kript_Api_SendVerificationResponse>

  func verifyUser(
    _ request: Kript_Api_VerifyUserRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_VerifyUserRequest, Kript_Api_VerifyUserResponse>

  func updatePassword(
    _ request: Kript_Api_UpdatePasswordRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_UpdatePasswordRequest, Kript_Api_UpdatePasswordResponse>

  func createAccount(
    _ request: Kript_Api_CreateAccountRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_CreateAccountRequest, Kript_Api_CreateAccountResponse>

  func refreshAuth(
    _ request: Kript_Api_RefreshAuthRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_RefreshAuthRequest, Kript_Api_RefreshAuthResponse>

  func getUser(
    _ request: Kript_Api_GetUserRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_GetUserRequest, Kript_Api_GetUserResponse>

  func addTwoFactor(
    _ request: Kript_Api_AddTwoFactorRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_AddTwoFactorRequest, Kript_Api_AddTwoFactorResponse>

  func verifyTwoFactor(
    _ request: Kript_Api_VerifyTwoFactorRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_VerifyTwoFactorRequest, Kript_Api_VerifyTwoFactorResponse>

}

extension Kript_Api_AccountServiceClientProtocol {

  /// Login the user. If the user has 2-factor authentication enabled,
  /// a verification code must be sent with SendVerification to complete the
  /// login process.
  ///
  /// - Parameters:
  ///   - request: Request to send to LoginUser.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func loginUser(
    _ request: Kript_Api_LoginUserRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_LoginUserRequest, Kript_Api_LoginUserResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.AccountService/LoginUser",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Send a verification code to the user using the specified method.
  ///
  /// - Parameters:
  ///   - request: Request to send to SendVerification.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func sendVerification(
    _ request: Kript_Api_SendVerificationRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_SendVerificationRequest, Kript_Api_SendVerificationResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.AccountService/SendVerification",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Complete logging in the user.
  ///
  /// - Parameters:
  ///   - request: Request to send to VerifyUser.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func verifyUser(
    _ request: Kript_Api_VerifyUserRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_VerifyUserRequest, Kript_Api_VerifyUserResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.AccountService/VerifyUser",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Change the user's password.
  ///
  /// - Parameters:
  ///   - request: Request to send to UpdatePassword.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func updatePassword(
    _ request: Kript_Api_UpdatePasswordRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_UpdatePasswordRequest, Kript_Api_UpdatePasswordResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.AccountService/UpdatePassword",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Create an account.
  ///
  /// - Parameters:
  ///   - request: Request to send to CreateAccount.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func createAccount(
    _ request: Kript_Api_CreateAccountRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_CreateAccountRequest, Kript_Api_CreateAccountResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.AccountService/CreateAccount",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Fetch a new access token.
  ///
  /// - Parameters:
  ///   - request: Request to send to RefreshAuth.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func refreshAuth(
    _ request: Kript_Api_RefreshAuthRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_RefreshAuthRequest, Kript_Api_RefreshAuthResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.AccountService/RefreshAuth",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Get the information of the user with the given username or user id.
  /// If the user is the logged in user, the private user information is
  /// included.
  ///
  /// - Parameters:
  ///   - request: Request to send to GetUser.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func getUser(
    _ request: Kript_Api_GetUserRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_GetUserRequest, Kript_Api_GetUserResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.AccountService/GetUser",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Request to add the given two-factor destination and send a confirmation
  /// code to the two-factor destination.
  ///
  /// - Parameters:
  ///   - request: Request to send to AddTwoFactor.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func addTwoFactor(
    _ request: Kript_Api_AddTwoFactorRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_AddTwoFactorRequest, Kript_Api_AddTwoFactorResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.AccountService/AddTwoFactor",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Verify a two-factor destination.
  ///
  /// - Parameters:
  ///   - request: Request to send to VerifyTwoFactor.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func verifyTwoFactor(
    _ request: Kript_Api_VerifyTwoFactorRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_VerifyTwoFactorRequest, Kript_Api_VerifyTwoFactorResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.AccountService/VerifyTwoFactor",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }
}

internal final class Kript_Api_AccountServiceClient: Kript_Api_AccountServiceClientProtocol {
  internal let channel: GRPCChannel
  internal var defaultCallOptions: CallOptions

  /// Creates a client for the kript.api.AccountService service.
  ///
  /// - Parameters:
  ///   - channel: `GRPCChannel` to the service host.
  ///   - defaultCallOptions: Options to use for each service call if the user doesn't provide them.
  internal init(channel: GRPCChannel, defaultCallOptions: CallOptions = CallOptions()) {
    self.channel = channel
    self.defaultCallOptions = defaultCallOptions
  }
}

