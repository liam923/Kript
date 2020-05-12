package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/liam923/Kript/server/pkg/account"
	"github.com/liam923/Kript/server/pkg/data"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
)

var (
	gRPCPort    = flag.Int("grpc-port", 10000, "The gRPC server port")
	gatewayPort = flag.Int("rest-port", 11000, "The rest server port")
)

var log grpclog.LoggerV2

func init() {
	log = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
}

func main() {
	flag.Parse()

	addr := fmt.Sprintf("localhost:%d", *gRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	s := grpc.NewServer()
	api.RegisterDataServiceServer(s, &data.Server{})
	api.RegisterAccountServiceServer(s, &account.Server{})

	// Serve gRPC Server
	log.Info("Serving gRPC on http://", addr)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	dialAddr := fmt.Sprintf("passthrough://localhost/%s", addr)
	conn, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	mux := http.NewServeMux()

	jsonpb := &gateway.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		// This is necessary to get error details properly
		// marshalled in unary requests.
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)
	err = api.RegisterDataServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	err = api.RegisterAccountServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	mux.Handle("/", gwmux)

	gatewayAddr := fmt.Sprintf("localhost:%d", *gatewayPort)
	log.Info("Serving gRPC-Gateway on ", gatewayAddr)
	gwServer := http.Server{
		Addr: gatewayAddr,
		Handler: mux,
	}
	log.Fatalln(gwServer.ListenAndServe())
}
