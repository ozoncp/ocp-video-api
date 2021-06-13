package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ocp-video-api/internal/repo"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	desc "ocp-video-api/pkg/ocp-video-api"
)

type api struct {
	desc.UnimplementedOcpVideoApiServer
	repo repo.Repo
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func (a *api) ListVideosV1(
	ctx context.Context,
	req *desc.ListVideosV1Request,
) (*desc.ListVideosV1Response, error) {
	log.Print("ListVideosV1", req)

	if err := req.Validate(); err != nil {
		log.Print("ListVideosV1 invalid request", req)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dst := make([]*desc.Video, req.Limit)
	for i := req.Offset; i < req.Offset+req.Limit; i++ {
		v, err := a.repo.GetVideo(ctx, i)
		if err != nil {
			log.Print("Error retrieving video with ID", i)
			break
		}
		dst = append(dst, &desc.Video{
			Id:      v.VideoId,
			SlideId: v.SlideId,
			Link:    v.Link,
		})
	}

	if len(dst) > 0 {
		log.Print("ListVideosV1 found", dst)
		return &desc.ListVideosV1Response{Videos: dst}, nil
	} else {
		log.Print("ListVideosV1 no videos found", req)
		return nil, status.Error(codes.NotFound, fmt.Sprintf("No videos fount, offset: %v", req.Offset))
	}
}

func (a *api) DescribeVideoV1(
	ctx context.Context,
	req *desc.DescribeVideoV1Request,
) (*desc.DescribeVideoV1Response, error) {
	log.Print("DescribeVideoV1", req)

	if err := req.Validate(); err != nil {
		log.Print("DescribeVideoV1 invalid request", req)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	v, err := a.repo.GetVideo(ctx, req.VideoId)
	if err != nil {
		log.Print("DescribeVideoV1 query error", req, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	log.Print("DescribeVideoV1 found", v)
	return &desc.DescribeVideoV1Response{
		Video: &desc.Video{
			Id:      v.VideoId,
			SlideId: v.SlideId,
			Link:    v.Link},
	}, nil
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

	if err := req.Validate(); err != nil {
		log.Print("RemoveVideoV1 invalid request", req)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := a.repo.RemoveVideo(ctx, req.VideoId)
	if err != nil {
		log.Print("RemoveVideoV1 video not removed due to error", req, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	log.Print("RemoveVideoV1 video removed")
	return &desc.RemoveVideoV1Response{
		Found: true,
	}, nil
}

func NewOcpVideoApi(r repo.Repo) desc.OcpVideoApiServer {
	return &api{repo: r}
}
