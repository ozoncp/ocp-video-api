package utils

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"ocp-video-api/internal/models"
	"testing"
)

func TestVideoIdToMap(t *testing.T) {
	type TestFixture struct {
		in   []models.Video
		want map[uint64]models.Video
	}

	testSuccessFixtures := []TestFixture{
		{
			[]models.Video{
				models.Video{VideoId: 1, SlideId: 2, Link: "test link 1"},
				models.Video{VideoId: 2, SlideId: 2, Link: "test link 1"},
			},
			map[uint64]models.Video{
				uint64(1): models.Video{VideoId: 1, SlideId: 2, Link: "test link 1"},
				uint64(2): models.Video{VideoId: 2, SlideId: 2, Link: "test link 1"},
			},
		},
	}

	for _, fixture := range testSuccessFixtures {
		got, err := SliceToMapVideo(fixture.in)
		if err != nil {
			t.Error(err)
		}
		if cmp.Equal(got, fixture.want) != true {
			t.Errorf(spew.Sprintf("got:\n%v\nwant:\n%+v", got, fixture.want))
		}
	}
}