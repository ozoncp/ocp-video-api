package utils

import (
	"github.com/google/go-cmp/cmp"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestSliceChunked(t *testing.T) {
	type TestFixture struct {
		chunkSize int
		in []TSliceUtils
		want [][]TSliceUtils
	}

	testFixtures := []TestFixture {
		{ 1,
			[]TSliceUtils{1, 2, 3},
			[][]TSliceUtils {
				[]TSliceUtils{ 1 },
				[]TSliceUtils{ 2 },
				[]TSliceUtils{ 3 },
			},
		},
	}

	for _, fixture := range testFixtures {
		got, _ := SliceChunked(fixture.in, fixture.chunkSize)
		if cmp.Equal(got, fixture.want) == false {
			t.Errorf(spew.Sprintf("got:\n%v\nwant:\n%+v", got, fixture.want))
		}
	}
}

func TestSliceFilter(t *testing.T) {
	type TestFixture struct {
		in []TSliceUtils
		ban []TSliceUtils
		want []TSliceUtils
	}

	testFixtures := []TestFixture {
		{ []TSliceUtils{1, 1, 2, 3},
			[]TSliceUtils{1, 2},
			[]TSliceUtils{3},
		},
	}

	for _, fixture := range testFixtures {
		got, _ := SliceFilter(fixture.in, fixture.ban)
		if cmp.Equal(got, fixture.want) != true {
			t.Errorf("Error")
		}
	}
}

func TestMapKVSwapped(t *testing.T) {
	type TestFixture struct {
		in map[TMapKey]TMapValue
		want map[TMapValue]TMapKey
	}

	testFixtures := []TestFixture {
		{ map[TMapKey]TMapValue { 1: 10 },
			map[TMapValue]TMapKey { 1: 10 },
		},
	}

	for _, fixture := range testFixtures {
		got, _ := MapKVSwapped(fixture.in)
		if cmp.Equal(got, fixture.want) != true {
			t.Errorf("Error")
		}
	}
}