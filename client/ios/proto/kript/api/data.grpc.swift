//
// DO NOT EDIT.
//
// Generated by the protocol buffer compiler.
// Source: kript/api/data.proto
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


/// Usage: instantiate Kript_Api_DataServiceClient, then call methods of this protocol to make API calls.
internal protocol Kript_Api_DataServiceClientProtocol: GRPCClient {
  func getData(
    _ request: Kript_Api_GetDataRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_GetDataRequest, Kript_Api_GetDataResponse>

  func updateDatum(
    _ request: Kript_Api_UpdateDatumRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_UpdateDatumRequest, Kript_Api_UpdateDatumResponse>

  func createDatum(
    _ request: Kript_Api_CreateDatumRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_CreateDatumRequest, Kript_Api_CreateDatumResponse>

  func deleteDatum(
    _ request: Kript_Api_DeleteDatumRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_DeleteDatumRequest, Kript_Api_DeleteDatumResponse>

  func shareDatum(
    _ request: Kript_Api_ShareDatumRequest,
    callOptions: CallOptions?
  ) -> UnaryCall<Kript_Api_ShareDatumRequest, Kript_Api_ShareDatumResponse>

}

extension Kript_Api_DataServiceClientProtocol {

  /// Get the list of all data for the logged in user, or a specific datum if
  /// specified.
  ///
  /// - Parameters:
  ///   - request: Request to send to GetData.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func getData(
    _ request: Kript_Api_GetDataRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_GetDataRequest, Kript_Api_GetDataResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.DataService/GetData",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Update the specified datum with new data.
  ///
  /// - Parameters:
  ///   - request: Request to send to UpdateDatum.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func updateDatum(
    _ request: Kript_Api_UpdateDatumRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_UpdateDatumRequest, Kript_Api_UpdateDatumResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.DataService/UpdateDatum",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Create a new datum.
  ///
  /// - Parameters:
  ///   - request: Request to send to CreateDatum.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func createDatum(
    _ request: Kript_Api_CreateDatumRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_CreateDatumRequest, Kript_Api_CreateDatumResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.DataService/CreateDatum",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Delete the specified datum.
  ///
  /// - Parameters:
  ///   - request: Request to send to DeleteDatum.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func deleteDatum(
    _ request: Kript_Api_DeleteDatumRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_DeleteDatumRequest, Kript_Api_DeleteDatumResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.DataService/DeleteDatum",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }

  /// Share a datum with another user by granting them new permission(s) on it.
  ///
  /// - Parameters:
  ///   - request: Request to send to ShareDatum.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func shareDatum(
    _ request: Kript_Api_ShareDatumRequest,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Kript_Api_ShareDatumRequest, Kript_Api_ShareDatumResponse> {
    return self.makeUnaryCall(
      path: "/kript.api.DataService/ShareDatum",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions
    )
  }
}

internal final class Kript_Api_DataServiceClient: Kript_Api_DataServiceClientProtocol {
  internal let channel: GRPCChannel
  internal var defaultCallOptions: CallOptions

  /// Creates a client for the kript.api.DataService service.
  ///
  /// - Parameters:
  ///   - channel: `GRPCChannel` to the service host.
  ///   - defaultCallOptions: Options to use for each service call if the user doesn't provide them.
  internal init(channel: GRPCChannel, defaultCallOptions: CallOptions = CallOptions()) {
    self.channel = channel
    self.defaultCallOptions = defaultCallOptions
  }
}
