package main

import (
	Stl "./lib"
	"fmt"
	"golang.org/x/image/tiff"
	"image"
	"os"
)

func main() {
	img := tiffToImage()

	/*fmt.Println(img.At(0, 0))
	fmt.Println(img.At(0, 5999))
	fmt.Println(img.At(5999, 5999))
	fmt.Println(img.At(5999, 0))*/

	heightMap := getHeightMapOfImage(img, 49.0, 9.0, 47.0, 7.0)

	Stl.GenerateSTLMapFromHeightMap(heightMap, 5000)

}

func getHeightMapOfImage(img image.Image, upper float32, right float32, lower float32, left float32) [][]float32 {
	imgUpper := float32(50.0)
	imgRight := float32(10.0)
	imgLower := float32(45.0)
	imgLeft := float32(5.0)

	imgXPoints := float32(img.Bounds().Max.X)
	imgYPoints := float32(img.Bounds().Max.Y)

	yScale := (imgUpper - imgLower) / imgYPoints
	xScale := (imgRight - imgLeft) / imgXPoints

	xSize := int((right - left) / xScale)
	ySize := int((upper - lower) / yScale)

	fmt.Printf("xSize: %d, ySize: %d\n", xSize, ySize)

	heightMap := make([][]float32, ySize)
	for i := range heightMap {
		heightMap[i] = make([]float32, xSize)
	}

	yOffset := -int((imgLower - lower) / yScale)
	xOffset := -int((imgLeft - left) / xScale)
	for yHeightMap := 0; yHeightMap < ySize; yHeightMap++ {
		for xHeightMap := 0; xHeightMap < xSize; xHeightMap++ {
			xImg := yOffset + yHeightMap
			yImg := xOffset + xHeightMap

			if xImg > 5999 || xImg < 0 || yImg > 5999 || yImg < 0 {
				fmt.Print("")
			}

			r, _, _, _ := img.At(xImg, yImg).RGBA()
			heightMap[yHeightMap][xHeightMap] = float32(r)
		}
	}

	return heightMap
}

/* 	links := 0
	rechts := img.Bounds().Max.X - 1 //5999
	unten := img.Bounds().Max.Y - 1 //5999
 	oben := 0
	*Left 5.0000000 //Längengrad
	Lower 45.0000000 //Breitengrad
	Right 10.0000000 //Längengrad
	Upper 50.0000000 //Breitengrad

	1Pixel = 0.00083333333
*/
func tiffToImage() image.Image {
	file, err := os.Open("./srtm/srtm_38_03/srtm_38_03.tif")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	img, err := tiff.Decode(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return img
}
