package account

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/liam923/Kript/server/internal/secure"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc/grpclog"
	"time"
)

// An implementation of AccountService.
type server struct {
	database              database
	Logger                *grpclog.LoggerV2
	signer                secure.JwtSigner
	validator             secure.JwtValidator
	refreshTokenLife      time.Duration
	accessTokenLife       time.Duration
	verificationTokenLife time.Duration
}

// Create a new server. projectId is the id of the project that the Firestore database is housed in, and generate
// is used to sign tokens.
func Server(logger *grpclog.LoggerV2, privateKey []byte, publicKey []byte) (*server, error) {
	googleCreds, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		return nil, err
	}
	client, err := firestore.NewClient(context.Background(), googleCreds.ProjectID)
	if err != nil {
		return nil, err
	}
	signer, err := secure.NewJwtSigner(privateKey, secure.IssuerId)
	if err != nil {
		return nil, err
	}
	validator, err := secure.NewJwtValidator(publicKey, secure.IssuerId)
	if err != nil {
		return nil, err
	}
	return &server{
		database:              &fs{client.Collection("account")},
		Logger:                logger,
		signer:                signer,
		validator:             validator,
		refreshTokenLife:      secure.RefreshTokenLife,
		accessTokenLife:       secure.AccessTokenLife,
		verificationTokenLife: secure.VerificationTokenLife,
	}, nil
}
