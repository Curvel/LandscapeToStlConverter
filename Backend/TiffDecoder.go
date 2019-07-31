package main

import (
	Stl "LandscapeToStlConverter/Backend/lib"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/tiff"
	"image"
	"os"
)

type STRM struct {
	name   string
	top    float32
	right  float32
	bottom float32
	left   float32
	image  image.Image
}

var strmMaps = [...]STRM{
	{"srtm_38_03", 50.0, 10.0, 45.0, 5.0, nil},
	{"srtm_38_02", 55.0, 10.0, 50.0, 5.0, nil},
	{"srtm_39_02", 55.0, 15.0, 50.0, 10.0, nil},
	{"srtm_39_03", 50.0, 15.0, 45.0, 10.0, nil}}

func main() {
	// TODO change to console params
	var top float32 = 49.27
	var right float32 = 8.55
	var bottom float32 = 48.91
	var left float32 = 8.19

	heightMap := getHeightMap(top, right, bottom, left)

	Stl.GenerateSTLMapFromHeightMap(heightMap, 50)

}

func getMaxBorders() (maxTop float32, maxRight float32, maxBottom float32, maxLeft float32) {
	maxTop = -100000
	maxRight = -100000
	maxBottom = 100000
	maxLeft = 100000

	for _, strmMap := range strmMaps {
		if strmMap.top > maxTop {
			maxTop = strmMap.top
		}
		if strmMap.right > maxRight {
			maxRight = strmMap.right
		}
		if strmMap.left < maxLeft {
			maxLeft = strmMap.left
		}
		if strmMap.bottom < maxBottom {
			maxBottom = strmMap.bottom
		}
	}

	return
}

func isSelectionInRange(top float32, right float32, bottom float32, left float32) bool {
	maxTop, maxRight, maxBottom, maxLeft := getMaxBorders()

	return top <= maxTop && right <= maxRight && left >= maxLeft && bottom >= maxBottom
}

func getHeightMap(top float32, right float32, bottom float32, left float32) [][]float32 {

	if !isSelectionInRange(top, right, bottom, left) {
		panic("Selection out of range!")
	}

	loadImagesForRange(top, right, bottom, left)

	maxTop, maxRight, maxBottom, maxLeft := getMaxBorders()

	imgXPoints := 6000 * ((maxTop - maxBottom) / 5)
	imgYPoints := 6000 * ((maxRight - maxLeft) / 5)

	// Coordinates to SRTM-Map scales
	yScale := (maxTop - maxBottom) / imgYPoints
	xScale := (maxRight - maxLeft) / imgXPoints

	// Size of generated Height Map
	xSize := int( mgl32.Round((right - left) / xScale,0))
	ySize := int( mgl32.Round((top - bottom) / yScale,0))

	fmt.Printf("xSize: %d, ySize: %d\n", xSize, ySize)

	heightMap := make([][]float32, ySize)
	for i := range heightMap {
		heightMap[i] = make([]float32, xSize)
	}

	yOffset := int(-(maxBottom - bottom) / yScale)
	xOffset := int(-(maxLeft - left) / xScale)
	for yHeightMap := 0; yHeightMap < ySize; yHeightMap++ {
		for xHeightMap := 0; xHeightMap < xSize; xHeightMap++ {
			xImg := yOffset + yHeightMap
			yImg := xOffset + xHeightMap

			height := getHeight(xImg, yImg, xScale, yScale, maxLeft, maxBottom)

			heightMap[yHeightMap][xHeightMap] = float32(height)
		}
	}

	return heightMap
}

func getHeight(x int, y int, xScale float32, yScale float32, maxLeft float32, maxBottom float32) uint32 {
	imageNeeded := getNeededImage(x, y, xScale, yScale, maxLeft, maxBottom)

	r, _, _, _ := imageNeeded.At((y%6000), 5999 - (x%6000)).RGBA()
	if r > 10000 {
		r = 0
	}
	return r
}

func getNeededImage(x int, y int, xScale float32, yScale float32, maxLeft float32, maxBottom float32) (neededImage image.Image) {
	xCoordinate := float32(y)*yScale + maxLeft
	yCoordinate := float32(x)*xScale + maxBottom

	//fmt.Printf("xCoor: %f, yCoord: %f", xCoordinate, yCoordinate)
	for _, strmMap := range strmMaps {
		//fmt.Printf("top: %f, right: %f, bottom: %f, left: %f, %d \n", strmMap.top, strmMap.right, strmMap.bottom, strmMap.left, i)
		if strmMap.top > yCoordinate && strmMap.right > xCoordinate && strmMap.bottom <= yCoordinate && strmMap.left <= xCoordinate {
			return strmMap.image
		}
	}

	panic("Needed Image was not loaded!")
}

func loadImagesForRange(top float32, right float32, bottom float32, left float32) {
	for i, m := range strmMaps {
		if bottom <= m.top && top >= m.bottom && left <= m.right && right >= m.left {
			img := srtmTiffToImage(m.name)
			strmMaps[i].image = img
		}
	}
}

func srtmTiffToImage(name string) image.Image {
	uri := fmt.Sprintf("C:/Users/maxgt/go/src/LandscapeToStlConverter/Backend/srtm/%s/%s.tif", name, name)

	file, err := os.Open(uri)
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
