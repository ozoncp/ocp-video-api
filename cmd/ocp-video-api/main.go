package main

import (
	"fmt"
	"ocp-video-api/internal/models"
	"ocp-video-api/internal/utils"
	"os"
)

const DESCRIPTION = `
Video API provider
Author: Mosolov Sergey
Start date: 2021/05/13

`

func main() {
	in := map[int]int{1: 10, 2: 20}
	flipped := utils.MapKIntVIntSwapped(in)
	fmt.Println(flipped)

	loopFileOpenClose()

	v := models.Video{VideoId: 1, SlideId: 1, Link: "link_to_video"}
	print("Video id: ", v.VideoId, " to_string: ", v.String())

	fmt.Print(DESCRIPTION)
}

// TODO: add tests /proc/<PID>/fd/... single text.txt opened file on each iteration
func loopFileOpenClose() {
	const FileDeferForTimes = 3
	for i := 0; i < FileDeferForTimes; i++ {
		func() {
			f, err := os.Open("/tmp/test.txt")
			if err != nil {
				panic(err)
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					panic(err)
				}
			}(f)
		}()
	}
}
