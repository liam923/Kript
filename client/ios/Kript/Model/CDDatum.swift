//
//  CDDatum.swift
//  Kript
//
//  Created by Liam Stevenson on 8/2/20.
//  Copyright © 2020 Liam Stevenson. All rights reserved.
//

extension CDDatum {
    func getProto() throws -> Kript_Api_Datum? {
        if let data = self.protoData {
            return try Kript_Api_Datum(serializedData: data)
        } else {
            return nil
        }
    }
    
    func setProto(fromDatum datum: Kript_Api_Datum) throws {
        let data = try datum.serializedData()
        self.protoData = data
    }
}
