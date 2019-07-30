package main

import (
	Stl "LandscapeToStlConverter/Backend/lib"
	"fmt"
	"golang.org/x/image/tiff"
	"image"
	"os"
)

func main() {
	img := tiffToImage()

	fmt.Println(img.At(0, 0))
	fmt.Println(img.At(0, 5999))
	fmt.Println(img.At(5999, 5999))
	fmt.Println(img.At(5999, 0))

	heightMap := getHeightMapOfImage(img, 50.0, 10.0, 45.0, 5.0)

	fmt.Println(heightMap[0][0])
	fmt.Println(heightMap[0][5999])
	fmt.Println(heightMap[5999][5999])
	fmt.Println(heightMap[5999][0])
	size := 6000
	heightMap2 := make([][]float32, size)
	for i := range heightMap2 {
		heightMap2[i] = make([]float32, size)
	}
	max := heightMap[0][0]
	for i:=0; i<size; i++{
		for j := 0 ; j < size; j++ {
			heightMap2[i][j] = heightMap[i][j]
			if heightMap2[i][j] < max{
				max = heightMap2[i][j]
			}
		}
	}
	for i:=0; i<size; i++{
		for j := 0 ; j < size; j++ {
			heightMap2[i][j] = heightMap2[i][j] - max
		}
	}

	fmt.Println("start generating")
	Stl.GenerateSTLMapFromHeightMap(heightMap2, 5000)

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

	for yHeightMap := 0; yHeightMap < ySize; yHeightMap++ {
		for xHeightMap := 0; xHeightMap < xSize; xHeightMap++ {
			xImg := int(float32(yHeightMap))
			yImg := int(float32(xHeightMap))

			if xImg > ySize || xImg < 0 || yImg > xSize || yImg < 0 {
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
	file, err := os.Open("C:/Users/maxgt/go/src/LandscapeToStlConverter/Backend/srtm/srtm_38_03/srtm_38_03.tif")
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
