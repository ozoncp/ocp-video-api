.PHONY: build
build: vendor-proto generate .build

.PHONY: .all
.all: dependencies build

.PHONY: all
all: requirements .all

PHONY: generate
generate:
		mkdir -p pkg/ocp-video-api

		protoc -I vendor.protogen \
				--go_out=pkg/ocp-video-api --go_opt=paths=import \
				--go-grpc_out=pkg/ocp-video-api --go-grpc_opt=paths=import \
				--grpc-gateway_out=pkg/ocp-video-api \
				--grpc-gateway_opt=logtostderr=true \
				--grpc-gateway_opt=paths=import \
				--validate_out lang=go:pkg/ocp-video-api \
				api/ocp-video-api/ocp-video-api.proto
		mv pkg/ocp-video-api/ocp-video-api/pkg/ocp-video-api/* pkg/ocp-video-api/
		rm -rf pkg/ocp-video-api/ocp-video-api
		mkdir -p cmd/ocp-video-api

PHONY: .build
.build:
	CGO_ENABLED=0 GOOS=linux go build -o /server cmd/ocp-video-api/main.go

.PHONY: requirements
requirements:
		apt-get update -q && apt-get install -y protobuf-compiler

.PHONY: dependencies
dependencies:
		ls go.mod || go mod init

		go get -u github.com/envoyproxy/protoc-gen-validate
		go get -u github.com/golang/protobuf/proto
		go get -u github.com/golang/protobuf/protoc-gen-go
		go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
		go get -u github.com/rs/zerolog/log
		go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go get -u google.golang.org/protobuf/runtime/protoimpl
		go get -u google.golang.org/protobuf/types/known/anypb

		go mod download

		go install github.com/envoyproxy/protoc-gen-validate
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

PHONY: vendor-proto
vendor-proto:
		mkdir -p vendor.protogen/api/ocp-video-api

		cp -rf api/ocp-video-api/ocp-video-api.proto vendor.protogen/api/ocp-video-api/

		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi

		@if [ ! -d vendor.protogen/github.com/envoyproxy ]; then \
			mkdir -p vendor.protogen/github.com/envoyproxy &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/github.com/envoyproxy/protoc-gen-validate ;\
		fi

.PHONY: clean
clean:
		rm -rf bin pkg vendor.protogen

.PHONY: tidy
tidy:
		go mod tidy