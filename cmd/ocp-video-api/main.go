package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	goruntime "runtime"

	"ocp-video-api/internal/api"
	desc "ocp-video-api/pkg/ocp-video-api"

	"github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":7002"
	httpPort = ":7000"
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

	s := grpc.NewServer()
	desc.RegisterOcpVideoApiServer(s, api.NewOcpVideoApi())

	fmt.Printf("Server listening on %s\n", *grpcEndpoint)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func printIsContainered() {
	myOS, myArch := goruntime.GOOS, goruntime.GOARCH
	inContainer := "inside"
	if _, err := os.Lstat("/.dockerenv"); err != nil && os.IsNotExist(err) {
		inContainer = "outside"
	}
	fmt.Println("2... I'm running on:\n", myOS, myArch)
	fmt.Println("I'm running container status:\n", inContainer)
}

type WidgetSerialized struct {
	Name   string `json:"name,omitempty"`
	Weight int64  `json:"weight,omitempty"`
}

func testComposedDb() {
	ctx := context.Background()
	db, err := sqlx.Connect("postgres",
		"postgres://goland:goland@db:5432/goland?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}

	query := squirrel.Select("name", "weight").
		From("widgets").
		Where(squirrel.Eq{"id": 42}).
		RunWith(db).
		PlaceholderFormat(squirrel.Dollar)

	var widget WidgetSerialized
	if err := query.QueryRowContext(ctx).Scan(&widget.Name, &widget.Weight); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(widget)
}

func main() {
	printIsContainered()

	testComposedDb()

	flag.Parse()

	go runGrpc()

	if err := runHttp(); err != nil {
		log.Fatal(err)
	}
}
