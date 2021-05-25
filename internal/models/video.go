package models

import "fmt"

func NewVideo(videoId uint64, slideId uint64, link string) *video {
	return &video{videoId, slideId, link}
}

// FIXME: все id uint в соответствии с ER-диаграмой, полагаю что uint на ER диаграмме это uint64
type video struct {
	VideoId uint64
	SlideId uint64
	Link string
}

func (v *video) String() string {
	return fmt.Sprintf("Video [id = %v, slideId=%v, link=%v]", v.VideoId, v.SlideId, v.Link)
}