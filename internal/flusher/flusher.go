package flusher

import (
	"ocp-video-api/internal"
	"ocp-video-api/internal/models"
	"ocp-video-api/internal/repo"
)

type Flusher interface {
	Flush(tasks []models.Video) ([]models.Video, error)
}

type flusher struct {
	chunkSize int
	repo repo.Repo
}

func NewFlusher(chunkSize int, repo repo.Repo) Flusher {
	return &flusher{
		chunkSize: chunkSize, repo: repo,
	}
}

func (f *flusher) Flush(vs []models.Video) ([]models.Video, error)  {
	if f.chunkSize <= 0 {
		return nil, internal.ErrInvalidSize
	}

	for len(vs) > 0 {
		// TODO: test what if 0 <= len(vs) <= f.chunkSize
		if err := f.repo.AddVideos(vs[0:f.chunkSize:f.chunkSize]); err != nil {
			return vs, err
		}
		vs = vs[f.chunkSize:]		// TODO: test it carefully
	}

	return nil, nil
}