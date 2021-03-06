package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"ocp-video-api/internal/api"
	"ocp-video-api/internal/metrics"
	"ocp-video-api/internal/producer"
	"ocp-video-api/internal/repo"
	"ocp-video-api/internal/utils"
	desc "ocp-video-api/pkg/ocp-video-api"
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
		"postgres://goland:goland@localhost:5432/goland?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}

	s := grpc.NewServer()
	r := repo.NewRepo(db, chunkSize)
	p, err := producer.NewProducer(16, producer.NewSaramaSender("localhost:9094", "video"))
	if err != nil {
		log.Fatalf("can't create producer, error: %v", err)
	}
	err = p.Init()
	if err != nil {
		log.Fatalf("can't init producer, error: %v", err)
	}
	m := metrics.New()
	m.Init()
	desc.RegisterOcpVideoApiServer(s, api.NewOcpVideoApi(r, p, m))

	fmt.Printf("Server listening on %s\n", *grpcEndpoint)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	flag.Parse()

	err := utils.InitTracing()
	if err != nil {
		log.Fatalf("can't init tracing: %v", err)
	}

	go runGrpc()

	if err = runHttp(); err != nil {
		log.Fatalf("can't init http server: %v", err)
	}
}
