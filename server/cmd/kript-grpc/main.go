package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/liam923/Kript/server/internal/forward"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"net"
	"os"
)

var (
	accountServerAddress = flag.String("account-address", "", "The address of the account microservice")
	dataServerAddress    = flag.String("data-address", "", "The address of the data microservice")
)

var log grpclog.LoggerV2

func init() {
	log = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
}

// Provide a gRPC API gateway including both account and data endpoints that calls account and data microservices.
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

	accountConn, err := newConn(*accountServerAddress, false)
	if err != nil {
		log.Fatalln("Failed to connect to account server: ", err)
	}
	defer accountConn.Close()
	accountClient := api.NewAccountServiceClient(accountConn)
	accountServer := forward.NewAccountForwarder(accountClient)
	api.RegisterAccountServiceServer(s, accountServer)

	dataConn, err := newConn(*dataServerAddress, false)
	if err != nil {
		log.Fatalln("Failed to connect to data server: ", err)
	}
	defer dataConn.Close()
	dataClient := api.NewDataServiceClient(dataConn)
	dataServer := forward.NewDataForwarder(dataClient)
	api.RegisterDataServiceServer(s, dataServer)

	// Serve gRPC Server
	log.Info("Serving gRPC on http://", addr)
	log.Fatal(s.Serve(lis))
	return
}

func newConn(host string, insecure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	return grpc.Dial(host, opts...)
}
