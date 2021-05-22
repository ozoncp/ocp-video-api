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

func SliceFilter(in []TSliceUtils, ban []TSliceUtils) ([][]TSliceUtils, error) {
	panic("Not implemented!")
}

type TMapKey generic.Type
type TMapValue generic.Type

func MapKVSwapped(in map[TMapKey]TMapValue) (map[TMapValue]TMapKey, error) {
	panic("Not implemented!")
}