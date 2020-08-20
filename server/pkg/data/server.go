package data

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/liam923/Kript/server/internal/secure"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc/grpclog"
)

// An implementation of AccountService.
type server struct {
	database  database
	Logger    *grpclog.LoggerV2
	validator secure.JwtValidator
}

// Create a new server. projectId is the id of the project that the Firestore database is housed in, and generate
// is used to sign tokens.
func Server(logger *grpclog.LoggerV2, publicKey []byte) (*server, error) {
	googleCreds, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		return nil, err
	}
	client, err := firestore.NewClient(context.Background(), googleCreds.ProjectID)
	if err != nil {
		return nil, err
	}
	validator, err := secure.NewJwtValidator(publicKey, secure.IssuerId)
	if err != nil {
		return nil, err
	}
	return &server{
		database:  &fs{client.Collection("data")},
		Logger:    logger,
		validator: validator,
	}, nil
}
