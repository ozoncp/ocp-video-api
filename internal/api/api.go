package api

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ocp-video-api/internal/models"
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

	vs, err := a.repo.GetVideos(ctx, req.Limit, req.Offset)
	if err != nil {
		log.Print("ListVideosV1 error", vs)
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Print("ListVideosV1 found", vs)
	rval := make([]*desc.Video, len(vs))
	innerRval := make([]desc.Video, len(vs))
	for i, v := range vs {
		innerRval[i] = desc.Video{
			Id:      v.VideoId,
			SlideId: v.SlideId,
			Link:    v.Link,
		}
		rval[i] = &innerRval[i]
	}
	return &desc.ListVideosV1Response{Videos: rval}, nil
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

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ID, err := a.repo.AddVideo(ctx, models.Video{
		VideoId: 0,
		SlideId: req.SlideId,
		Link:    req.Link,
	})
	if err != nil {
		log.Print("CreateVideoV1 video is not created due to error", req, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.CreateVideoV1Response{VideoId: ID}, nil
}

func (a *api) MultiCreateVideoV1(
	ctx context.Context,
	req *desc.MultiCreateVideoV1Request,
) (*desc.MultiCreateVideoV1Response, error) {
	log.Print("MultiCreateVideoV1", req)

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	vs := make([]models.Video, len(req.Videos))
	for i, v := range req.Videos {
		vs[i].SlideId = v.SlideId
		vs[i].Link = v.Link
	}

	cnt, err := a.repo.AddVideos(ctx, vs)
	if err != nil {
		log.Print("MultiCreateVideoV1 video is not created due to error", req, err, "created", cnt)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.MultiCreateVideoV1Response{Count: cnt}, nil
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
		log.Print("RemoveVideoV1 video is not removed due to error", req, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	log.Print("RemoveVideoV1 video removed")
	return &desc.RemoveVideoV1Response{
		Found: true,
	}, nil
}

func (a *api) UpdateVideoV1(
	ctx context.Context,
	req *desc.UpdateVideoV1Request,
) (*desc.UpdateVideoV1Response, error) {
	log.Print("UpdateVideoV1", req)

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	v := models.Video{
		VideoId: req.Video.Id,
		SlideId: req.Video.SlideId,
		Link:    req.Video.Link,
	}
	err := a.repo.UpdateVideo(ctx, v)
	if err != nil {
		log.Print("UpdateVideoV1 video is not updated due to error", req, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.UpdateVideoV1Response{
		Found: true,
	}, nil
}

func NewOcpVideoApi(r repo.Repo) desc.OcpVideoApiServer {
	return &api{repo: r}
}
