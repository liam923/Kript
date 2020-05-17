package account

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/api/iterator"
	"reflect"
)

type password struct {
	Hash          string            `firestore:"hash,omitempty"`
	Salt          string            `firestore:"salt,omitempty"`
	HashAlgorithm api.HashAlgorithm `firestore:"hashAlgorithm,omitempty"`
}

type keys struct {
	PublicKey                     string                   `firestore:"publicKey,omitempty"`
	PrivateKey                    string                   `firestore:"privateKey,omitempty"`
	PrivateKeyEncryptionAlgorithm api.SEncryptionAlgorithm `firestore:"privateKeyEncryptionAlgorithm,omitempty"`
	DataEncryptionAlgorithm       api.AEncryptionAlgorithm `firestore:"dataEncryptionAlgorithm,omitempty"`
}

type twoFactorOption struct {
	Id          string                      `firestore:"id,omitempty"`
	Type        api.TwoFactor_TwoFactorType `firestore:"type,omitempty"`
	Destination string                      `firestore:"destination,omitempty"`
}

type user struct {
	Username  string            `firestore:"username,omitempty"`
	Password  password          `firestore:"password,omitempty"`
	Keys      keys              `firestore:"keys,omitempty"`
	TwoFactor []twoFactorOption `firestore:"twoFactorOption,omitempty"`
}

func (u user) toApiUser(id string, includePrivate bool) *api.User {
	private := &api.PrivateUser{}
	if includePrivate {
		private = &api.PrivateUser{
			PrivateKey:                    u.Keys.PrivateKey,
			PrivateKeyEncryptionAlgorithm: u.Keys.PrivateKeyEncryptionAlgorithm,
		}
	}

	public := &api.PublicUser{
		Id:                      id,
		Username:                u.Username,
		PublicKey:               u.Keys.PublicKey,
		Salt:                    u.Password.Salt,
		PasswordHashAlgorithm:   u.Password.HashAlgorithm,
		DataEncryptionAlgorithm: u.Keys.DataEncryptionAlgorithm,
	}

	return &api.User{
		Public:  public,
		Private: private,
	}
}

//go:generate mockgen -source firebase.go -destination mocks_test.go -package account

// database is a mockable interface for interacting with Firestore or any other database
type database interface {
	fetchUserById(ctx context.Context, userId string) (user *user, err error)
	fetchUserByUsername(ctx context.Context, username string) (user *user, userId string, err error)
	isUsernameAvailable(ctx context.Context, username string) (bool, error)
	updateUser(ctx context.Context, userId string, user *user) error
	createUser(ctx context.Context, user *user) (userId string, err error)
}

// A database implementation for Firestore.
type fs struct {
	db *firestore.CollectionRef
}

func (s *fs) fetchUserById(ctx context.Context, userId string) (user *user, err error) {
	doc, err := s.db.Doc(userId).Get(ctx)
	if err == nil {
		err = doc.DataTo(user)
	}
	return
}

func (s *fs) fetchUserByUsername(ctx context.Context, username string) (user *user, userId string, err error) {
	if field, ok := reflect.TypeOf(user).Elem().FieldByName("Username"); ok {
		usernameField := string(field.Tag)
		iter := s.db.Where(usernameField, "==", username).Limit(1).Documents(ctx)
		doc, err := iter.Next()
		if err == nil {
			err = doc.DataTo(user)
		}
		return nil, "", err
	} else {
		err = fmt.Errorf("could not find field")
		return
	}
}

func (s *fs) isUsernameAvailable(ctx context.Context, username string) (bool, error) {
	if field, ok := reflect.TypeOf(&user{}).Elem().FieldByName("Username"); ok {
		usernameField := string(field.Tag)
		iter := s.db.Where(usernameField, "==", username).Limit(1).Documents(ctx)
		_, err := iter.Next()
		return err == iterator.Done, nil
	} else {
		return false, fmt.Errorf("could not find field")
	}
}

func (s *fs) updateUser(ctx context.Context, userId string, user *user) error {
	_, err := s.db.Doc(userId).Set(ctx, user, firestore.MergeAll)
	return err
}

func (s *fs) createUser(ctx context.Context, user *user) (userId string, err error) {
	docRef, _, err := s.db.Add(ctx, user)
	if err == nil {
		userId = docRef.ID
	}
	return
}
