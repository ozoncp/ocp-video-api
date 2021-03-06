package flusher_test

import (
	"context"
	"errors"
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
			mockRepo.EXPECT().AddVideos(gomock.Any(), gomock.Any()).Return(uint64(1), nil).Times(1)
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
			mockRepo.EXPECT().AddVideos(gomock.Any(), gomock.Any()).Return(uint64(2), nil).Times(1)
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
			mockRepo.EXPECT().AddVideos(gomock.Any(), gomock.Any()).Return(uint64(2), nil).Times(1)
			mockRepo.EXPECT().AddVideos(gomock.Any(), gomock.Any()).Return(uint64(1), nil).Times(1)
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
			mockRepo.EXPECT().AddVideos(gomock.Any(), gomock.Any()).Return(uint64(0), nil).Times(0)
			rest, err := f.Flush(in)
			Expect(rest).Should(BeEquivalentTo(in))
			Expect(err).Should(BeIdenticalTo(internal.ErrInvalidSize))
		})

		It("returns original slice and error on nil, passed as slice if chunk size is ok", func() {
			f = flusher.New(2, mockRepo)
			mockRepo.EXPECT().AddVideos(gomock.Any(), gomock.Any()).Return(uint64(0), nil).Times(0)
			rest, err := f.Flush(nil)
			Expect(rest).Should(BeNil())
			Expect(err).Should(BeNil())
		})
	})

	Context("Flushes partially what is already processe on repo.AddVideos error", func() {
		It("", func() {
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
			expected := []models.Video{
				models.Video{
					VideoId: 3,
					SlideId: 2,
					Link:    "video/3",
				},
			}
			someErr := errors.New("some error")
			f = flusher.New(1, mockRepo)

			idx, failOn := 0, 2
			mockRepo.EXPECT().AddVideos(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, vs []models.Video) (uint64, error) {
					if idx == failOn {
						return uint64(0), someErr
					}
					idx++
					return uint64(1), nil
				},
			).Times(3)

			rest, err := f.Flush(in)
			Expect(rest).Should(BeEquivalentTo(expected))
			Expect(err).Should(BeIdenticalTo(someErr))
		})
	})
})
