package api

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	desc "ocp-video-api/pkg/ocp-video-api"
)

type api struct {
	desc.UnimplementedOcpVideoApiServer
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func (a *api) ListVideosV1(
	ctx context.Context,
	req *desc.ListVideosV1Request,
) (*desc.ListVideosV1Response, error) {
	log.Print("ListVideosV1", req)
	return nil, nil
}

func (a *api) DescribeVideoV1(
	ctx context.Context,
	req *desc.DescribeVideoV1Request,
) (*desc.DescribeVideoV1Response, error) {
	log.Print("DescribeTaskV1", req)
	return nil, nil
}

func (a *api) CreateVideoV1(
	ctx context.Context,
	req *desc.CreateVideoV1Request,
) (*desc.CreateVideoV1Response, error) {
	log.Print("CreateVideoV1", req)
	return nil, nil
}

func (a *api) RemoveVideoV1(
	ctx context.Context,
	req *desc.RemoveVideoV1Request,
) (*desc.RemoveVideoV1Response, error) {
	log.Print("RemoveVideoV1", req)
	return nil, nil
}

func NewOcpVideoApi() desc.OcpVideoApiServer {
	return &api{}
}
