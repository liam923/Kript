package data

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/liam923/Kript/server/internal/generate"
	"github.com/liam923/Kript/server/internal/jwt"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"testing"
)

// Create a server that writes to a dummy log.
func createServer(t *testing.T, db database) (Server, jwt.Signer) {
	issuerId := "kript.api"
	keyPair := generate.Keys(4096)
	logger := grpclog.NewLoggerV2(&dummyWriter{}, &dummyWriter{}, &dummyWriter{})
	signer, err := jwt.NewSigner(keyPair.Private, issuerId)
	if err != nil {
		t.Errorf("Failed to initialize signer")
	}
	validator, err := jwt.NewValidator(keyPair.Public, issuerId)
	if err != nil {
		t.Errorf("Failed to initialize validator")
	}
	return Server{
		database:  db,
		Logger:    &logger,
		validator: validator,
	}, *signer
}

type dummyWriter struct{}

func (w *dummyWriter) Write([]byte) (n int, err error) { return }

func TestGetData(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server, signer := createServer(t, db)

	validToken, invalidTokens := generate.JWT(&signer, "liam923", jwt.AccessTokenType)

	type fetchForUser struct {
		user string
		data []idedDatum
		err  error
	}
	type fetchData struct {
		id    string
		datum datum
		err   error
	}
	tests := []struct {
		testName     string
		request      []string
		err          error
		fetchForUser *fetchForUser
		fetchData    *[]fetchData
	}{
		{
			testName: "valid individual",
			request:  []string{"1", "2"},
			err:      nil,
			fetchData: &[]fetchData{
				{
					id: "1",
					datum: datum{
						Owner:                   "liam923",
						Data:                    []byte("23dyoirbeu9"),
						DataEncryptionAlgorithm: 0,
						Accessors: map[string]accessor{
							"liam923": {
								UserId:      "liam923",
								DataKey:     []byte("a0iojcdopi"),
								Permissions: []api.Permission{api.Permission_ADMIN},
							},
						},
						Metadata: metadata{},
					},
					err: nil,
				},
				{
					id: "2",
					datum: datum{
						Owner:                   "other person",
						Data:                    []byte("saofinfwuoirc"),
						DataEncryptionAlgorithm: 0,
						Accessors: map[string]accessor{
							"other person": {
								UserId:      "other person",
								DataKey:     []byte("23ewionpwje"),
								Permissions: []api.Permission{api.Permission_ADMIN},
							},
							"liam923": {
								UserId:      "liam923",
								DataKey:     []byte("23idjoqnwuiacjs"),
								Permissions: []api.Permission{api.Permission_READ},
							},
						},
						Metadata: metadata{},
					},
					err: nil,
				},
			},
		},
		{
			testName: "valid all",
			request:  []string{},
			err:      nil,
			fetchForUser: &fetchForUser{
				user: "liam923",
				data: []idedDatum{
					{
						Id: "1",
						Datum: datum{
							Owner:                   "liam923",
							Data:                    []byte("23dyoirbeu9"),
							DataEncryptionAlgorithm: 0,
							Accessors: map[string]accessor{
								"liam923": {
									UserId:      "liam923",
									DataKey:     []byte("a0iojcdopi"),
									Permissions: []api.Permission{api.Permission_ADMIN},
								},
							},
							Metadata: metadata{},
						},
					},
					{
						Id: "2",
						Datum: datum{
							Owner:                   "other person",
							Data:                    []byte("saofinfwuoirc"),
							DataEncryptionAlgorithm: 0,
							Accessors: map[string]accessor{
								"other person": {
									UserId:      "other person",
									DataKey:     []byte("23ewionpwje"),
									Permissions: []api.Permission{api.Permission_ADMIN},
								},
								"liam923": {
									UserId:      "liam923",
									DataKey:     []byte("23idjoqnwuiacjs"),
									Permissions: []api.Permission{api.Permission_READ},
								},
								"third person": {
									UserId:      "third person",
									DataKey:     []byte("29huawosincd"),
									Permissions: []api.Permission{api.Permission_SHARE},
								},
							},
							Metadata: metadata{},
						},
					},
				},
				err: nil,
			},
		},
		{
			testName: "invalid individual",
			request:  []string{"1", "2"},
			err:      status.Error(codes.Internal, "error"),
			fetchData: &[]fetchData{
				{
					id: "1",
					datum: datum{
						Owner:                   "liam923",
						Data:                    []byte("23dyoirbeu9"),
						DataEncryptionAlgorithm: 0,
						Accessors: map[string]accessor{
							"liam923": {
								UserId:      "liam923",
								DataKey:     []byte("a0iojcdopi"),
								Permissions: []api.Permission{api.Permission_ADMIN},
							},
						},
						Metadata: metadata{},
					},
					err: nil,
				},
				{
					id:  "2",
					err: status.Error(codes.Internal, "error"),
				},
			},
		},
		{
			testName: "invalid all",
			request:  []string{},
			err:      status.Error(codes.Internal, "error"),
			fetchForUser: &fetchForUser{
				user: "liam923",
				data: nil,
				err:  status.Error(codes.Internal, "error"),
			},
		},
		{
			testName: "unauthorized individual",
			request:  []string{"1", "2"},
			err:      status.Errorf(codes.PermissionDenied, "access denied for datum 2"),
			fetchData: &[]fetchData{
				{
					id: "1",
					datum: datum{
						Owner:                   "liam923",
						Data:                    []byte("23dyoirbeu9"),
						DataEncryptionAlgorithm: 0,
						Accessors: map[string]accessor{
							"liam923": {
								UserId:      "liam923",
								DataKey:     []byte("a0iojcdopi"),
								Permissions: []api.Permission{api.Permission_ADMIN},
							},
						},
						Metadata: metadata{},
					},
					err: nil,
				},
				{
					id: "2",
					datum: datum{
						Owner:                   "other person",
						Data:                    []byte("saofinfwuoirc"),
						DataEncryptionAlgorithm: 0,
						Accessors: map[string]accessor{
							"third party": {
								UserId:      "third party",
								DataKey:     []byte("23idjoqnwuiacjs"),
								Permissions: []api.Permission{api.Permission_READ},
							},
							"other person": {
								UserId:      "other person",
								DataKey:     []byte("23ewionpwje"),
								Permissions: []api.Permission{api.Permission_ADMIN},
							},
						},
						Metadata: metadata{},
					},
					err: nil,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf(tt.testName, i), func(t *testing.T) {
			data := []idedDatum{}
			if tt.fetchForUser != nil {
				db.
					EXPECT().
					fetchDataForUser(context.Background(), tt.fetchForUser.user).
					Return(&tt.fetchForUser.data, tt.fetchForUser.err)
				data = tt.fetchForUser.data
			} else if tt.fetchData != nil {
				for _, fetch := range *tt.fetchData {
					d := fetch.datum
					db.EXPECT().
						fetchDatum(context.Background(), fetch.id).
						Return(&d, fetch.err)
					data = append(data, idedDatum{
						Datum: d,
						Id:    fetch.id,
					})
				}
			}

			request := &api.GetDataRequest{
				AccessToken: &api.AccessToken{
					Jwt: &api.JWT{Token: validToken},
				},
				DatumIds: tt.request,
			}
			response, err := server.GetData(context.Background(), request)
			if tt.err != nil && (err == nil || tt.err.Error() != err.Error()) {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.err == nil {
				if response == nil {
					t.Errorf("unexpected response: %v", response)
				}
				for i, actual := range response.Datums {
					expectedId := data[i].Id
					expected := data[i].Datum
					if expectedId != actual.Id ||
						expected.Owner != actual.Owner ||
						bytes.Compare(expected.Data, actual.Data.Data) != 0 ||
						expected.DataEncryptionAlgorithm != actual.DataEncryptionAlgorithm ||
						len(expected.Accessors) != len(actual.Accessors) {
						t.Errorf("unexpected response: %v", response)
					}
				}
			}

			for _, invalidToken := range invalidTokens {
				request = &api.GetDataRequest{
					AccessToken: &api.AccessToken{
						Jwt: &api.JWT{Token: invalidToken},
					},
					DatumIds: tt.request,
				}
				response, err := server.GetData(context.Background(), request)
				if err == nil {
					t.Error("nil error for invalid token")
				} else if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
					t.Errorf("invalid error for invalid token: %s", err)
				}
				if response != nil {
					t.Errorf("response should have been nil for invalid token: %v", response)
				}
			}
		})
	}

	invalidRequests := []*api.GetDataRequest{
		nil,
		{
			AccessToken: nil,
			DatumIds:    nil,
		},
	}
	t.Run("invalid requests", func(t *testing.T) {
		for _, request := range invalidRequests {
			response, err := server.GetData(context.Background(), request)
			if response != nil {
				t.Errorf("expected response to be nil: %v", response)
			}
			if err == nil {
				t.Errorf("expected err not to be nil")
			}
			if s, ok := status.FromError(err); !ok || s.Code() != codes.InvalidArgument {
				t.Errorf("unexpected err: %v", err)
			}
		}
	})
}

type partialDatumMatcher struct {
	datum datum
}

func (p partialDatumMatcher) Matches(x interface{}) bool {
	if q, ok := x.(*datum); ok {
		return bytes.Compare(p.datum.Data, q.Data) == 0 &&
			len(p.datum.Accessors) == len(q.Accessors) &&
			p.datum.Owner == q.Owner
	}
	return false
}

func (p partialDatumMatcher) String() string {
	return fmt.Sprintf("%v", p.datum)
}

func TestUpdateDatum(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server, signer := createServer(t, db)

	validToken, invalidTokens := generate.JWT(&signer, "liam923", jwt.AccessTokenType)

	type fetchDatum struct {
		id    string
		datum datum
		err   error
	}
	tests := []struct {
		testName   string
		request    api.UpdateDatumRequest
		fetchDatum fetchDatum
		updateErr  error
		err        error
	}{
		{
			testName: "valid",
			request: api.UpdateDatumRequest{
				Id:     "1",
				Title:  "new title",
				Data:   &api.ESecret{Data: []byte("new data")},
				DataIv: []byte("iodnjewq"),
			},
			err: nil,
			fetchDatum: fetchDatum{
				id: "1",
				datum: datum{
					Owner:                   "liam923",
					Data:                    []byte("23dyoirbeu9"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"liam923": {
							UserId:      "liam923",
							DataKey:     []byte("a0iojcdopi"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
		},
		{
			testName: "invalid id",
			request: api.UpdateDatumRequest{
				Id:     "1",
				Title:  "new title",
				Data:   &api.ESecret{Data: []byte("new data")},
				DataIv: []byte("iodnjewq"),
			},
			err: status.Error(codes.Internal, "error"),
			fetchDatum: fetchDatum{
				id:    "1",
				datum: datum{},
				err:   status.Error(codes.Internal, "error"),
			},
		},
		{
			testName: "unauthorized write",
			request: api.UpdateDatumRequest{
				Id:     "2",
				Title:  "new title",
				Data:   &api.ESecret{Data: []byte("new data")},
				DataIv: []byte("iodnjewq"),
			},
			err: status.Errorf(codes.PermissionDenied, "write access denied for datum 2"),
			fetchDatum: fetchDatum{
				id: "2",
				datum: datum{
					Owner:                   "other person",
					Data:                    []byte("saofinfwuoirc"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"third party": {
							UserId:      "third party",
							DataKey:     []byte("23idjoqnwuiacjs"),
							Permissions: []api.Permission{api.Permission_READ},
						},
						"other person": {
							UserId:      "other person",
							DataKey:     []byte("23ewionpwje"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
		},
		{
			testName: "write error",
			request: api.UpdateDatumRequest{
				Id:     "2",
				Title:  "new title",
				Data:   &api.ESecret{Data: []byte("new data")},
				DataIv: []byte("iodnjewq"),
			},
			err: status.Errorf(codes.PermissionDenied, "write access denied for datum 2"),
			fetchDatum: fetchDatum{
				id: "2",
				datum: datum{
					Owner:                   "liam923",
					Data:                    []byte("saofinfwuoirc"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"liam923": {
							UserId:      "liam923",
							DataKey:     []byte("a0iojcdopi"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
						"third party": {
							UserId:      "third party",
							DataKey:     []byte("23idjoqnwuiacjs"),
							Permissions: []api.Permission{api.Permission_READ},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
			updateErr: status.Errorf(codes.PermissionDenied, "write access denied for datum 2"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf(tt.testName, i), func(t *testing.T) {
			d := tt.fetchDatum.datum
			db.EXPECT().
				fetchDatum(context.Background(), tt.fetchDatum.id).
				Return(&d, tt.fetchDatum.err)

			if tt.fetchDatum.err == nil {
				expectedDatum := tt.fetchDatum.datum
				expectedDatum.Data = tt.request.Data.Data
				expectedDatum.DataIv = tt.request.DataIv
				db.EXPECT().
					updateDatum(context.Background(), partialDatumMatcher{datum: expectedDatum}, tt.fetchDatum.id).
					Return(tt.updateErr)
			}

			request := tt.request
			request.AccessToken = &api.AccessToken{
				Jwt: &api.JWT{Token: validToken},
			}

			response, err := server.UpdateDatum(context.Background(), &request)
			if tt.err != nil && (err == nil || tt.err.Error() != err.Error()) {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.err == nil {
				if response == nil {
					t.Errorf("unexpected response: %v", response)
				}
				expectedId := tt.fetchDatum.id
				expected := tt.fetchDatum.datum
				actual := response.Datum
				if expectedId != actual.Id ||
					expected.Owner != actual.Owner ||
					bytes.Compare(tt.request.Data.Data, actual.Data.Data) != 0 ||
					bytes.Compare(tt.request.DataIv, actual.DataIv) != 0 ||
					len(expected.Accessors) != len(actual.Accessors) {
					t.Errorf("unexpected response: %v, expected: %v", response, expected)
				}
			}

			for _, invalidToken := range invalidTokens {
				request := tt.request
				request.AccessToken = &api.AccessToken{
					Jwt: &api.JWT{Token: invalidToken},
				}
				response, err := server.UpdateDatum(context.Background(), &request)
				if err == nil {
					t.Error("nil error for invalid token")
				} else if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
					t.Errorf("invalid error for invalid token: %s", err)
				}
				if response != nil {
					t.Errorf("response should have been nil for invalid token: %v", response)
				}
			}
		})
	}

	invalidRequests := []*api.UpdateDatumRequest{
		nil,
		{
			AccessToken: nil,
		},
		{
			AccessToken: &api.AccessToken{Jwt: &api.JWT{Token: ""}},
		},
	}
	t.Run("invalid requests", func(t *testing.T) {
		for _, request := range invalidRequests {
			response, err := server.UpdateDatum(context.Background(), request)
			if response != nil {
				t.Errorf("expected response to be nil: %v", response)
			}
			if err == nil {
				t.Errorf("expected err not to be nil")
			}
			if s, ok := status.FromError(err); !ok || s.Code() != codes.InvalidArgument {
				t.Errorf("unexpected err: %v", err)
			}
		}
	})
}

func TestCreateDatum(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server, signer := createServer(t, db)

	validToken, invalidTokens := generate.JWT(&signer, "liam923", jwt.AccessTokenType)

	type createResponse struct {
		id  string
		err error
	}
	tests := []struct {
		testName       string
		datum          datum
		createResponse createResponse
	}{
		{
			testName: "valid",
			datum: datum{
				Owner:                   "liam923",
				Data:                    []byte("DATA"),
				DataEncryptionAlgorithm: 0,
				DataIv:                  []byte("suh dude"),
				Accessors: map[string]accessor{
					"liam923": {
						UserId:      "liam923",
						DataKey:     []byte("DATA KEY"),
						Permissions: []api.Permission{api.Permission_ADMIN},
					},
				},
				Metadata: metadata{},
			},
			createResponse: createResponse{
				id: "1",
			},
		},
		{
			testName: "create error",
			datum: datum{
				Owner:                   "liam923",
				Data:                    []byte("DATA"),
				DataEncryptionAlgorithm: 0,
				DataIv:                  []byte("suh dude"),
				Accessors: map[string]accessor{
					"liam923": {
						UserId:      "liam923",
						DataKey:     []byte("DATA KEY"),
						Permissions: []api.Permission{api.Permission_ADMIN},
					},
				},
				Metadata: metadata{},
			},
			createResponse: createResponse{
				err: status.Errorf(codes.Internal, "failed to create"),
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf(tt.testName, i), func(t *testing.T) {
			db.EXPECT().
				createDatum(context.Background(), partialDatumMatcher{datum: tt.datum}).
				Return(tt.createResponse.id, tt.createResponse.err)

			request := api.CreateDatumRequest{
				AccessToken: &api.AccessToken{
					Jwt: &api.JWT{Token: validToken},
				},
				Data:                    &api.ESecret{Data: tt.datum.Data},
				DataKey:                 &api.EBytes{Data: tt.datum.Accessors["liam923"].DataKey},
				DataEncryptionAlgorithm: tt.datum.DataEncryptionAlgorithm,
				DataIv:                  tt.datum.DataIv,
			}

			response, err := server.CreateDatum(context.Background(), &request)
			if tt.createResponse.err != nil && (err == nil || tt.createResponse.err.Error() != err.Error()) {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.createResponse.err == nil {
				if response == nil {
					t.Errorf("unexpected response: %v, got err: %v", response, err)
				}
				expectedId := tt.createResponse.id
				expected := tt.datum
				actual := response.Datum
				if expectedId != actual.Id ||
					expected.Owner != actual.Owner ||
					bytes.Compare(expected.Data, actual.Data.Data) != 0 ||
					bytes.Compare(expected.DataIv, actual.DataIv) != 0 ||
					expected.DataEncryptionAlgorithm != actual.DataEncryptionAlgorithm ||
					len(expected.Accessors) != len(actual.Accessors) {
					t.Errorf("unexpected response: %v", response)
				}
			}

			for _, invalidToken := range invalidTokens {
				request := api.CreateDatumRequest{
					AccessToken: &api.AccessToken{
						Jwt: &api.JWT{Token: invalidToken},
					},
					Data:                    &api.ESecret{Data: tt.datum.Data},
					DataKey:                 &api.EBytes{Data: tt.datum.Accessors["liam923"].DataKey},
					DataEncryptionAlgorithm: tt.datum.DataEncryptionAlgorithm,
					DataIv:                  tt.datum.DataIv,
				}
				response, err := server.CreateDatum(context.Background(), &request)
				if err == nil {
					t.Error("nil error for invalid token")
				} else if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
					t.Errorf("invalid error for invalid token: %s", err)
				}
				if response != nil {
					t.Errorf("response should have been nil for invalid token: %v", response)
				}
			}
		})
	}

	invalidRequests := []*api.CreateDatumRequest{
		nil,
		{
			AccessToken: nil,
		},
		{
			AccessToken: &api.AccessToken{Jwt: &api.JWT{Token: ""}},
			DataKey:     nil,
			Data:        nil,
		},
		{
			AccessToken: &api.AccessToken{Jwt: &api.JWT{Token: ""}},
			DataKey:     &api.EBytes{},
			Data:        nil,
		},
		{
			AccessToken: &api.AccessToken{Jwt: &api.JWT{Token: ""}},
			DataKey:     nil,
			Data:        &api.ESecret{},
		},
	}
	t.Run("invalid requests", func(t *testing.T) {
		for _, request := range invalidRequests {
			response, err := server.CreateDatum(context.Background(), request)
			if response != nil {
				t.Errorf("expected response to be nil: %v", response)
			}
			if err == nil {
				t.Errorf("expected err not to be nil")
			}
			if s, ok := status.FromError(err); !ok || s.Code() != codes.InvalidArgument {
				t.Errorf("unexpected err: %v", err)
			}
		}
	})
}

func TestDeleteDatum(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server, signer := createServer(t, db)

	validToken, invalidTokens := generate.JWT(&signer, "liam923", jwt.AccessTokenType)

	type fetchDatum struct {
		id    string
		datum datum
		err   error
	}
	tests := []struct {
		testName   string
		request    api.DeleteDatumRequest
		err        error
		fetchDatum fetchDatum
	}{
		{
			testName: "valid",
			request:  api.DeleteDatumRequest{Id: "1"},
			err:      nil,
			fetchDatum: fetchDatum{
				id: "1",
				datum: datum{
					Owner:                   "liam923",
					Data:                    []byte("23dyoirbeu9"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"liam923": {
							UserId:      "liam923",
							DataKey:     []byte("a0iojcdopi"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
		},
		{
			testName: "invalid fetch",
			request:  api.DeleteDatumRequest{Id: "2"},
			err:      status.Error(codes.Internal, "error"),
			fetchDatum: fetchDatum{
				id:  "2",
				err: status.Error(codes.Internal, "error"),
			},
		},
		{
			testName: "no permission",
			request:  api.DeleteDatumRequest{Id: "1"},
			err:      status.Error(codes.PermissionDenied, "delete access denied for datum 1"),
			fetchDatum: fetchDatum{
				id: "1",
				datum: datum{
					Owner:                   "other person",
					Data:                    []byte("saofinfwuoirc"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"third party": {
							UserId:      "third party",
							DataKey:     []byte("23idjoqnwuiacjs"),
							Permissions: []api.Permission{api.Permission_READ},
						},
						"other person": {
							UserId:      "other person",
							DataKey:     []byte("23ewionpwje"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf(tt.testName, i), func(t *testing.T) {
			d := tt.fetchDatum.datum
			db.EXPECT().
				fetchDatum(context.Background(), tt.fetchDatum.id).
				Return(&d, tt.fetchDatum.err)

			if tt.fetchDatum.err == nil {
				db.EXPECT().
					deleteDatum(context.Background(), tt.fetchDatum.id).
					Return(tt.fetchDatum.err)
			}

			request := tt.request
			request.AccessToken = &api.AccessToken{
				Jwt: &api.JWT{Token: validToken},
			}

			response, err := server.DeleteDatum(context.Background(), &request)
			if tt.err != nil && (err == nil || tt.err.Error() != err.Error()) {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.err == nil {
				if response == nil {
					t.Errorf("unexpected response: %v", response)
				}
				expectedId := tt.fetchDatum.id
				expected := tt.fetchDatum.datum
				actual := response.Datum
				if expectedId != actual.Id ||
					expected.Owner != actual.Owner ||
					bytes.Compare(expected.Data, actual.Data.Data) != 0 ||
					expected.DataEncryptionAlgorithm != actual.DataEncryptionAlgorithm ||
					len(expected.Accessors) != len(actual.Accessors) {
					t.Errorf("unexpected response: %v", response)
				}
			}

			for _, invalidToken := range invalidTokens {
				request := tt.request
				request.AccessToken = &api.AccessToken{
					Jwt: &api.JWT{Token: invalidToken},
				}
				response, err := server.DeleteDatum(context.Background(), &request)
				if err == nil {
					t.Error("nil error for invalid token")
				} else if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
					t.Errorf("invalid error for invalid token: %s", err)
				}
				if response != nil {
					t.Errorf("response should have been nil for invalid token: %v", response)
				}
			}
		})
	}

	invalidRequests := []*api.DeleteDatumRequest{
		nil,
		{
			AccessToken: nil,
		},
	}
	t.Run("invalid requests", func(t *testing.T) {
		for _, request := range invalidRequests {
			response, err := server.DeleteDatum(context.Background(), request)
			if response != nil {
				t.Errorf("expected response to be nil: %v", response)
			}
			if err == nil {
				t.Errorf("expected err not to be nil")
			}
			if s, ok := status.FromError(err); !ok || s.Code() != codes.InvalidArgument {
				t.Errorf("unexpected err: %v", err)
			}
		}
	})
}

func TestShareDatum(t *testing.T) {
	// Create mock database
	ctrl := gomock.NewController(t)
	db := NewMockdatabase(ctrl)

	// Initialize server
	server, signer := createServer(t, db)

	validToken, invalidTokens := generate.JWT(&signer, "liam923", jwt.AccessTokenType)

	type fetchDatum struct {
		id    string
		datum datum
		err   error
	}
	tests := []struct {
		testName        string
		request         api.ShareDatumRequest
		fetchDatum      *fetchDatum
		willReachUpdate bool
		updateErr       error
		err             error
	}{
		{
			testName: "valid",
			request: api.ShareDatumRequest{
				Id:          "1",
				TargetId:    "friend",
				DataKey:     &api.EBytes{Data: []byte("234567890")},
				Permissions: []api.Permission{api.Permission_READ, api.Permission_SHARE},
			},
			err: nil,
			fetchDatum: &fetchDatum{
				id: "1",
				datum: datum{
					Owner:                   "liam923",
					Data:                    []byte("23dyoirbeu9"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"liam923": {
							UserId:      "liam923",
							DataKey:     []byte("a0iojcdopi"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
			willReachUpdate: true,
		},
		{
			testName: "valid additional",
			request: api.ShareDatumRequest{
				Id:          "1",
				TargetId:    "friend",
				DataKey:     &api.EBytes{Data: []byte("234567890")},
				Permissions: []api.Permission{api.Permission_READ, api.Permission_SHARE},
			},
			err: nil,
			fetchDatum: &fetchDatum{
				id: "1",
				datum: datum{
					Owner:                   "liam923",
					Data:                    []byte("23dyoirbeu9"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"friend": {
							UserId:      "friend",
							DataKey:     []byte("adhuo289uo"),
							Permissions: []api.Permission{api.Permission_READ},
						},
						"liam923": {
							UserId:      "liam923",
							DataKey:     []byte("a0iojcdopi"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
			willReachUpdate: true,
		},
		{
			testName: "valid non owner",
			request: api.ShareDatumRequest{
				Id:          "1",
				TargetId:    "third party",
				DataKey:     &api.EBytes{Data: []byte("234567890")},
				Permissions: []api.Permission{api.Permission_READ, api.Permission_DELETE},
			},
			err: nil,
			fetchDatum: &fetchDatum{
				id: "1",
				datum: datum{
					Owner:                   "friend",
					Data:                    []byte("23dyoirbeu9"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"liam923": {
							UserId:      "liam923",
							DataKey:     []byte("adhuo289uo"),
							Permissions: []api.Permission{api.Permission_READ, api.Permission_DELETE, api.Permission_SHARE},
						},
						"friend": {
							UserId:      "friend",
							DataKey:     []byte("a0iojcdopi"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
			willReachUpdate: true,
		},
		{
			testName: "invalid id",
			request: api.ShareDatumRequest{
				Id:          "1",
				TargetId:    "friend",
				DataKey:     &api.EBytes{Data: []byte("234567890")},
				Permissions: []api.Permission{api.Permission_READ, api.Permission_SHARE},
			},
			err: status.Error(codes.Internal, "error"),
			fetchDatum: &fetchDatum{
				id:    "1",
				datum: datum{},
				err:   status.Error(codes.Internal, "error"),
			},
		},
		{
			testName: "unauthorized share",
			request: api.ShareDatumRequest{
				Id:          "2",
				TargetId:    "friend",
				DataKey:     &api.EBytes{Data: []byte("234567890")},
				Permissions: []api.Permission{api.Permission_READ, api.Permission_SHARE},
			},
			err: status.Errorf(codes.PermissionDenied, "share access denied for datum 2"),
			fetchDatum: &fetchDatum{
				id: "2",
				datum: datum{
					Owner:                   "other person",
					Data:                    []byte("saofinfwuoirc"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"liam923": {
							UserId:      "liam923",
							DataKey:     []byte("23idjoqnwuiacjs"),
							Permissions: []api.Permission{api.Permission_READ},
						},
						"other person": {
							UserId:      "other person",
							DataKey:     []byte("23ewionpwje"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
		},
		{
			testName: "unauthorized share additional",
			request: api.ShareDatumRequest{
				Id:          "2",
				TargetId:    "friend",
				DataKey:     &api.EBytes{Data: []byte("234567890")},
				Permissions: []api.Permission{api.Permission_READ, api.Permission_SHARE, api.Permission_DELETE},
			},
			err: status.Errorf(codes.PermissionDenied, "attempted to share more access than the user has for datum 2"),
			fetchDatum: &fetchDatum{
				id: "2",
				datum: datum{
					Owner:                   "other person",
					Data:                    []byte("saofinfwuoirc"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"liam923": {
							UserId:      "liam923",
							DataKey:     []byte("23idjoqnwuiacjs"),
							Permissions: []api.Permission{api.Permission_READ, api.Permission_SHARE},
						},
						"other person": {
							UserId:      "other person",
							DataKey:     []byte("23ewionpwje"),
							Permissions: []api.Permission{api.Permission_ADMIN},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
		},
		{
			testName: "share unknown",
			request: api.ShareDatumRequest{
				Id:          "2",
				TargetId:    "friend",
				DataKey:     &api.EBytes{Data: []byte("234567890")},
				Permissions: []api.Permission{api.Permission_UNKNOWN},
			},
			err:        status.Errorf(codes.InvalidArgument, "unable to share \"UNKNOWN\" permission"),
			fetchDatum: nil,
		},
		{
			testName: "share error",
			request: api.ShareDatumRequest{
				Id:          "2",
				TargetId:    "friend",
				DataKey:     &api.EBytes{Data: []byte("234567890")},
				Permissions: []api.Permission{api.Permission_READ, api.Permission_SHARE},
			},
			err: status.Errorf(codes.Internal, "error connecting to database"),
			fetchDatum: &fetchDatum{
				id: "2",
				datum: datum{
					Owner:                   "liam923",
					Data:                    []byte("saofinfwuoirc"),
					DataEncryptionAlgorithm: 0,
					Accessors: map[string]accessor{
						"third party": {
							UserId:      "third party",
							DataKey:     []byte("23idjoqnwuiacjs"),
							Permissions: []api.Permission{api.Permission_READ},
						},
						"liam923": {
							UserId:      "liam923",
							DataKey:     []byte("a0iojcdopi"),
							Permissions: []api.Permission{api.Permission_ADMIN, api.Permission_SHARE},
						},
					},
					Metadata: metadata{},
				},
				err: nil,
			},
			willReachUpdate: true,
			updateErr:       status.Errorf(codes.Internal, "error connecting to database"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf(tt.testName, i), func(t *testing.T) {
			if tt.fetchDatum != nil {
				d := tt.fetchDatum.datum
				db.EXPECT().
					fetchDatum(context.Background(), tt.fetchDatum.id).
					Return(&d, tt.fetchDatum.err)
			}

			if tt.willReachUpdate {
				expectedDatum := tt.fetchDatum.datum
				if expectedDatum.Accessors == nil {
					expectedDatum.Accessors = map[string]accessor{}
				}
				expectedDatum.Accessors[tt.request.TargetId] = accessor{}
				db.EXPECT().
					updateDatum(context.Background(), partialDatumMatcher{datum: expectedDatum}, tt.fetchDatum.id).
					Return(tt.updateErr)
			}

			request := tt.request
			request.AccessToken = &api.AccessToken{
				Jwt: &api.JWT{Token: validToken},
			}

			response, err := server.ShareDatum(context.Background(), &request)
			if tt.err != nil && (err == nil || tt.err.Error() != err.Error()) {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.err == nil {
				if response == nil {
					t.Error("unexpected nil response")
				} else {
					expectedId := tt.fetchDatum.id
					expected := tt.fetchDatum.datum
					if expected.Accessors == nil {
						expected.Accessors = map[string]accessor{}
					}
					expected.Accessors[tt.request.TargetId] = accessor{}
					actual := response.Datum
					if expectedId != actual.Id ||
						expected.Owner != actual.Owner ||
						bytes.Compare(expected.Data, actual.Data.Data) != 0 ||
						expected.DataEncryptionAlgorithm != actual.DataEncryptionAlgorithm ||
						len(expected.Accessors) != len(actual.Accessors) {
						t.Errorf("unexpected response: %v", response)
					}
				}
			}

			for _, invalidToken := range invalidTokens {
				request := tt.request
				request.AccessToken = &api.AccessToken{
					Jwt: &api.JWT{Token: invalidToken},
				}
				response, err := server.ShareDatum(context.Background(), &request)
				if err == nil {
					t.Error("nil error for invalid token")
				} else if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
					t.Errorf("invalid error for invalid token: %s", err)
				}
				if response != nil {
					t.Errorf("response should have been nil for invalid token: %v", response)
				}
			}
		})
	}

	invalidRequests := []*api.ShareDatumRequest{
		nil,
		{
			AccessToken: nil,
		},
		{
			AccessToken: &api.AccessToken{Jwt: &api.JWT{Token: ""}},
			DataKey:     nil,
		},
	}
	t.Run("invalid requests", func(t *testing.T) {
		for _, request := range invalidRequests {
			response, err := server.ShareDatum(context.Background(), request)
			if response != nil {
				t.Errorf("expected response to be nil: %v", response)
			}
			if err == nil {
				t.Errorf("expected err not to be nil")
			}
			if s, ok := status.FromError(err); !ok || s.Code() != codes.InvalidArgument {
				t.Errorf("unexpected err: %v", err)
			}
		}
	})
}
