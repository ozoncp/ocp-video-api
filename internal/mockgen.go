package internal

//go:generate mockgen -destination=./mock/repo_mock.go -package=mock ocp-video-api/internal/repo Repo
//go:generate mockgen -destination=./mock/flusher_mock.go -package=mock ocp-video-api/internal/flusher Flusher
