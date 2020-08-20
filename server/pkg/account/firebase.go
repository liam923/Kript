package account

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/liam923/Kript/server/internal/encode"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	firestoreTag                   = "firestore"
	firestoreTimeout               = time.Minute
	verificationCodesSubcollection = "verificationCodes"
)

type password struct {
	Hash          []byte            `firestore:"hash,omitempty"`
	Salt          []byte            `firestore:"salt,omitempty"`
	HashAlgorithm api.HashAlgorithm `firestore:"hashAlgorithm,omitempty"`
}

type keys struct {
	PublicKey                     []byte                   `firestore:"publicKey,omitempty"`
	PrivateKey                    []byte                   `firestore:"privateKey,omitempty"`
	PrivateKeyEncryptionAlgorithm api.SEncryptionAlgorithm `firestore:"privateKeyEncryptionAlgorithm,omitempty"`
	PrivateKeyIv                  []byte                   `firestore:"privateKeyIv,omitempty"`
	PrivateKeyKeySalt             []byte                   `firestore:"privateKeyKeySalt,omitempty"`
	PrivateKeyKeyHashAlgorithm    api.HashAlgorithm        `firestore:"privateKeyKeyHashAlgorithm,omitempty"`
	DataEncryptionAlgorithm       api.AEncryptionAlgorithm `firestore:"dataEncryptionAlgorithm,omitempty"`
}

type twoFactorOption struct {
	Type        api.TwoFactorType `firestore:"type,omitempty"`
	Destination string            `firestore:"destination,omitempty"`
}

type user struct {
	Username  string                     `firestore:"username,omitempty"`
	Password  password                   `firestore:"password,omitempty"`
	Keys      keys                       `firestore:"keys,omitempty"`
	TwoFactor map[string]twoFactorOption `firestore:"twoFactorOption,omitempty"`
}

func (u user) toApiUser(id string, includePrivate bool) *api.User {
	var private *api.PrivateUser
	if includePrivate {
		private = &api.PrivateUser{
			PrivateKey:                    &api.EBytes{Data: u.Keys.PrivateKey},
			PrivateKeyEncryptionAlgorithm: u.Keys.PrivateKeyEncryptionAlgorithm,
			PrivateKeyIv:                  u.Keys.PrivateKeyIv,
			PrivateKeyKeySalt:             u.Keys.PrivateKeyKeySalt,
			PrivateKeyKeyHashAlgorithm:    u.Keys.PrivateKeyKeyHashAlgorithm,
		}
	}

	public := &api.PublicUser{
		Id:                      id,
		Username:                u.Username,
		PublicKey:               u.Keys.PublicKey,
		PasswordSalt:            u.Password.Salt,
		PasswordHashAlgorithm:   u.Password.HashAlgorithm,
		DataEncryptionAlgorithm: u.Keys.DataEncryptionAlgorithm,
	}

	return &api.User{
		Public:  public,
		Private: private,
	}
}

type verificationCode struct {
	Code                  string          `firestore:"Code,omitempty"`
	HasConfirmDestination bool            `firestore:"hasDestination,omitempty"`
	ConfirmDestination    twoFactorOption `firestore:"destination,omitempty"`
}

//go:generate mockgen -source firebase.go -destination mocks_test.go -package account

// database is a mockable interface for interacting with Firestore or any other database
type database interface {
	fetchUserById(ctx context.Context, userId string) (user *user, err error)
	fetchUserByUsername(ctx context.Context, username string) (user *user, userId string, err error)
	isUsernameAvailable(ctx context.Context, username string) (bool, error)
	updateUser(ctx context.Context, userId string, user *user) error
	createUser(ctx context.Context, user *user) (userId string, err error)
	addVerificationTokenCode(ctx context.Context, userId string, tokenId string, code string, confirmDestination *twoFactorOption) error
	verifyVerificationTokenCode(ctx context.Context, userId string, tokenId string, code string) (confirmDestination *twoFactorOption, err error)
}

// A database implementation for Firestore.
type fs struct {
	db *firestore.CollectionRef
}

func (s *fs) fetchUserById(ctx context.Context, userId string) (*user, error) {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	doc, err := s.db.Doc(userId).Get(ctx)
	user := &user{}
	if err == nil {
		err = doc.DataTo(user)
	}
	err = richError(err)
	if err != nil {
		user = nil
	}
	return user, err
}

func (s *fs) fetchUserByUsername(ctx context.Context, username string) (*user, string, error) {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	iter := s.db.Where("username", "==", username).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == nil {
		user := &user{}
		err = doc.DataTo(user)
		if err == nil {
			return user, doc.Ref.ID, nil
		}
	} else if err == iterator.Done {
		return nil, "",
			status.Errorf(codes.NotFound, "could not find user with username %s", username)
	}
	return nil, "", status.Error(codes.Internal, err.Error())
}

func (s *fs) isUsernameAvailable(ctx context.Context, username string) (bool, error) {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	iter := s.db.Where("username", "==", username).Limit(1).Documents(ctx)
	_, err := iter.Next()
	if err == iterator.Done {
		return true, nil
	} else if err == nil {
		return false, nil
	} else {
		return false, richError(err)
	}
}

func (s *fs) updateUser(ctx context.Context, userId string, user *user) error {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	data, err := encode.ToMap(*user, firestoreTag)
	if err != nil {
		return richError(err)
	}

	_, err = s.db.Doc(userId).Set(ctx, data, firestore.MergeAll)
	return richError(err)
}

func (s *fs) createUser(ctx context.Context, user *user) (userId string, err error) {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	data, err := encode.ToMap(*user, firestoreTag)
	if err != nil {
		return "", richError(err)
	}

	docRef, _, err := s.db.Add(ctx, data)
	if err != nil {
		return "", richError(err)
	}

	return docRef.ID, nil
}

func (s *fs) addVerificationTokenCode(ctx context.Context, userId string, tokenId string, code string, confirmDestination *twoFactorOption) error {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	codeStruct := verificationCode{
		Code:                  code,
		HasConfirmDestination: false,
	}
	if confirmDestination != nil {
		codeStruct.ConfirmDestination = *confirmDestination
		codeStruct.HasConfirmDestination = true
	}

	data, err := encode.ToMap(codeStruct, firestoreTag)
	if err != nil {
		return richError(err)
	}

	_, err = s.db.Doc(userId).Collection(verificationCodesSubcollection).Doc(tokenId).Create(ctx, data)
	if err != nil {
		return richError(err)
	}

	return nil
}

func (s *fs) verifyVerificationTokenCode(ctx context.Context, userId string, tokenId string, code string) (confirmDestination *twoFactorOption, err error) {
	ctx, _ = context.WithTimeout(ctx, firestoreTimeout)

	doc, err := s.db.Doc(userId).Collection(verificationCodesSubcollection).Doc(tokenId).Get(ctx)
	verificationCode := &verificationCode{}
	if err == nil {
		err = richError(doc.DataTo(verificationCode))
		if verificationCode.Code != code {
			err = status.Error(codes.Unauthenticated, "invalid verification Code")
		} else if verificationCode.HasConfirmDestination {
			confirmDestination = &verificationCode.ConfirmDestination
		}
	} else {
		err = richError(err)
	}

	if err != nil {
		confirmDestination = nil
	}
	return
}

func richError(err error) error {
	if err == nil || status.Code(err) == codes.NotFound {
		return err
	} else {
		return status.Error(codes.Internal, err.Error())
	}
}
