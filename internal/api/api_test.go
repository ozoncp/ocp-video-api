package api_test

import (
	"context"
	"database/sql"
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
			createRequest  *desc.CreateVideoV1Request
			createResponse *desc.CreateVideoV1Response
			grpcApi        desc.OcpVideoApiServer
		)

		BeforeEach(func() {
			r := repo.NewRepo(sqlxDB, 2)
			grpcApi = api.NewOcpVideoApi(r)

			createRequest = &desc.CreateVideoV1Request{
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

})
