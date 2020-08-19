package data

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/liam923/Kript/server/internal/encode"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	firestoreTag     = "firestore"
	firestoreTimeout = time.Minute
)

type accessor struct {
	UserId      string           `firestore:"userId,omitempty"`
	DataKey     []byte           `firestore:"dataKey,omitempty"`
	Permissions []api.Permission `firestore:"permissions,omitempty"`
}

type permissionGrantMetadata struct {
	GranterId  string         `firestore:"granterId,omitempty"`
	Permission api.Permission `firestore:"permission,omitempty"`
	IsGrant    bool           `firestore:"isGrant,omitempty"`
}

type accessMetadata struct {
	GrantMetadata []permissionGrantMetadata `firestore:"grantMetadata,omitempty"`
}

type metadata struct {
	OwnerId        string                    `firestore:"ownerId,omitempty"`
	CreatedTime    time.Time                 `firestore:"createdTime,omitempty"`
	LastEdited     time.Time                 `firestore:"lastEdited,omitempty"`
	AccessMetadata map[string]accessMetadata `firestore:"accessMetadata,omitempty"`
}

type datum struct {
	Owner                   string                   `firestore:"owner,omitempty"`
	Data                    []byte                   `firestore:"data,omitempty"`
	DataEncryptionAlgorithm api.SEncryptionAlgorithm `firestore:"dataEncryptionAlgorithm,omitempty"`
	DataIv                  []byte                   `firestore:"dataIv,omitempty"`
	Accessors               map[string]accessor      `firestore:"accessors,omitempty"`
	Metadata                metadata                 `firestore:"metadata,omitempty"`
}

type idedDatum struct {
	Datum datum
	Id    string
}

func (d datum) toApiDatum(id string) *api.Datum {
	accessors := make(map[string]*api.Datum_Access, len(d.Accessors))
	for key, accessor := range d.Accessors {
		accessors[key] = &api.Datum_Access{
			UserId:      accessor.UserId,
			DataKey:     &api.EBytes{Data: accessor.DataKey},
			Permissions: accessor.Permissions,
		}
	}

	accessMetadata := make(map[string]*api.Datum_Metadata_AccessMetadata, len(d.Metadata.AccessMetadata))
	for key, accessMetadatum := range d.Metadata.AccessMetadata {
		grantMetadata := make([]*api.Datum_Metadata_AccessMetadata_PermissionGrantMetadata, len(accessMetadatum.GrantMetadata))
		for i, grantMetadatum := range accessMetadatum.GrantMetadata {
			grantMetadata[i] = &api.Datum_Metadata_AccessMetadata_PermissionGrantMetadata{
				GranterId:  grantMetadatum.GranterId,
				Permission: grantMetadatum.Permission,
				IsGrant:    grantMetadatum.IsGrant,
			}
		}
		accessMetadata[key] = &api.Datum_Metadata_AccessMetadata{
			GrantMetadata: grantMetadata,
		}
	}

	createdTime, _ := ptypes.TimestampProto(d.Metadata.CreatedTime)
	lastEdited, _ := ptypes.TimestampProto(d.Metadata.LastEdited)
	return &api.Datum{
		Id:                      id,
		Owner:                   d.Owner,
		Data:                    &api.ESecret{Data: d.Data},
		DataEncryptionAlgorithm: d.DataEncryptionAlgorithm,
		DataIv:                  d.DataIv,
		Accessors:               accessors,
		Metadata: &api.Datum_Metadata{
			OwnerId: d.Metadata.OwnerId,
			CreatedTime: &types.Timestamp{
				Seconds: createdTime.Seconds,
				Nanos:   createdTime.Nanos,
			},
			LastEdited: &types.Timestamp{
				Seconds: lastEdited.Seconds,
				Nanos:   lastEdited.Nanos,
			},
			AccessMetadata: accessMetadata,
		},
	}
}

//go:generate mockgen -source firebase.go -destination mocks_test.go -package data

// database is a mockable interface for interacting with Firestore or any other database
type database interface {
	fetchDatum(ctx context.Context, id string) (datum *datum, err error)
	fetchDataForUser(ctx context.Context, userId string) (data *[]idedDatum, err error)
	createDatum(ctx context.Context, datum *datum) (id string, err error)
	updateDatum(ctx context.Context, datum *datum, id string) error
	deleteDatum(ctx context.Context, id string) error
}

// A database implementation for Firestore.
type fs struct {
	db *firestore.CollectionRef
}

func (s *fs) fetchDatum(ctx context.Context, id string) (*datum, error) {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	doc, err := s.db.Doc(id).Get(ctx)
	if err != nil {
		return nil, richError(err)
	}
	datum := &datum{}
	err = doc.DataTo(datum)
	if err != nil {
		return nil, richError(err)
	}
	return datum, nil
}

func (s *fs) fetchDataForUser(ctx context.Context, userId string) (data *[]idedDatum, err error) {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	permissionsPath := fmt.Sprintf("accessors.%s.permissions", userId)
	permissions := []api.Permission{
		api.Permission_READ,
		api.Permission_WRITE,
		api.Permission_DELETE,
		api.Permission_SHARE,
	}
	queries := []firestore.Query{
		s.db.Where("owner", "==", userId),
		s.db.Where(permissionsPath, "array-contains-any", permissions),
	}
	var dataDeref []idedDatum
	for _, query := range queries {
		iter := query.Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				return nil, richError(err)
			}
			datum := &datum{}
			err = doc.DataTo(datum)
			if err != nil {
				return nil, richError(err)
			}
			dataDeref = append(dataDeref, idedDatum{
				Datum: *datum,
				Id:    doc.Ref.ID,
			})
		}
	}
	return &dataDeref, nil
}

func (s *fs) createDatum(ctx context.Context, datum *datum) (id string, err error) {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	data, err := encode.ToMap(*datum, firestoreTag)
	if err != nil {
		return "", richError(err)
	}

	docRef, _, err := s.db.Add(ctx, data)
	if err != nil {
		return "", richError(err)
	}

	return docRef.ID, nil
}

func (s *fs) updateDatum(ctx context.Context, datum *datum, id string) error {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	data, err := encode.ToMap(*datum, firestoreTag)
	if err != nil {
		return richError(err)
	}

	_, err = s.db.Doc(id).Set(ctx, data)
	return richError(err)
}

func (s *fs) deleteDatum(ctx context.Context, id string) error {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	_, err := s.db.Doc(id).Delete(ctx)
	return richError(err)
}

func richError(err error) error {
	if err == nil || status.Code(err) == codes.NotFound {
		return err
	} else {
		return status.Error(codes.Internal, err.Error())
	}
}
