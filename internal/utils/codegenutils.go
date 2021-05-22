package utils

import "github.com/cheekybits/genny/generic"

type TSliceUtils generic.Type

func SliceChunked(in []TSliceUtils, chunkSize uint) ([][]TSliceUtils, error) {
	panic("Not implemented!")
}

func SliceFilter(in []TSliceUtils, ban []TSliceUtils) ([][]TSliceUtils, error) {
	panic("Not implemented!")
}

type TMapKey generic.Type
type TMapValue generic.Type

func MapKVSwapped(in map[TMapKey]TMapValue) (map[TMapValue]TMapKey, error) {
	panic("Not implemented!")
}