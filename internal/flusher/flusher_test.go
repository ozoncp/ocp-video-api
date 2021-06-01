package flusher_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"ocp-video-api/internal"
	"ocp-video-api/internal/flusher"
	"ocp-video-api/internal/mock"
	"ocp-video-api/internal/models"
)

var _ = Describe("Flusher", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mock.MockRepo
		f        flusher.Flusher
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mock.NewMockRepo(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Flusher successfully saves videos less than batch size len", func() {
		in := []models.Video{
			models.Video{
				VideoId: 1,
				SlideId: 1,
				Link:    "video/1",
			},
		}
		It("", func() {
			f = flusher.New(2, mockRepo)
			mockRepo.EXPECT().AddVideos(gomock.Any()).Return(nil).MinTimes(1).MaxTimes(1)
			rest, err := f.Flush(in)
			Expect(rest).Should(BeNil())
			Expect(err).Should(BeNil())
		})
	})

	Context("Flusher successfully saves videos exact as batch size len", func() {
		in := []models.Video{
			models.Video{
				VideoId: 1,
				SlideId: 1,
				Link:    "video/1",
			},
			models.Video{
				VideoId: 2,
				SlideId: 1,
				Link:    "video/2",
			},
		}
		It("", func() {
			f = flusher.New(2, mockRepo)
			mockRepo.EXPECT().AddVideos(gomock.Any()).Return(nil).MinTimes(1).MaxTimes(1)
			rest, err := f.Flush(in)
			Expect(rest).Should(BeNil())
			Expect(err).Should(BeNil())
		})
	})

	Context("Flusher successfully saves videos more than batch size len", func() {
		in := []models.Video{
			models.Video{
				VideoId: 1,
				SlideId: 1,
				Link:    "video/1",
			},
			models.Video{
				VideoId: 2,
				SlideId: 1,
				Link:    "video/2",
			},
			models.Video{
				VideoId: 3,
				SlideId: 2,
				Link:    "video/3",
			},
		}
		It("", func() {
			f = flusher.New(2, mockRepo)
			mockRepo.EXPECT().AddVideos(gomock.Any()).Return(nil).MinTimes(2).MaxTimes(2)
			rest, err := f.Flush(in)
			Expect(rest).Should(BeNil())
			Expect(err).Should(BeNil())
		})
	})

	Context("Flusher return error on insufficient arguments", func() {
		It("returns original slice and error on batch size equal to 0", func() {
			in := []models.Video{
				models.Video{
					VideoId: 1,
					SlideId: 1,
					Link:    "video/1",
				},
				models.Video{
					VideoId: 2,
					SlideId: 1,
					Link:    "video/2",
				},
				models.Video{
					VideoId: 3,
					SlideId: 2,
					Link:    "video/3",
				},
			}

			f = flusher.New(0, mockRepo)
			mockRepo.EXPECT().AddVideos(gomock.Any()).Return(internal.ErrInvalidSize).MinTimes(1).MaxTimes(1)
			rest, err := f.Flush(in)
			Expect(rest).Should(BeEquivalentTo(in))
			Expect(err).Should(BeIdenticalTo(internal.ErrInvalidSize))
		})

		It("returns original slice and error on nil, passed as slice if chunk size is ok", func() {
			f = flusher.New(2, mockRepo)
			mockRepo.EXPECT().AddVideos(gomock.Any()).Return(internal.ErrInvalidArg).MinTimes(1).MaxTimes(1)
			rest, err := f.Flush(nil)
			Expect(rest).Should(BeNil())
			Expect(err).Should(BeIdenticalTo(internal.ErrInvalidArg))
		})

		It("returns original slice and error on empty slice, passed as slice if chunk size is ok", func() {
			in := []models.Video{{}}
			f = flusher.New(2, mockRepo)
			mockRepo.EXPECT().AddVideos(gomock.Any()).Return(internal.ErrInvalidArg).MinTimes(1).MaxTimes(1)
			rest, err := f.Flush(in)
			Expect(rest).Should(BeEquivalentTo([]models.Video{{}}))
			Expect(err).Should(BeIdenticalTo(internal.ErrInvalidArg))
		})
	})
})
