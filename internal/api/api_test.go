package api_test

import (
	"context"
	"database/sql"
	"errors"
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

})
