package flusher

import (
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

	if vs == nil || len(vs) <= 0 {
		return vs, internal.ErrInvalidArg
	}

	s := f.chunkSize
	for len(vs) > 0 {
		if len(vs) < s {
			s = len(vs)
		}
		if err := f.repo.AddVideos(vs[0:s:s]); err != nil {
			return vs, err
		}
		vs = vs[s:] // TODO: test it carefully
	}

	return nil, nil
}
