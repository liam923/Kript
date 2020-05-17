package account

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/liam923/Kript/server/internal/jwt"
	"google.golang.org/grpc/grpclog"
	"time"
)

// An implementation of AccountService.
type Server struct {
	database              database
	Logger                *grpclog.LoggerV2
	signer                *jwt.Signer
	validator             *jwt.Validator
	refreshTokenLife      time.Duration
	accessTokenLife       time.Duration
	verificationTokenLife time.Duration
}

// Create a new Server. projectId is the id of the project that the Firestore database is housed in, and generate
// is used to sign tokens.
func NewServer(projectId string, logger *grpclog.LoggerV2, privateKey []byte, publicKey []byte) (*Server, error) {
	client, err := firestore.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, err
	}
	signer, err := jwt.NewSigner(privateKey, jwt.IssuerId)
	if err != nil {
		return nil, err
	}
	validator, err := jwt.NewValidator(publicKey, jwt.IssuerId)
	if err != nil {
		return nil, err
	}
	return &Server{
		database:              &fs{client.Collection("account")},
		Logger:                logger,
		signer:                signer,
		validator:             validator,
		refreshTokenLife:      jwt.RefreshTokenLife,
		accessTokenLife:       jwt.AccessTokenLife,
		verificationTokenLife: jwt.VerificationTokenLife,
	}, nil
}
