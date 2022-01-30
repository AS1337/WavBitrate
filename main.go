package main

import (
	"fmt"

	"github.com/AS1337/WavBitrate/bitrate"
)

func testBitrate() {
	a, err := bitrate.CheckBitrate("a.wav")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}

func main() {
	testBitrate()
}
