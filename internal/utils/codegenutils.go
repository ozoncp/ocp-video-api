package utils

import (
	"fmt"
	"github.com/cheekybits/genny/generic"
)

type TSliceUtils generic.Type

func SliceChunked(in []TSliceUtils, chunkSize int) ([][]TSliceUtils, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("chunkSize must be >0")
	}

	retval := make([][]TSliceUtils, 0, (len(in) + chunkSize - 1) / chunkSize)
	for chunkSize < len(in) {
		in, retval = in[chunkSize:], append(retval, in[0 : chunkSize : chunkSize])
	}
	retval = append(retval, in)

	return retval, nil
}

func SliceFilter(in []TSliceUtils, ban []TSliceUtils) ([]TSliceUtils, error) {
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

	return retval, nil
}

type TMapKey generic.Type
type TMapValue generic.Type

func MapKVSwapped(in map[TMapKey]TMapValue) (map[TMapValue]TMapKey, error) {
	retval := make(map[TMapValue]TMapKey, len(in))
	for k, v := range in {
		if _, found := retval[v]; found {
			panic("input map should not contains duplicate values!")
		}
		retval[v] = k
	}

	return retval, nil
}