package repo

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"ocp-video-api/internal"
	"ocp-video-api/internal/models"
	"ocp-video-api/internal/utils"
)

const tableName = "videos"

type Repo interface {
	AddVideos(ctx context.Context, v []models.Video) (int64, error)
	AddVideo(ctx context.Context, v *models.Video) (uint64, error)
	RemoveVideo(ctx context.Context, ID uint64) error
	GetVideo(ctx context.Context, ID uint64) (*models.Video, error)
}

func NewRepo(db *sqlx.DB, chunkSize int) Repo {
	return &repo{db: db, chunkSize: chunkSize}
}

type repo struct {
	chunkSize int
	db        *sqlx.DB
}

func (r *repo) AddVideos(ctx context.Context, vs []models.Video) (int64, error) {
	batches := utils.SliceChunkedModelsVideo(vs, r.chunkSize)

	var pushedCnt int64
	for _, batch := range batches {
		query := squirrel.Insert(tableName).
			Columns("slide_id", "link").
			RunWith(r.db).
			PlaceholderFormat(squirrel.Dollar)

		for _, v := range batch {
			query = query.Values(v.SlideId, v.Link)
		}

		rc, err := query.ExecContext(ctx)
		if err != nil {
			return pushedCnt, err
		}

		added, err := rc.RowsAffected()
		if err != nil {
			return pushedCnt, err
		}
		pushedCnt += added
	}
	return pushedCnt, nil
}

func (r *repo) AddVideo(ctx context.Context, v *models.Video) (uint64, error) {
	query := squirrel.Insert(tableName).
		Columns("slide_id", "link").
		Values(v.SlideId, v.Link).
		Suffix("RETURNING \"id\"").
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	err := query.QueryRowContext(ctx).Scan(&v.VideoId)
	if err != nil {
		return 0, err
	}

	return v.VideoId, nil
}

func (r *repo) RemoveVideo(ctx context.Context, ID uint64) error {
	query := squirrel.Delete(tableName).
		Where(squirrel.Eq{"id": ID}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	result, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	cnt, err := result.RowsAffected()
	if cnt > 1 {
		return nil
	} else {
		return internal.ErrIDNotFound
	}
}

func (r *repo) GetVideo(ctx context.Context, ID uint64) (*models.Video, error) {
	query := squirrel.Select("id", "slide_id", "link").
		From(tableName).
		Where(squirrel.Eq{"id": ID}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	rval := models.Video{}
	if err := query.QueryRowContext(ctx).Scan(&rval.VideoId, rval.SlideId, rval.Link); err != nil {
		return nil, err
	}

	return &rval, nil
}
