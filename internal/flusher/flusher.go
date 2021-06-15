package flusher

import (
	"context"
	"ocp-video-api/internal"
	"ocp-video-api/internal/models"
	"ocp-video-api/internal/repo"
)

type Flusher interface {
	Flush(vs []models.Video) ([]models.Video, error)
}

type flusher struct {
	chunkSize int
	repo      repo.Repo
}

func New(chunkSize int, repo repo.Repo) Flusher {
	return &flusher{
		chunkSize: chunkSize, repo: repo,
	}
}

func (f *flusher) Flush(vs []models.Video) ([]models.Video, error) {
	if f.chunkSize <= 0 {
		return vs, internal.ErrInvalidSize
	}

	s := f.chunkSize
	for len(vs) > 0 {
		if len(vs) < s {
			s = len(vs)
		}
		var (
			handledCnt uint64
			err error
		)
		if handledCnt, err = f.repo.AddVideos(context.Background(), vs[0:s:s]); err != nil {
			return vs[handledCnt:], err
		}
		if handledCnt != uint64(s) {
			panic("logic error: handledCnt != batch && err = nil")
		}
		vs = vs[handledCnt:]
	}

	return nil, nil
}
