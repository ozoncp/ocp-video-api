package api_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"ocp-video-api/internal/repo"
	desc "ocp-video-api/pkg/ocp-video-api"

	"ocp-video-api/internal/api"
)

var _ = Describe("Api", func() {
	var (
		ctx context.Context

		db     *sql.DB
		sqlxDB *sqlx.DB
		mock   sqlmock.Sqlmock

		err error
	)

	BeforeEach(func() {
		ctx = context.Background()

		db, mock, err = sqlmock.New()
		Expect(err).Should(BeNil())

		sqlxDB = sqlx.NewDb(db, "sqlmock")
	})

	JustBeforeEach(func() {
	})

	AfterEach(func() {
		mock.ExpectClose()
		err = db.Close()
		Expect(err).Should(BeNil())
	})

	Context("create video", func() {
		var (
			ID             uint64 = 1
			createResponse *desc.CreateVideoV1Response
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)

			createRequest := &desc.CreateVideoV1Request{
				SlideId: 42,
				Link:    "link/to/42",
			}

			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(ID)
			mock.ExpectQuery("INSERT INTO videos").
				WithArgs(createRequest.SlideId, createRequest.Link).
				WillReturnRows(rows)

			createResponse, err = grpcApi.CreateVideoV1(ctx, createRequest)
		})

		It("", func() {
			Expect(err).Should(BeNil())
			Expect(createResponse.VideoId).Should(Equal(ID))
		})
	})

	Context("create video invalid argument", func() {
		var (
			createResponse *desc.CreateVideoV1Response
			err            error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)
			createRequest := &desc.CreateVideoV1Request{SlideId: 0, Link: "invalid/slide/id"}
			createResponse, err = grpcApi.CreateVideoV1(ctx, createRequest)
		})

		It("", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(createResponse).Should(BeNil())
		})
	})

	Context("create video sql error", func() {
		var (
			createResponse *desc.CreateVideoV1Response
			err            error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)

			createRequest := &desc.CreateVideoV1Request{
				SlideId: 7,
				Link:    "42",
			}

			mock.ExpectQuery("INSERT INTO videos").
				WithArgs(createRequest.SlideId, createRequest.Link).
				WillReturnError(errors.New("mock db with different table schema error"))

			createResponse, err = grpcApi.CreateVideoV1(ctx, createRequest)
		})

		It("", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(createResponse).Should(BeNil())
		})
	})

	Context("get video simple", func() {
		var (
			videoId uint64 = 1
			slideId uint64 = 42
			link           = "/link/to/42"

			getResponse *desc.DescribeVideoV1Response
			err         error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)

			getRequest := &desc.DescribeVideoV1Request{VideoId: 1}

			rows := sqlmock.NewRows([]string{"id", "slide_id", "link"}).
				AddRow(videoId, slideId, link)
			mock.ExpectQuery("SELECT (.+) FROM videos WHERE").
				WithArgs(getRequest.VideoId).
				WillReturnRows(rows)

			getResponse, err = grpcApi.DescribeVideoV1(ctx, getRequest)
		})

		It("", func() {
			Expect(err).Should(BeNil())
			Expect(getResponse.Video.Id).Should(Equal(videoId))
			Expect(getResponse.Video.SlideId).Should(Equal(slideId))
			Expect(getResponse.Video.Link).Should(Equal(link))
		})
	})

	Context("get video invalid argument", func() {
		var (
			getResponse *desc.DescribeVideoV1Response
			err         error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)

			getRequest := &desc.DescribeVideoV1Request{VideoId: 0}

			getResponse, err = grpcApi.DescribeVideoV1(ctx, getRequest)
		})

		It("", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(getResponse).Should(BeNil())
		})
	})

	Context("get video sql error", func() {
		var (
			getResponse *desc.DescribeVideoV1Response
			err         error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)

			getRequest := &desc.DescribeVideoV1Request{VideoId: 1}

			mock.ExpectQuery("SELECT (.+) FROM videos WHERE").
				WithArgs(getRequest.VideoId).
				WillReturnError(errors.New("mock db with different table schema error"))

			getResponse, err = grpcApi.DescribeVideoV1(ctx, getRequest)
		})

		It("", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(getResponse).Should(BeNil())
		})
	})

	Context("remove video", func() {
		var (
			videoId uint64 = 1
			removeResponse *desc.RemoveVideoV1Response
			err error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)
			removeRequest := &desc.RemoveVideoV1Request{ VideoId: videoId }

			mock.ExpectExec("DELETE FROM videos").
				WithArgs(removeRequest.VideoId).WillReturnResult(sqlmock.NewResult(0, 1))

			removeResponse, err = grpcApi.RemoveVideoV1(ctx, removeRequest)
		})

		It("", func() {
			Expect(err).Should(BeNil())
			Expect(removeResponse.Found).Should(Equal(true))
		})
	})

	// RowsAffected() комментарии: Not every database or database driver may support this
	// тестируем пессимистичный случай
	Context("remove video, rowsAffected returns 0", func() {
		var (
			videoId uint64 = 1
			removeResponse *desc.RemoveVideoV1Response
			err error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)
			removeRequest := &desc.RemoveVideoV1Request{ VideoId: videoId }

			mock.ExpectExec("DELETE FROM videos").
				WithArgs(removeRequest.VideoId).WillReturnResult(sqlmock.NewResult(0, 0))

			removeResponse, err = grpcApi.RemoveVideoV1(ctx, removeRequest)
		})

		It("", func() {
			Expect(err).Should(BeNil())
			Expect(removeResponse.Found).Should(Equal(true))
		})
	})

	Context("remove video invalid argument", func() {
		var (
			removeResponse *desc.RemoveVideoV1Response
			err error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)
			removeRequest := &desc.RemoveVideoV1Request{ VideoId: 0 }

			removeResponse, err = grpcApi.RemoveVideoV1(ctx, removeRequest)
		})

		It("", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(removeResponse).Should(BeNil())
		})
	})

	Context("remove video sql error", func() {
		var (
			removeResponse *desc.RemoveVideoV1Response
			err error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)
			removeRequest := &desc.RemoveVideoV1Request{ VideoId: 1 }

			mock.ExpectExec("DELETE FROM videos").
				WithArgs(removeRequest.VideoId).
				WillReturnError(errors.New("mock db with different table schema error"))

			removeResponse, err = grpcApi.RemoveVideoV1(ctx, removeRequest)
		})

		It("", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(removeResponse).Should(BeNil())
		})
	})

	Context("list videos", func() {
		var (
			limit  uint64 = 10
			offset uint64 = 0

			listResponse *desc.ListVideosV1Response
			err error
			videos = [2]desc.Video{
				desc.Video{
					Id: 1,
					SlideId: 7,
					Link:    "/link/to/7",
				},
				desc.Video{
					Id: 2,
					SlideId: 42,
					Link:    "/link/to/42",
				},
			}
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)
			listRequest := &desc.ListVideosV1Request{Limit: limit, Offset: offset}

			rows := sqlmock.NewRows([]string{"id", "slide_id", "link"}).
				AddRow(videos[0].Id, videos[0].SlideId, videos[0].Link).
				AddRow(videos[1].Id, videos[1].SlideId, videos[1].Link)

			mock.ExpectQuery(
				fmt.Sprintf("SELECT id, slide_id, link FROM videos LIMIT %d OFFSET %d", limit, offset)).
				WillReturnRows(rows)

			listResponse, err = grpcApi.ListVideosV1(ctx, listRequest)
		})

		It("", func() {
			Expect(err).Should(BeNil())
			Expect(len(listResponse.Videos)).Should(Equal(len(videos)))
		})
	})

	Context("list videos sql error", func() {
		var (
			limit  uint64 = 10
			offset uint64 = 0

			listResponse *desc.ListVideosV1Response
			err error
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi := api.NewOcpVideoApi(r)
			listRequest := &desc.ListVideosV1Request{Limit: limit, Offset: offset}

			mock.ExpectQuery(
				fmt.Sprintf("SELECT id, slide_id, link FROM videos LIMIT %d OFFSET %d", limit, offset)).
				WillReturnError(errors.New("mock db with different table schema error"))

			listResponse, err = grpcApi.ListVideosV1(ctx, listRequest)
		})

		It("", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(listResponse).Should(BeNil())
		})
	})
})
