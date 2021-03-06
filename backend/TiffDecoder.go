package main

import (
	Stl "./lib"
	"errors"
	"flag"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/tiff"
	"image"
	"os"
	"path/filepath"
	"runtime"
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

var topFlag = flag.Float64("neLat", 10000.0, "top coordinate of selected area")
var rightFlag = flag.Float64("neLng", 10000.0, "right coordinate of selected area")
var bottomFlag = flag.Float64("swLat", 10000.0, "bottom coordinate of selected area")
var leftFlag = flag.Float64("swLng", 10000.0, "left coordinate of selected area")
var modelTypeFlag = flag.String("model", "", "surface|section")
var croppingFlag = flag.String("cropping", "", "sqr|hex|rnd")
var lengthFlag = flag.Int("length", 0, "length of the largest side in mm")
var fileNameFlag = flag.String("name", "", "name of the stl file")
var heightFactorFlag = flag.Float64("heightFactor", 0.0, "smaller is bigger")

var projectPath string

func main() {
	flag.Parse()

	var top= float32(*topFlag)
	var right= float32(*rightFlag)
	var bottom= float32(*bottomFlag)
	var left= float32(*leftFlag)
	var modelType= *modelTypeFlag
	var cropping= *croppingFlag
	var length= *lengthFlag
	var fileName= *fileNameFlag
	var heightFactor= float32(*heightFactorFlag)

	//fmt.Printf("naLat: %f, neLang: %f, swLat: %f, swLng: %f, model: %s, cropping: %s, length: %d, name: %s, height: %f",
	//	top, right, bottom, left, modelType, cropping, length, fileName, heightFactor)

	var mErr error

	if runtime.GOOS == "darwin" {
		projectPath = "./"
	} else {
		projectPath, mErr = filepath.Abs("./")
	}


	if modelType == "surface" && mErr == nil {
		heightMap, err := getHeightMap(top, right, bottom, left)
		mErr = err
		if heightMap != nil && err == nil {
			if cropping == "sqr" {
				mErr = Stl.GenerateSTLMapFromHeightMap(heightMap, uint32(length), heightFactor, fileName)
			} else if cropping == "hex" {
				mErr = Stl.GenerateSettlerOfCatan(heightMap, uint32(length), heightFactor, fileName)
			}
		}
	} else if modelType == "section" && mErr == nil {
		profileMap, err := getProfileMap(top, right, bottom, left)
		mErr = err
		if profileMap != nil && err == nil {
			mErr = Stl.GenerateSTLMapFromSideMap(profileMap, uint32(length), heightFactor, fileName)
		}
	}

	if mErr == nil {
		os.Exit(0)
	} else {
		fmt.Print(mErr)
		os.Exit(1)
	}
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

func getProfileMap(top float32, right float32, bottom float32, left float32) ([]float32, error) {
	heightMap, err := getHeightMap(top, right, bottom, left)
	if err != nil {
		return nil, err
	}
	lengthX := float32(len(heightMap[0]))
	lengthY := float32(len(heightMap))
	var stepX float32
	var stepY float32
	var size int
	if (lengthY > lengthX) {
		size = int(lengthY)
		stepY = 1.0
		stepX = lengthX / lengthY
	} else {
		size = int(lengthX)
		stepX = 1.0
		stepY = lengthY / lengthX
	}
	profileMap := make([]float32, size)
	for i := 0; i < size-1; i++ {
		profileMap[i] = heightMap[int(float32(i)*stepY)][int(float32(i)*stepX)]
	}
	return profileMap, nil
}

func getHeightMap(top float32, right float32, bottom float32, left float32) ([][]float32, error) {
	//fmt.Println("Load HeightMap")

	if !isSelectionInRange(top, right, bottom, left) {
		return nil, errors.New("selection out of range")
	}

	fmt.Println("0;0;0")
	err := loadImagesForRange(top, right, bottom, left)
	if err != nil {
		return nil, err
	}
	fmt.Println("10;0;0")
	maxTop, maxRight, maxBottom, maxLeft := getMaxBorders()

	imgXPoints := 6000 * ((maxTop - maxBottom) / 5)
	imgYPoints := 6000 * ((maxRight - maxLeft) / 5)

	// Coordinates to SRTM-Map scales
	yScale := (maxTop - maxBottom) / imgYPoints
	xScale := (maxRight - maxLeft) / imgXPoints

	// Size of generated Height Map
	xSize := int(mgl32.Round((right-left)/xScale, 0))
	ySize := int(mgl32.Round((top-bottom)/yScale, 0))

	//fmt.Printf("xSize: %d, ySize: %d\n", xSize, ySize)

	heightMap := make([][]float32, ySize)
	for i := range heightMap {
		heightMap[i] = make([]float32, xSize)
	}

	yOffset := int(-(maxBottom - bottom) / yScale)
	xOffset := int(-(maxLeft - left) / xScale)
	for yHeightMap := 0; yHeightMap < ySize; yHeightMap++ {
		percentage := int(10 + ((float32(yHeightMap) / float32(ySize)) * 85))
		fmt.Printf("%d;0;0\n", percentage) // percent
		for xHeightMap := 0; xHeightMap < xSize; xHeightMap++ {
			xImg := yOffset + yHeightMap
			yImg := xOffset + xHeightMap

			height, err := getHeight(xImg, yImg, xScale, yScale, maxLeft, maxBottom)
			if err != nil {
				return nil, err
			}

			heightMap[yHeightMap][xHeightMap] = float32(height)
		}
	}

	fmt.Println("95;0;0")
	return flipMapX(heightMap), nil
}

func flipMapX(heightMap [][]float32) [][]float32 {
	flippedMap := make([][]float32, len(heightMap))
	for i := range flippedMap {
		flippedMap[i] = make([]float32, len(heightMap[0]))
	}
	for i := 0; i < len(heightMap); i++ {
		for j := 0; j < len(heightMap[0]); j++ {
			flippedMap[i][(len(heightMap[0])-1)-j] = heightMap[i][j]
		}
	}
	fmt.Println("100;0;0")
	return flippedMap
}

func getHeight(x int, y int, xScale float32, yScale float32, maxLeft float32, maxBottom float32) (uint32, error) {
	img, err := getNeededImage(x, y, xScale, yScale, maxLeft, maxBottom)
	if err != nil {
		return 0, err
	}

	r, _, _, _ := img.At((y % 6000), 5999-(x%6000)).RGBA()
	if r > 10000 {
		r = 0
	}
	return r, nil
}

func getNeededImage(x int, y int, xScale float32, yScale float32, maxLeft float32, maxBottom float32) (image.Image, error) {
	xCoordinate := float32(y)*yScale + maxLeft
	yCoordinate := float32(x)*xScale + maxBottom

	//fmt.Printf("xCoor: %f, yCoord: %f", xCoordinate, yCoordinate)
	for _, strmMap := range strmMaps {
		//fmt.Printf("top: %f, right: %f, bottom: %f, left: %f, %d \n", strmMap.top, strmMap.right, strmMap.bottom, strmMap.left, i)
		if strmMap.top > yCoordinate && strmMap.right > xCoordinate && strmMap.bottom <= yCoordinate && strmMap.left <= xCoordinate {
			return strmMap.image, nil
		}
	}

	return nil, errors.New("needed image was not loaded")
}

func loadImagesForRange(top float32, right float32, bottom float32, left float32) error {
	for i, m := range strmMaps {
		if bottom <= m.top && top >= m.bottom && left <= m.right && right >= m.left {
			img, err := srtmTiffToImage(m.name)
			if err != nil {
				return err
			}
			strmMaps[i].image = img
		}
	}
	return nil
}

func srtmTiffToImage(name string) (image.Image, error) {
	uri := fmt.Sprintf("%s/srtm/%s/%s.tif", projectPath, name, name)
	file, err := os.Open(uri)
	if err != nil {
		return nil, err
	}

	img, err := tiff.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
