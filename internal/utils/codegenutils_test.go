package utils

import (
	"github.com/google/go-cmp/cmp"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

//TODO: протестировать граничные случаи
func TestSliceChunkedTSliceUtils(t *testing.T) {
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
		{ 2,
			[]TSliceUtils{1, 2, 3},
			[][]TSliceUtils {
				[]TSliceUtils{ 1, 2 },
				[]TSliceUtils{ 3 },
			},
		},
		{ 3,
			[]TSliceUtils{1, 2, 3},
			[][]TSliceUtils {
				[]TSliceUtils{ 1, 2, 3 },
			},
		},
		{ 4,
			[]TSliceUtils{1, 2, 3},
			[][]TSliceUtils {
				[]TSliceUtils{ 1, 2, 3 },
			},
		},
	}

	for _, fixture := range testFixtures {
		got := SliceChunkedTSliceUtils(fixture.in, fixture.chunkSize)
		if cmp.Equal(got, fixture.want) == false {
			t.Errorf(spew.Sprintf("got:\n%v\nwant:\n%+v", got, fixture.want))
		}
	}
}

//TODO: протестировать граничные случаи
func TestSliceFilterTSliceUtils(t *testing.T) {
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
		{ []TSliceUtils{1, 1, 2, 3},
			[]TSliceUtils{1, 2, 2},
			[]TSliceUtils{3},
		},
		{ []TSliceUtils{1, 1, 1, 1},
			[]TSliceUtils{1, 2, 2},
			[]TSliceUtils{},
		},
		{ []TSliceUtils{1, 1, 1, 1},
			[]TSliceUtils{},
			[]TSliceUtils{1, 1, 1, 1},
		},
	}

	for _, fixture := range testFixtures {
		got := SliceFilterTSliceUtils(fixture.in, fixture.ban)
		if cmp.Equal(got, fixture.want) != true {
			t.Errorf(spew.Sprintf("got:\n%v\nwant:\n%+v", got, fixture.want))
		}
	}
}

//TODO: протестировать граничные случаи
//TODO: протестировать панику
func TestMapKTMapKeyVTMapValueSwapped(t *testing.T) {
	type TestFixture struct {
		in map[TMapKey]TMapValue
		want map[TMapValue]TMapKey
	}

	testFixtures := []TestFixture {
		{ map[TMapKey]TMapValue { 1: 10 },
			map[TMapValue]TMapKey { 10: 1 },
		},
		{ map[TMapKey]TMapValue { 1: 10, 2: 20 },
			map[TMapValue]TMapKey { 10: 1, 20: 2 },
		},
	}

	for _, fixture := range testFixtures {
		got := MapKTMapKeyVTMapValueSwapped(fixture.in)
		if cmp.Equal(got, fixture.want) != true {
			t.Errorf(spew.Sprintf("got:\n%v\nwant:\n%+v", got, fixture.want))
		}
	}
}