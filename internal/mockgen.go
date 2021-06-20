package internal

//go:generate mockgen -destination=./mock/repo_mock.go -package=mock ocp-video-api/internal/repo Repo
//go:generate mockgen -destination=./mock/flusher_mock.go -package=mock ocp-video-api/internal/flusher Flusher
//go:generate mockgen -destination=./mock/sender_mock.go -package=mock ocp-video-api/internal/producer Sender
//go:generate mockgen -destination=./mock/producer_mock.go -package=mock ocp-video-api/internal/producer Producer
//go:generate mockgen -destination=./mock/metrics_mock.go -package=mock ocp-video-api/internal/metrics Metrics
