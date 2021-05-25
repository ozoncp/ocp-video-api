package main

import (
	"fmt"
	"ocp-video-api/internal/utils"
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

	fmt.Print(DESCRIPTION)
}
