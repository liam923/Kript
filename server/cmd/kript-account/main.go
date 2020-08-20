package main

import (
	"flag"
	"fmt"
	"github.com/liam923/Kript/server/internal/secure"
	"github.com/liam923/Kript/server/pkg/account"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"net"
	"os"
)

var (
	privateJwtKeyPath  = flag.String("private-jwt", "", "The path to the private RSA key used to sign JWTs")
	publicJwtKeyPath   = flag.String("public-jwt", "", "The path to the public RSA key used to verify JWTs")
	sendgridApiKeyPath = flag.String("sendgrid-api-key", "", "The path to the API key for SendGrid")
)

var log grpclog.LoggerV2

func init() {
	log = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
}

// Provide a gRPC API with data endpoints.
func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s", port)
	}

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
		return
	}
	s := grpc.NewServer()
	reflection.Register(s)

	privateKey, err := secure.ReadRSAKeyFile(*privateJwtKeyPath)
	if err != nil {
		log.Fatalln("Failed to read jwt private key:", err)
		return
	}
	publicKey, err := secure.ReadRSAKeyFile(*publicJwtKeyPath)
	if err != nil {
		log.Fatalln("Failed to read jwt public key:", err)
		return
	}
	sendGridApiKeyData, err := ioutil.ReadFile(*sendgridApiKeyPath)
	if err != nil {
		log.Fatalln("Failed to read SendGrid api key:", err)
		return
	}
	sendgridApiKey := string(sendGridApiKeyData)

	accountServer, err := account.Server(&log, privateKey, publicKey, sendgridApiKey)
	if err != nil {
		log.Fatalln("Failed to create data database client:", err)
		return
	}
	api.RegisterAccountServiceServer(s, accountServer)

	// Serve gRPC server
	log.Info("Serving gRPC on http://", addr)
	log.Fatal(s.Serve(lis))
}
