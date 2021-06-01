package flusher

import (
	"ocp-video-api/internal/models"
)

type Flusher interface {
	Flush(tasks []models.Video) ([]models.Video, error)
}
