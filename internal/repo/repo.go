package repo

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"ocp-video-api/internal/models"
	"ocp-video-api/internal/utils"
)

const tableName = "videos"

type Repo interface {
	AddVideos(ctx context.Context, v []models.Video) ([]uint64, error)
	AddVideo(ctx context.Context, v models.Video) (uint64, error)
	RemoveVideo(ctx context.Context, ID uint64) error
	GetVideo(ctx context.Context, ID uint64) (*models.Video, error)
	GetVideos(ctx context.Context, limit, offset uint64) ([]models.Video, error)
	UpdateVideo(ctx context.Context, v models.Video) error
}

func NewRepo(db *sqlx.DB, chunkSize int) Repo {
	return &repo{db: db, chunkSize: chunkSize}
}

type repo struct {
	chunkSize int
	db        *sqlx.DB
}

func (r *repo) AddVideos(ctx context.Context, vs []models.Video) ([]uint64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("Adding %v videos", len(vs)))
	defer span.Finish()

	log.Print("Adding videos", vs)
	batches := utils.SliceChunkedModelsVideo(vs, r.chunkSize)

	ids := make([]uint64, 0, len(vs))
	for _, batch := range batches {
		err := r.insertBatch(ctx, batch, &ids)
		if err != nil {
			return ids, err
		}
	}
	log.Print("Videos succesfully pushed, ids:", ids)
	return ids, nil
}

func (r *repo) insertBatch(ctx context.Context, batch []models.Video, dstIDs *[]uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("Create batch with %v videos", len(batch)))
	defer span.Finish()

	query := squirrel.Insert(tableName).
		Columns("slide_id", "link").
		Suffix("RETURNING \"id\"").
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	for _, v := range batch {
		query = query.Values(v.SlideId, v.Link)
	}

	rows, err := query.Query()
	if err != nil {
		log.Print("Error pushing batch", batch, "error", err)
		return err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Print("Error closing rows in batch insert", err)
		}
	}()

	for rows.Next() {
		var id uint64
		errScan := rows.Scan(&id)
		if errScan != nil {
			log.Print(errScan)
			return errScan
		}
		*dstIDs = append(*dstIDs, id)
	}

	return nil
}

func (r *repo) AddVideo(ctx context.Context, v models.Video) (uint64, error) {
	log.Print("Adding single video", v)
	query := squirrel.Insert(tableName).
		Columns("slide_id", "link").
		Values(v.SlideId, v.Link).
		Suffix("RETURNING \"id\"").
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	err := query.QueryRowContext(ctx).Scan(&v.VideoId)
	if err != nil {
		log.Print("Error pushing video", err)
		return 0, err
	}

	return v.VideoId, nil
}

func (r *repo) RemoveVideo(ctx context.Context, ID uint64) error {
	log.Print("Removing single video with ID", ID)
	query := squirrel.Delete(tableName).
		Where(squirrel.Eq{"id": ID}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	_, err := query.ExecContext(ctx)
	if err != nil {
		log.Print("Error deleting video ID", ID, "error", err)
		return err
	}

	return nil
}

func (r *repo) GetVideo(ctx context.Context, ID uint64) (*models.Video, error) {
	log.Print("Get single video with ID", ID)
	query := squirrel.Select("id", "slide_id", "link").
		From(tableName).
		Where(squirrel.Eq{"id": ID}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	rval := models.Video{}
	if err := query.QueryRowContext(ctx).Scan(&rval.VideoId, &rval.SlideId, &rval.Link); err != nil {
		log.Print("Error retrieving video with ID", ID, "error", err)
		return nil, err
	}

	return &rval, nil
}

func (r *repo) GetVideos(ctx context.Context, limit, offset uint64) ([]models.Video, error) {
	query := squirrel.Select("id", "slide_id", "link").
		From(tableName).
		RunWith(r.db).
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	vs := make([]models.Video, 0, limit)
	for rows.Next() {
		var v models.Video
		if err = rows.Scan(&v.VideoId, &v.SlideId, &v.Link); err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs[:len(vs):len(vs)], nil
}

func (r *repo) UpdateVideo(ctx context.Context, v models.Video) error {
	query := squirrel.Update(tableName).
		Set("slide_id", v.SlideId).
		Set("link", v.Link).
		Where(squirrel.Eq{"id": v.VideoId}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	_, err := query.ExecContext(ctx)
	return err
}
