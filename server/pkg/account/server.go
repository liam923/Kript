package account

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/liam923/Kript/server/internal/jwt"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/grpc/grpclog"
)

// The issuer that will be used in JWTs
const issuerId = "kript.api.account"

// An implementation of AccountService.
type Server struct {
	db        *firestore.CollectionRef
	logger    *grpclog.LoggerV2
	signer    *jwt.Signer
	validator *jwt.Validator
}

// Create a new Server. projectId is the id of the project that the Firestore database is housed in, and key
// is used to sign tokens.
func NewServer(projectId string, logger *grpclog.LoggerV2, privateKey []byte, publicKey []byte) (*Server, error) {
	client, err := firestore.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, err
	}
	signer, err := jwt.NewSigner(privateKey, issuerId)
	if err != nil {
		return nil, err
	}
	validator, err := jwt.NewValidator(publicKey, issuerId)
	if err != nil {
		return nil, err
	}
	return &Server{
		db:     client.Collection("account"),
		logger: logger,
		signer: signer,
		validator: validator,
	}, nil
}

func (s *Server) UpdatePassword(context.Context, *api.UpdatePasswordRequest) (*api.UpdatePasswordResponse, error) {
	panic("implement me")
}

func (s *Server) CreateAccount(context.Context, *api.CreateAccountRequest) (*api.CreateAccountResponse, error) {
	panic("implement me")
}

func (s *Server) RefreshAuth(context.Context, *api.RefreshAuthRequest) (*api.RefreshAuthResponse, error) {
	panic("implement me")
}

func (s *Server) GetUser(context.Context, *api.GetUserRequest) (*api.GetUserResponse, error) {
	panic("implement me")
}
