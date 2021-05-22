// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Package utils provides building blocks for main ocp-video-api microservice
package utils

// int is the placeholder for genny to monomorphise for

// SliceChunkedInt return slice of slices, every inner slice points to original input slice mem
// in - input slice
// chunkSize - size of each chunk
// Panics if chunkSize <= 0!
func SliceChunkedInt(in []int, chunkSize int) [][]int {
	if chunkSize <= 0 {
		panic("chunkSize must be >0")
	}

	retval := make([][]int, 0, (len(in)+chunkSize-1)/chunkSize)
	for chunkSize < len(in) {
		in, retval = in[chunkSize:], append(retval, in[0:chunkSize:chunkSize])
	}
	retval = append(retval, in)

	return retval
}

// SliceFilterInt return slice filtered from ban elements, filter is performed inplace!
// in - input slice
// ban - slice of elements which would be filtered
func SliceFilterInt(in []int, ban []int) []int {
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

// MapKIntVIntSwapped return new map with flipped key to value
// in - input map
// Panics if input map contains value duplicates
func MapKIntVIntSwapped(in map[int]int) map[int]int {
	retval := make(map[int]int, len(in))
	for k, v := range in {
		if _, found := retval[v]; found {
			panic("input map should not contains duplicate values!")
		}
		retval[v] = k
	}

	return retval
}