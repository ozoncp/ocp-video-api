package saver_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"ocp-video-api/internal/mock"
	"ocp-video-api/internal/models"
	"ocp-video-api/internal/saver"
)

const testCapacity = 2

var _ = Describe("Saver", func() {

	var (
		err error

		ctrl *gomock.Controller
		ctx  context.Context

		mockFlusher *mock.MockFlusher

		v models.Video
		s saver.Saver
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())

		mockFlusher = mock.NewMockFlusher(ctrl)
	})

	JustBeforeEach(func() {
		s.Init(ctx)
	})

	AfterEach(func() {
		s.Close()
		ctrl.Finish()
	})

	Context("ctx canceled", func() {
		var (
			cancelFunc context.CancelFunc
		)

		BeforeEach(func() {
			s = saver.New(testCapacity, time.Millisecond*10, mockFlusher)
			v = models.Video{
				VideoId: 1,
				SlideId: 42,
				Link:    "/video/42",
			}
			ctx, cancelFunc = context.WithCancel(ctx)
			mockFlusher.EXPECT().Flush([]models.Video{models.Video{
				VideoId: 1,
				SlideId: 42,
				Link:    "/video/42",
			}}).Return(nil, nil)
		})

		JustBeforeEach(func() {
			err = s.Save(ctx, v)
			cancelFunc()
		})

		It("close do not hangs", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("saver works", func() {
		var (
			cancelFunc context.CancelFunc
		)

		BeforeEach(func() {
			s = saver.New(testCapacity, time.Millisecond*100, mockFlusher)
			ctx, cancelFunc = context.WithCancel(ctx)
			mockFlusher.EXPECT().Flush([]models.Video{models.Video{
				VideoId: 1,
				SlideId: 42,
				Link:    "/video/42",
			}}).Return(nil, nil).Times(1)
			mockFlusher.EXPECT().Flush(gomock.Any()).Times(1)
		})

		JustBeforeEach(func() {
			err = s.Save(ctx, models.Video{
				VideoId: 1,
				SlideId: 42,
				Link:    "/video/42",
			})
		})

		It("drop", func() {
			<-time.After(time.Millisecond * 200)
			cancelFunc()

			Expect(err).Should(BeNil())
		})
	})

	Context("drop all on overflow test", func() {
		var (
			cancelFunc context.CancelFunc
		)

		BeforeEach(func() {
			s = saver.New(testCapacity, time.Millisecond*10, mockFlusher)
			ctx, cancelFunc = context.WithCancel(ctx)
			mockFlusher.EXPECT().Flush([]models.Video{models.Video{
				VideoId: 3,
				SlideId: 42,
				Link:    "/video/42",
			}}).Return(nil, nil)
		})

		JustBeforeEach(func() {
			err = s.Save(ctx, models.Video{
				VideoId: 1,
				SlideId: 42,
				Link:    "/video/42",
			})
			err = s.Save(ctx, models.Video{
				VideoId: 2,
				SlideId: 42,
				Link:    "/video/42",
			})
			err = s.Save(ctx, models.Video{
				VideoId: 3,
				SlideId: 42,
				Link:    "/video/42",
			})
			cancelFunc()
		})

		It("drop", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("drop single on overflow test", func() {
		var (
			cancelFunc context.CancelFunc
		)

		BeforeEach(func() {
			s = saver.New(testCapacity, time.Millisecond*10, mockFlusher, saver.DropFirstPolicy)
			ctx, cancelFunc = context.WithCancel(ctx)
			mockFlusher.EXPECT().Flush([]models.Video{models.Video{
				VideoId: 2,
				SlideId: 42,
				Link:    "/video/42",
			},
				models.Video{
					VideoId: 3,
					SlideId: 42,
					Link:    "/video/42",
				}}).Return(nil, nil)
		})

		JustBeforeEach(func() {
			err = s.Save(ctx, models.Video{
				VideoId: 1,
				SlideId: 42,
				Link:    "/video/42",
			})
			err = s.Save(ctx, models.Video{
				VideoId: 2,
				SlideId: 42,
				Link:    "/video/42",
			})
			err = s.Save(ctx, models.Video{
				VideoId: 3,
				SlideId: 42,
				Link:    "/video/42",
			})
			cancelFunc()
		})

		It("drop", func() {
			Expect(err).Should(BeNil())
		})
	})

	//TODO: test error in flusher.Flush
})
