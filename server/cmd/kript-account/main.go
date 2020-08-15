package main

import (
	"flag"
	"fmt"
	"github.com/liam923/Kript/server/internal/jwt"
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
	privateJwtKeyPath = flag.String("private-jwt", "", "The path to the private RSA key used to sign JWTs")
	publicJwtKeyPath  = flag.String("public-jwt", "", "The path to the public RSA key used to verify JWTs")
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

	privateKey, err := jwt.ReadRSAKeyFile(*privateJwtKeyPath)
	if err != nil {
		log.Fatalln("Failed to read jwt private key:", err)
		return
	}
	publicKey, err := jwt.ReadRSAKeyFile(*publicJwtKeyPath)
	if err != nil {
		log.Fatalln("Failed to read jwt public key:", err)
		return
	}

	accountServer, err := account.NewServer(&log, privateKey, publicKey)
	if err != nil {
		log.Fatalln("Failed to create data database client:", err)
		return
	}
	api.RegisterAccountServiceServer(s, accountServer)

	// Serve gRPC Server
	log.Info("Serving gRPC on http://", addr)
	log.Fatal(s.Serve(lis))
}