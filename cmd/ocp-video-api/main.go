package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"ocp-video-api/internal/api"
	"ocp-video-api/internal/producer"
	"ocp-video-api/internal/repo"
	desc "ocp-video-api/pkg/ocp-video-api"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":7002"
	httpPort = ":7000"

	chunkSize = 7
)

var (
	grpcEndpoint = flag.String("grpc-server-endpoint", "0.0.0.0"+grpcPort, "gRPC server endpoint")
	httpEndpoint = flag.String("http-server-endpoint", "0.0.0.0"+httpPort, "HTTP server endpoint")
)

func runHttp() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := desc.RegisterOcpVideoApiHandlerFromEndpoint(ctx, gwmux, *grpcEndpoint, opts)
	if err != nil {
		return err
	}

	//TODO: add swagger?
	// mux.HandleFunc("/swagger/", serveSwagger)
	mux.Handle("/", gwmux)

	fmt.Printf("Server listening on %s\n", *httpEndpoint)
	return http.ListenAndServe(*httpEndpoint, mux)
}

func runGrpc() {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := sqlx.Connect("postgres",
		"postgres://goland:goland@db:5432/goland?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}

	s := grpc.NewServer()
	p, err := producer.NewProducer(16, producer.NewSaramaSender("localhost:9094", "video"))
	if err != nil {
		panic(fmt.Sprintf("can't create producer, error: %v", err))
	}
	desc.RegisterOcpVideoApiServer(s, api.NewOcpVideoApi(repo.NewRepo(db, chunkSize), p))

	fmt.Printf("Server listening on %s\n", *grpcEndpoint)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	flag.Parse()

	go runGrpc()

	if err := runHttp(); err != nil {
		log.Fatal(err)
	}
}
