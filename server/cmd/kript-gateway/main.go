package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"net/http"
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

// Provide a REST API gateway including both account and data endpoints that calls account and data microservices.
func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s", port)
	}

	accountConn, err := newConn(*accountServerAddress, false)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
		return
	}

	dataConn, err := newConn(*dataServerAddress, false)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
		return
	}

	mux := http.NewServeMux()

	jsonpb := &gateway.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)
	err = api.RegisterDataServiceHandler(context.Background(), gwmux, dataConn)
	if err != nil {
		log.Fatalln("Failed to register data gateway:", err)
		return
	}
	err = api.RegisterAccountServiceHandler(context.Background(), gwmux, accountConn)
	if err != nil {
		log.Fatalln("Failed to register account gateway:", err)
		return
	}

	mux.Handle("/", gwmux)

	gatewayAddr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Info("Serving gRPC-Gateway on ", gatewayAddr)
	gwServer := http.Server{
		Addr:    gatewayAddr,
		Handler: mux,
	}
	log.Fatalln(gwServer.ListenAndServe())
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
