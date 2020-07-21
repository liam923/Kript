package data

import (
	"context"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (s *Server) GetData(ctx context.Context, request *api.GetDataRequest) (*api.GetDataResponse, error) {
	// validate the request and authenticate the user before doing anything
	if request == nil || !s.validateAccessTokenFormat(request.AccessToken) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	userId, err := s.loginUserWithAccessToken(*request.AccessToken)
	if err != nil {
		return nil, err
	}

	var data []*api.Datum
	if request.DatumIds == nil || len(request.DatumIds) == 0 {
		// fetch all datums that the user can read
		fData, err := s.database.fetchDataForUser(ctx, userId)
		if err != nil {
			return nil, err
		}
		data = make([]*api.Datum, len(*fData))
		for i, idedDatum := range *fData {
			data[i] = idedDatum.Datum.toApiDatum(idedDatum.Id)
		}
	} else {
		data = make([]*api.Datum, len(request.DatumIds))
		for i, id := range request.DatumIds {
			datum, err := s.database.fetchDatum(ctx, id)
			if err != nil {
				return nil, err
			}

			// check that the user has read permission on the datum
			if confirmPermission(userId, datum, newRawPermissionSet(rawPermissionRead)) {
				data[i] = datum.toApiDatum(id)
			} else {
				return nil, status.Errorf(codes.PermissionDenied, "access denied for datum %s", id)
			}
		}
	}

	return &api.GetDataResponse{
		Datums: data,
	}, nil
}

func (s *Server) UpdateDatum(ctx context.Context, request *api.UpdateDatumRequest) (*api.UpdateDatumResponse, error) {
	startTime := time.Now()

	// validate the request and authenticate the user before doing anything
	if request == nil || request.Data == nil || !s.validateAccessTokenFormat(request.AccessToken) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	userId, err := s.loginUserWithAccessToken(*request.AccessToken)
	if err != nil {
		return nil, err
	}

	// fetch the datum as it is prior to updating
	oldDatum, err := s.database.fetchDatum(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	// confirm the user has write access
	if !confirmPermission(userId, oldDatum, newRawPermissionSet(rawPermissionRead, rawPermissionWrite)) {
		return nil, status.Errorf(codes.PermissionDenied, "write access denied for datum %s", request.Id)
	}

	newDatum := *oldDatum
	newDatum.Title = request.Title
	newDatum.Data = request.Data.Data
	newDatum.DataIv = request.DataIv
	newDatum.DataEncryptionAlgorithm = request.DataEncryptionAlgorithm
	newDatum.Metadata.LastEdited = startTime

	err = s.database.updateDatum(ctx, &newDatum, request.Id)
	if err != nil {
		return nil, err
	}

	return &api.UpdateDatumResponse{
		Datum: newDatum.toApiDatum(request.Id),
	}, nil
}

func (s *Server) CreateDatum(ctx context.Context, request *api.CreateDatumRequest) (*api.CreateDatumResponse, error) {
	// validate the request and authenticate the user before doing anything
	if request == nil || request.Data == nil || request.DataKey == nil ||
		!s.validateAccessTokenFormat(request.AccessToken) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	userId, err := s.loginUserWithAccessToken(*request.AccessToken)
	if err != nil {
		return nil, err
	}

	newDatum := datum{
		Owner:                   userId,
		Title:                   request.Title,
		Data:                    request.Data.Data,
		DataEncryptionAlgorithm: request.DataEncryptionAlgorithm,
		DataIv:                  request.DataIv,
		Accessors: map[string]accessor{
			userId: {
				UserId:      userId,
				DataKey:     request.DataKey.Data,
				Permissions: []api.Permission{api.Permission_ADMIN},
			},
		},
		Metadata: metadata{
			OwnerId:     userId,
			CreatedTime: time.Now(),
			LastEdited:  time.Now(),
			AccessMetadata: map[string]accessMetadata{
				userId: {
					GrantMetadata: []permissionGrantMetadata{
						{
							GranterId:  userId,
							Permission: api.Permission_ADMIN,
							IsGrant:    true,
						},
					},
				},
			},
		},
	}

	id, err := s.database.createDatum(ctx, &newDatum)
	if err != nil {
		return nil, err
	}

	return &api.CreateDatumResponse{
		Datum: newDatum.toApiDatum(id),
	}, nil
}

func (s *Server) DeleteDatum(ctx context.Context, request *api.DeleteDatumRequest) (*api.DeleteDatumResponse, error) {
	// validate the request and authenticate the user before doing anything
	if request == nil || !s.validateAccessTokenFormat(request.AccessToken) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	userId, err := s.loginUserWithAccessToken(*request.AccessToken)
	if err != nil {
		return nil, err
	}

	// fetch the datum as it is before deleting in order to verify the user has delete access
	datum, err := s.database.fetchDatum(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if !confirmPermission(userId, datum, newRawPermissionSet(rawPermissionRead, rawPermissionDelete)) {
		return nil, status.Errorf(codes.PermissionDenied, "delete access denied for datum %s", request.Id)
	}

	err = s.database.deleteDatum(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return &api.DeleteDatumResponse{
		Datum: datum.toApiDatum(request.Id),
	}, nil
}

func (s *Server) ShareDatum(ctx context.Context, request *api.ShareDatumRequest) (*api.ShareDatumResponse, error) {
	// validate the request and authenticate the user before doing anything
	if request == nil || request.DataKey == nil || !s.validateAccessTokenFormat(request.AccessToken) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	userId, err := s.loginUserWithAccessToken(*request.AccessToken)
	if err != nil {
		return nil, err
	}

	// if there is an unknown permission, throw an error
	if contains(request.Permissions, api.Permission_UNKNOWN) {
		return nil, status.Errorf(codes.InvalidArgument, "unable to share \"UNKNOWN\" permission")
	}

	// fetch the datum as it is before sharing
	oldDatum, err := s.database.fetchDatum(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	// validate that the user has permission to share the datum
	if !confirmPermission(userId, oldDatum, newRawPermissionSet(rawPermissionRead, rawPermissionShare)) {
		return nil, status.Errorf(codes.PermissionDenied, "share access denied for datum %s", request.Id)
	}
	// validate that the user has all permissions that they are trying to share
	allRawPermissionsGranted := newRawPermissionSet()
	for _, apiPermission := range request.Permissions {
		for rawPermission := range apiPermissionToRaw[apiPermission].m {
			allRawPermissionsGranted.m[rawPermission] = struct{}{}
		}
	}
	if !confirmPermission(userId, oldDatum, allRawPermissionsGranted) {
		return nil, status.Errorf(codes.PermissionDenied,
			"attempted to share more access than the user has for datum %s",
			request.Id)
	}

	newDatum := *oldDatum
	// modify the access data
	var newAccessor accessor
	if oldAccessor, ok := newDatum.Accessors[request.TargetId]; ok {
		newAccessor = oldAccessor
		// append new permissions to list of permissions
		for _, permission := range request.Permissions {
			if !contains(newAccessor.Permissions, permission) {
				newAccessor.Permissions = append(newAccessor.Permissions, permission)
			}
		}
	} else {
		newAccessor = accessor{
			UserId:      request.TargetId,
			DataKey:     request.DataKey.Data,
			Permissions: request.Permissions,
		}
	}
	if newDatum.Accessors == nil {
		newDatum.Accessors = map[string]accessor{}
	}
	newDatum.Accessors[request.TargetId] = newAccessor

	var newAccessMetadata accessMetadata
	if oldAccessMetadata, ok := newDatum.Metadata.AccessMetadata[request.TargetId]; ok {
		newAccessMetadata = oldAccessMetadata
	} else {
		newAccessMetadata = accessMetadata{
			GrantMetadata: []permissionGrantMetadata{},
		}
	}
	for _, permission := range request.Permissions {
		newAccessMetadata.GrantMetadata = append(newAccessMetadata.GrantMetadata, permissionGrantMetadata{
			GranterId:  userId,
			Permission: permission,
			IsGrant:    true,
		})
	}
	if newDatum.Metadata.AccessMetadata == nil {
		newDatum.Metadata.AccessMetadata = map[string]accessMetadata{}
	}
	newDatum.Metadata.AccessMetadata[request.TargetId] = newAccessMetadata

	err = s.database.updateDatum(ctx, &newDatum, request.Id)
	if err != nil {
		return nil, err
	}

	return &api.ShareDatumResponse{
		Datum: newDatum.toApiDatum(request.Id),
	}, nil
}
