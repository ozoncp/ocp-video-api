package utils

import (
	"fmt"
	"ocp-video-api/internal/models"
)

//TODO: переписать через слайс интерфейсов Ider { Id() uint64 } позже?
// через индирекцию будет медленне чем через специализированную функцию?
// добавить для генерализации в аргументы функтор select который извлекает
// ключ мапы из значения и мономорфизировать?
func SliceToMapVideo(in []models.Video) (map[uint64]models.Video, error) {
	if in == nil {
		return nil, fmt.Errorf("empty slice provided")
	}

	retval := make(map[uint64]models.Video, len(in))
	for _, v := range in {
		if _, found := retval[v.VideoId]; found {
			return nil, fmt.Errorf("input slice contains video with duplicate ids")
		}
		retval[v.VideoId] = v
	}

	return retval, nil
}
