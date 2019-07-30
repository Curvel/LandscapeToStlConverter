package main

import (
	"fmt"
	"golang.org/x/image/tiff"
	"os"
)

func main() {
	file, err := os.Open("./srtm/srtm_38_03/srtm_38_03.tif")
	if err != nil {
		fmt.Println(err)
		return
	}

	img, err := tiff.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(img.At(0, 300))
}
