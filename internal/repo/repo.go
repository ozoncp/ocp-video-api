package repo

import "ocp-video-api/internal/models"

type Repo interface {
	AddVideos(v []models.Video) error
	AddVideo(v *models.Video) error
	RemoveVideo(ID uint64) error
	GetVideo(ID uint64) (*models.Video, error)
}
