package models

import "fmt"

// FIXME: все id uint в соответствии с ER-диаграмой, полагаю что uint на ER диаграмме это uint64
type Video struct {
	VideoId uint64
	SlideId uint64
	Link    string
}

func (v *Video) String() string {
	return fmt.Sprintf("Video [id = %v, slideId=%v, link=%v]", v.VideoId, v.SlideId, v.Link)
}
