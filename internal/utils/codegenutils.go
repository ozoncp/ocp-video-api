// Package utils provides building blocks for main ocp-video-api microservice
package utils

import (
	"github.com/cheekybits/genny/generic"
)

// TSliceUtils is the placeholder for genny to monomorphise for
// specialised slice types utils funcions
type TSliceUtils generic.Type

// SliceChunked return slice of slices, every inner slice points to original input slice mem
// in - input slice
// chunkSize - size of each chunk
// Panics if chunkSize <= 0!
func SliceChunked(in []TSliceUtils, chunkSize int) [][]TSliceUtils {
	if chunkSize <= 0 {
		panic("chunkSize must be >0")
	}

	retval := make([][]TSliceUtils, 0, (len(in) + chunkSize - 1) / chunkSize)
	for chunkSize < len(in) {
		in, retval = in[chunkSize:], append(retval, in[0 : chunkSize : chunkSize])
	}
	retval = append(retval, in)

	return retval
}

// SliceChunked return slice filtered from ban elements, filter is performed inplace!
// in - input slice
// ban - slice of elements which would be filtered
func SliceFilter(in []TSliceUtils, ban []TSliceUtils) []TSliceUtils {
	retval := in[:0]
	for _, v := range in {
		isBanned := false
		for _, banned := range ban {
			if v == banned {
				isBanned = true
				break
			}
		}
		if !isBanned {
			retval = append(retval, v)
		}
	}

	return retval
}

type TMapKey generic.Type
type TMapValue generic.Type

// MapKVSwapped return new map with flipped key to value
// in - input map
// Panics if input map contains value duplicates
func MapKVSwapped(in map[TMapKey]TMapValue) map[TMapValue]TMapKey {
	retval := make(map[TMapValue]TMapKey, len(in))
	for k, v := range in {
		if _, found := retval[v]; found {
			panic("input map should not contains duplicate values!")
		}
		retval[v] = k
	}

	return retval
}