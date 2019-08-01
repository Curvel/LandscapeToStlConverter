package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/geo/r3"
	"math"
	//"golang.org/x/image/math/f32"
	//"encoding/binary"
	. "github.com/go-gl/mathgl/mgl32"
	//"github.com/golang/geo/r3"
	"log"
	"os"
)

const MapHeight = 2.0
const ProfileHeight = 5.0
const ProfileThickness = 2.0

type triangle struct {
	v1 Vec3
	v2 Vec3
	v3 Vec3
}


/*func GenerateSTLMapFromSideMap(sidemap []float32, sizeInMM uint32){
	var step float32
	size :=float32(sizeInMM)
	step = size / float32(len(sidemap)-1)
	heightstep := step / heightFactor
	c1:= Vec3{0,0,0}
	c2:= Vec3{ProfileHeight,0, 0}
	c3:= Vec3{ProfileHeight,size, 0}
	c4:= Vec3{0,size, 0}
	c5:= Vec3{0,0, ProfileThickness}
	c6:= Vec3{ProfileHeight,0, ProfileThickness}
	c7:= Vec3{ProfileHeight,size, ProfileThickness}
	c8:= Vec3{0,size, ProfileThickness}

	ct1 := triangle{c1, c2,c3}
	ct2 := triangle{c1, c3,c4}
	ct3 := triangle{c5, c6,c7}
	ct4 := triangle{c5, c7,c8}

	ct5 := triangle{c1, c5,c6}
	ct6 := triangle{c1, c6,c2}
	ct7 := triangle{c3, c7,c8}
	ct8 := triangle{c3, c8,c4}
	ct9 := triangle{c4, c8,c5}
	ct10 := triangle{c4, c5,c1}

	var triangles []triangle

	triangles = append(triangles, ct1,ct2,ct3,ct4,ct5,ct6,ct7,ct8,ct9,ct10)
	for i:= 0 ; i< len(sidemap); i++ {
		v1 := Vec3{sidemap[i]*heightstep+ProfileHeight, float32(i) * step, 0}
		v3 := Vec3{sidemap[i]*heightstep+ProfileHeight, float32(i) * step, ProfileThickness}
		v5 := Vec3{ProfileHeight, float32(i) * step,0}
		v7 := Vec3{ProfileHeight, float32(i) * step, ProfileThickness}
		if (i != len(sidemap)-1){
			v2 := Vec3{sidemap[i+1]*heightstep+ProfileHeight, float32(i+1) * step , 0}
			v4 := Vec3{sidemap[i+1]*heightstep+ProfileHeight, float32(i+1) * step, ProfileThickness}
			v6 := Vec3{ProfileHeight, float32(i+1) * step, 0}
			v8 := Vec3{ProfileHeight, float32(i+1) * step, ProfileThickness}

			t1 := triangle{v1,v3,v4}
			t2 := triangle{v1,v4,v2}
			t3 := triangle{v1,v2,v6}
			t4 := triangle{v1,v6,v5}
			t5 := triangle{v3,v4,v8}
			t6 := triangle{v3,v8,v7}
			triangles = append(triangles, t1,t2,t3,t4,t5,t6)
		}
		if(i == 0){
			ct1 := triangle{v1, v5, v7}
			ct2 := triangle{v1, v7,v3}
			triangles = append(triangles, ct1,ct2)
		}
		if (i == len(sidemap) -1){
			ct3:= triangle{v1,v3,v7}
			ct4:= triangle{v1,v7,v5}
			triangles = append (triangles, ct3,ct4)
		}
	}
	generateSTLMapFromTriangles(triangles)
}*/

func GenerateSTLMapFromSideMap(sideMap []float32, sizeInMM uint32, heightFactor float32, fileName string) {
	thickness :=  int((ProfileThickness/ float32(sizeInMM)) * float32(len(sideMap)))

	heightMap := make([][]float32, thickness)
	for i:= range heightMap{
		heightMap[i] = make ([]float32, len(sideMap))
	}
	for i := 0 ; i < thickness; i++ {
		for j := 0 ; j < len(sideMap); j++ {
			heightMap[i][j] = sideMap[j]
		}
	}

	GenerateSTLMapFromHeightMap(heightMap, sizeInMM, heightFactor, fileName)

}
func GenerateSettlerOfCatan(heightMap [][]float32, sizeInMM uint32, heightFactor float32, fileName string){
	var offset int
	var triangles []triangle
	var stepX float32
	var stepY float32
	size := float32(sizeInMM)

	heightMap = getSquareMap(heightMap)
	squareLength := len(heightMap)
	sideLength := int((float64(squareLength)/2)/ math.Sin((100.0/360.0) * 2.0 * math.Pi)) // keine Ahnung was hier genau rein muss
	startPosRowOne := (squareLength - sideLength) / 2
	startPosRel := size * float32(startPosRowOne)/float32(squareLength)
	stepsPerRow := float32(startPosRowOne)/(float32(squareLength)/2)

	stepX = size / float32(squareLength)
	stepY = size / float32(squareLength)
	heightstep := stepX /heightFactor

	c1 := Vec3{0,startPosRel, 0}
	c2 := Vec3{0,size- startPosRel, 0}
	c3 := Vec3{size/2,size, 0}
	c4 := Vec3{size-stepX,size-startPosRel, 0}
	c5 := Vec3{size-stepX,startPosRel , 0}
	c6 := Vec3{size/2,0,0}

	c7 := Vec3{0,startPosRel, MapHeight}
	c8 := Vec3{0,size- startPosRel, MapHeight}
	c9 := Vec3{size/2,size, MapHeight}
	c10 := Vec3{size-stepX,size-startPosRel, MapHeight}
	c11 := Vec3{size-stepX,startPosRel , MapHeight}
	c12 := Vec3{size/2,0,MapHeight}

	t1 := triangle{c6, c1,c2}
	t2 := triangle{c6, c2,c3}
	t3 := triangle{c6, c3,c4}
	t4 := triangle{c6, c4,c5}

	t11 := triangle{c1,c8,c2}
	t12 := triangle{c1,c7,c8}
	t21 := triangle{c2,c9,c3}
	t22 := triangle{c2,c8,c9}
	t31 := triangle{c3,c10,c4}
	t32 := triangle{c3,c9,c10}
	t41 := triangle{c4,c11,c5}
	t42 := triangle{c4,c10,c11}
	t51 := triangle{c5,c12,c6}
	t52 := triangle{c5,c11,c12}
	t61 := triangle{c6,c7,c1}
	t62 := triangle{c6,c12,c7}

	triangles = append(triangles, t1,t2,t3,t4,t11,t12,t21,t22,t31,t32,t41,t42,t51,t52,t61,t62)


	for x := 0; x < squareLength; x++ {
		percentage:= int((float32(x)/float32(squareLength))*100)
		fmt.Printf("100;%d;0\n", percentage)
		offset  = getCatanOffset(x, squareLength, stepsPerRow)
		for y:= startPosRowOne-offset; y < startPosRowOne + sideLength + offset ; y++{
			v1 := Vec3{float32(x) * stepX, float32(y) * stepY, heightMap[x][y]*heightstep + MapHeight}
			if x < squareLength - 1 && y < startPosRowOne + sideLength + offset -1 {
				v2 := Vec3{float32(x+1) * stepX, float32(y) * stepY, heightMap[x+1][y]*heightstep + MapHeight}
				v3 := Vec3{float32(x) * stepX, float32(y+1) * stepY, heightMap[x][y+1]*heightstep + MapHeight}
				v4 := Vec3{float32(x+1) * stepX, float32(y+1) * stepY, heightMap[x+1][y+1]*heightstep + MapHeight}
				t1 := triangle{v1, v2, v4}
				t2 := triangle{v1, v4, v3}
				triangles = append(triangles, t1, t2)
			}
			if y == startPosRowOne-offset{ // links
			//oben
				if (x!= 0 && x <= squareLength/2 && y < startPosRowOne - getCatanOffset(x-1,squareLength ,stepsPerRow )){
					xoff := 1
					for{
						if (!(y+1 >= startPosRowOne - getCatanOffset(x-xoff,squareLength, stepsPerRow ))|| x -xoff < 0){
							xoff =xoff-1
							break
						}
						xoff = xoff+1
					}

					vr := Vec3{float32(x) * stepX, float32(y+1) * stepY, heightMap[x][y+1] * heightstep + MapHeight}
					vu := Vec3{float32(x-xoff) * stepX, float32(y+1)* stepY, heightMap[x-xoff][y+1] * heightstep + MapHeight}
					vd := Vec3{float32(x) * stepX, startPosRel-((float32(x)/(float32(squareLength)/2)) * startPosRel), MapHeight}
					vud := Vec3{float32(x-xoff) * stepX, startPosRel -((float32(x-xoff)/(float32(squareLength)/2)) * startPosRel), MapHeight}

					tu := triangle{v1, vr,vu}
					td1 := triangle{v1, vud, vd}
					td2 := triangle{v1, vu, vud}
					triangles = append(triangles, tu, td1, td2)

				}
				//unten
				if (x > squareLength/2 && y < startPosRowOne - getCatanOffset(x+1, squareLength, stepsPerRow )|| x == squareLength-2){
					xoff := 0
					for{
						if ((y-1 >= startPosRowOne - getCatanOffset(x-xoff, squareLength, stepsPerRow ))|| x - xoff <= squareLength/2){
							xoff = xoff-1
							break
						}
						xoff = xoff+1
					}
					vb := Vec3{float32(x+1) * stepX, float32(y) * stepY, heightMap[x+1][y] * heightstep + MapHeight}
					vl := Vec3{float32(x-xoff) * stepX, float32(y-1) * stepY, heightMap[x-xoff][y-1] * heightstep + MapHeight}
					vu := Vec3{float32(x-xoff) * stepX, float32(y)* stepY, heightMap[x-xoff][y] * heightstep + MapHeight}
					vd := Vec3{float32(x+1) * stepX, ((float32(x+1)-float32(squareLength)/2)/(float32(squareLength)/2)) * startPosRel, MapHeight}
					vud := Vec3{float32(x-xoff) * stepX, ((float32(x-xoff)-float32(squareLength)/2)/(float32(squareLength)/2)) * startPosRel, MapHeight}

					tu := triangle{vb, vu,vl}
					td1 := triangle{vb, vl, vud}
					td2 := triangle{vb, vud, vd}
					triangles = append(triangles, tu , td1, td2)
				}
			}
			if(x == squareLength/2 && squareLength % 2 ==0) {
				if y == 0 {
					vb := Vec3{float32(x+1) * stepX, 0, heightMap[x+1][y]*heightstep + MapHeight}
					vd := Vec3{float32(x) * stepX, startPosRel - ((float32(x) / (float32(squareLength) / 2)) * startPosRel), MapHeight}
					vbd := Vec3{float32(x+1) * stepX, (float32(x+1)/(float32(squareLength)/2) - 1) * startPosRel, MapHeight}
					veck := Vec3{size / 2, 0, MapHeight}
					td1 := triangle{v1, vd, veck}
					td2 := triangle{v1, veck, vb}
					td3 := triangle{vb, veck, vbd}
					triangles = append(triangles, td1, td2, td3)
				}
				if y == squareLength {
					vb := Vec3{float32(x+1) * stepX, size, heightMap[x+1][y]*heightstep + MapHeight}
					vd := Vec3{float32(x) * stepX, size - (startPosRel -((float32(x) / (float32(squareLength) / 2)) * startPosRel)), MapHeight}
					vbd := Vec3{float32(x+1) * stepX, size - (float32(x+1)/(float32(squareLength)/2) - 1) * startPosRel, MapHeight}
					veck := Vec3{size / 2, size, MapHeight}

					td1:= triangle {v1,veck,vd}
					td2:= triangle {v1,vb,veck}
					td3:= triangle {vb,vbd,veck}
					triangles = append(triangles, td1, td2, td3)
				}
			}





			if y == (startPosRowOne + sideLength + offset - 1){ // rechts
			//oben
				if (x!= 0 && x <= squareLength/2 && y >= startPosRowOne + sideLength + getCatanOffset(x-1,squareLength ,stepsPerRow )){
					xoff := 1
					for{
						if (y-1 >= startPosRowOne +sideLength+ getCatanOffset(x-xoff,squareLength, stepsPerRow )|| x -xoff < 0){
							xoff =xoff-1
							break
						}
						xoff = xoff+1
					}

					vl := Vec3{float32(x) * stepX, float32(y-1) * stepY, heightMap[x][y-1] * heightstep + MapHeight}
					vu := Vec3{float32(x-xoff) * stepX, float32(y-1)* stepY, heightMap[x-xoff][y-1] * heightstep + MapHeight}
					vd := Vec3{float32(x) * stepX, size - (startPosRel-((float32(x)/(float32(squareLength)/2)) * startPosRel)), MapHeight}
					vud := Vec3{float32(x-xoff) * stepX, size-(startPosRel -((float32(x-xoff)/(float32(squareLength)/2)) * startPosRel)), MapHeight}

					tu := triangle{v1, vu,vl}
					td1 := triangle{v1, vd, vud}
					td2 := triangle{v1, vud, vu}
					triangles = append(triangles, tu , td1, td2)

				}
				//unten
				if (x > squareLength/2 && y >= startPosRowOne + sideLength + getCatanOffset(x+1, squareLength, stepsPerRow )|| x == squareLength-2){
					xoff := 0
					for{
						if ((y+1 < startPosRowOne + sideLength + getCatanOffset(x-xoff, squareLength, stepsPerRow ))|| x - xoff <= squareLength/2){
							break
						}
						xoff = xoff+1
					}
					vb := Vec3{float32(x+1) * stepX, float32(y) * stepY, heightMap[x+1][y] * heightstep + MapHeight}
					vr := Vec3{float32(x-xoff) * stepX, float32(y+1) * stepY, heightMap[x-xoff][y-1] * heightstep + MapHeight}
					vu := Vec3{float32(x-xoff) * stepX, float32(y)* stepY, heightMap[x-xoff][y] * heightstep + MapHeight}
					vd := Vec3{float32(x+1) * stepX, size- ((float32(x+1)-float32(squareLength)/2)/(float32(squareLength)/2)) * startPosRel, MapHeight}
					vud := Vec3{float32(x-xoff) * stepX, size - ((float32(x-xoff)-float32(squareLength)/2)/(float32(squareLength)/2)) * startPosRel, MapHeight}

					tu := triangle{vb, vr,vu}
					td1 := triangle{vb, vud, vr}
					td2 := triangle{vb, vd, vud}
					triangles = append(triangles, tu , td1, td2)
				}
			}






			if x == 0 && heightMap[x][y] != 0 {
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if y < startPosRowOne + sideLength + offset-1 {
					vr := Vec3{v1.X(), v1.Y() + stepY, MapHeight}
					triangles = append(triangles, triangle{v1, vd, vr})
				}
				if y > startPosRowOne-offset {
					vl := Vec3{v1.X(), v1.Y() - stepY, heightMap[x][y-1]*heightstep + MapHeight}
					triangles = append(triangles, triangle{v1, vl, vd})
				}
			}
			if x == squareLength-1 && heightMap[x][y] != 0 {
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if y < startPosRowOne + sideLength + offset-1 {
					vr := Vec3{v1.X(), v1.Y() + stepY, MapHeight}
					triangles = append(triangles, triangle{v1, vr, vd})
				}
				if y > startPosRowOne-offset {
					vl := Vec3{v1.X(), v1.Y() - stepY, heightMap[x][y-1]*heightstep + MapHeight} //falsch
					triangles = append(triangles, triangle{v1, vd, vl})
				}
			}

		}
	}

	generateSTLMapFromTriangles(triangles, fileName)

}

func getCatanOffset(x int, squareLength int, stepsPerRow float32) int{
	if x < squareLength/2 {
		return int(float32(x)* stepsPerRow )
	}else{
		return int(float32(squareLength-x) * stepsPerRow)
	}
}

func GenerateSTLMapFromHeightMap(heightMap [][]float32, sizeInMM uint32, heightFactor float32, fileName string){
	var stepX float32
	var stepY float32
	var sizeX float32
	var sizeY float32
	size := float32(sizeInMM)
	if (len(heightMap) > len(heightMap[0])){
		sizeX = size
		sizeY = size * (float32(len(heightMap[0])) / float32(len(heightMap)))
	} else {
		sizeY = size
		sizeX = size * (float32(len(heightMap)) / float32(len(heightMap[0])))
	}
	c1 := Vec3{0, 0, 0}               //bottom up left
	c2 := Vec3{sizeX, 0, 0}            //bottom up right
	c3 := Vec3{sizeX, sizeY, 0}         // bottom down right
	c4 := Vec3{0, sizeY, 0}            //bottom down left
	c5 := Vec3{0, 0, MapHeight}       // top up left
	c6 := Vec3{sizeX, 0, MapHeight}    //top up right
	c7 := Vec3{sizeX, sizeY, MapHeight} //top down right
	c8 := Vec3{0, sizeY, MapHeight}    //top down left

	ct1 := triangle{c1, c2, c3}  //bottom
	ct2 := triangle{c1, c3, c4}  //bottom
	ct3 := triangle{c1, c5, c6}  //up
	ct4 := triangle{c1, c6, c2}  //up
	ct5 := triangle{c2, c6, c7}  //right
	ct6 := triangle{c2, c7, c3}  //right
	ct7 := triangle{c3, c7, c8}  //down
	ct8 := triangle{c3, c8, c4}  //down
	ct9 := triangle{c4, c8, c5}  //left
	ct10 := triangle{c4, c5, c1} //left

	stepY = sizeY / float32(len(heightMap[0])-1)
	stepX = sizeX / float32(len(heightMap)-1)
	heightstep := (stepX + stepY) / (heightFactor * 2)
	var triangles []triangle
	triangles = append(triangles, ct1,ct2,ct3,ct4,ct5,ct6,ct7,ct8,ct9,ct10)
	for i := 0; i< len(heightMap); i++ {
		percentage:= int((float32(i)/float32(len(heightMap)))*100)
		fmt.Printf("100;%d;0\n", percentage)
		for j := 0; j<len(heightMap[0]) ; j++ {
			v1 := Vec3{float32(i) * stepX, float32(j) * stepY, heightMap[i][j]*heightstep + MapHeight}
			if i < len(heightMap)-1 && j < len(heightMap[0])-1 {
				v2 := Vec3{float32(i+1) * stepX, float32(j) * stepY, heightMap[i+1][j]*heightstep + MapHeight}
				v3 := Vec3{float32(i) * stepX, float32(j+1) * stepY, heightMap[i][j+1]*heightstep + MapHeight}
				v4 := Vec3{float32(i+1) * stepX, float32(j+1) * stepY, heightMap[i+1][j+1]*heightstep + MapHeight}
				t1 := triangle{v1, v2, v4}
				t2 := triangle{v1, v4, v3}
				triangles = append(triangles, t1, t2)
			}
			if i == 0 && heightMap[i][j] != 0 {
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if j < len(heightMap[0])-1 {
					vr := Vec3{v1.X(), v1.Y() + stepY, MapHeight}
					triangles = append(triangles, triangle{v1, vd, vr})
				}
				if j > 0 {
					vl := Vec3{v1.X(), v1.Y() - stepY, heightMap[i][j-1]*heightstep + MapHeight}
					triangles = append(triangles, triangle{v1, vl, vd})
				}
			}
			if i == len(heightMap)-1 && heightMap[i][j] != 0 {
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if j < len(heightMap[0])-1 {
					vr := Vec3{v1.X(), v1.Y() + stepY, MapHeight}
					triangles = append(triangles, triangle{v1, vr, vd})
				}
				if j > 0 {
					vl := Vec3{v1.X(), v1.Y() - stepY, heightMap[i][j-1]*heightstep + MapHeight} //falsch
					triangles = append(triangles, triangle{v1, vd, vl})
				}
			}
			if j == 0 && heightMap[i][j] != 0 {
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if i < len(heightMap)-1 {
					vr := Vec3{v1.X() + stepX, v1.Y(), MapHeight}
					triangles = append(triangles, triangle{v1, vd, vr})
				}
				if i > 0 {
					vl := Vec3{v1.X() - stepX, v1.Y(), heightMap[i-1][j]*heightstep + MapHeight}
					triangles = append(triangles, triangle{v1, vl, vd})
				}
			}

			if j == len(heightMap[0])-1 && heightMap[i][j] != 0 {
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if i < len(heightMap)-1 {
					vr := Vec3{v1.X() + stepX, v1.Y(), MapHeight}
					triangles = append(triangles, triangle{v1, vr, vd})
				}
				if i > 0 {
					vl := Vec3{v1.X() - stepX, v1.Y(), heightMap[i-1][j]*heightstep + MapHeight} //falsch
					triangles = append(triangles, triangle{v1, vd, vl})
				}
			}

		}
	}
	generateSTLMapFromTriangles(triangles, fileName)
}

func getSquareMap(heightMap [][]float32) [][]float32 {
	length := int(math.Min(float64(len(heightMap)), float64(len(heightMap[0]))))
	squaredMap := make([][]float32, length)
	for x := 0 ; x < length ; x++ {
		squaredMap[x] = heightMap[x][:length]
	}
	return squaredMap
}


func generateSTLMapFromTriangles(triangles []triangle, fileName string) {
	var header [80]byte
	var byteStl []byte
	byteStl = header[:80]

	triangleCount := convertLittleEndianInt(uint32(len(triangles)))
	byteStl = append(byteStl, triangleCount...)
	for i := 0; i < len(triangles); i++ {
		percentage:= int((float32(i)/float32(len(triangles)))*100)
		fmt.Printf("100;100;%d\n", percentage)

		triangleByte := triangleToByte(triangles[i].v1, triangles[i].v2, triangles[i].v3, fileName)
		byteStl = append(byteStl, triangleByte...)

	}
	fmt.Printf("100;100;100\n")

	writeByteToFile(byteStl, fileName)
}

func triangleToByte(vector1 Vec3, vector2 Vec3, vector3 Vec3, fileName string) []byte {

	var returnBytes []byte
	var colorBytes []byte
	var vectorN Vec3

	dif1 := vector2.Sub(vector1)
	dif2 := vector3.Sub(vector1)

	vectorN = dif1.Cross(dif2)

	colorBytes = []byte{1, 1}

	returnBytes = append(returnBytes, convertLittleEndianFloat(vectorN.X())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vectorN.Y())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vectorN.Z())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vector1.X())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vector1.Y())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vector1.Z())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vector2.X())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vector2.Y())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vector2.Z())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vector3.X())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vector3.Y())...)
	returnBytes = append(returnBytes, convertLittleEndianFloat(vector3.Z())...)
	returnBytes = append(returnBytes, colorBytes...)

	/* returnBytes =append(returnBytes, 	VectorNX[0], VectorNX[1], VectorNX[2], VectorNX[3],
	VectorNY[0], VectorNY[1], VectorNY[2], VectorNY[3],
	VectorNZ[0], VectorNZ[1], VectorNZ[2], VectorNZ[3],
	Vector1X[0], Vector1X[1], Vector1X[2], Vector1X[3],
	Vector1Y[0], Vector1Y[1], Vector1Y[2], Vector1Y[3],
	Vector1Z[0], Vector1Z[1], Vector1Z[2], Vector1Z[3],
	Vector2X[0], Vector2X[1], Vector2X[2], Vector2X[3],
	Vector2Y[0], Vector2Y[1], Vector2Y[2], Vector2Y[3],
	Vector2Z[0], Vector2Z[1], Vector2Z[2], Vector2Z[3],
	Vector3X[0], Vector3X[1], Vector3X[2], Vector3X[3],
	Vector3Y[0], Vector3Y[1], Vector3Y[2], Vector3Y[3],
	Vector3Z[0], Vector3Z[1], Vector3Z[2], Vector3Z[3])

	*/
	return returnBytes
}

func convertLittleEndianFloat(inputEndian float32) []byte {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, inputEndian)
	return buf.Bytes()
}
func convertLittleEndianInt(inputEndian uint32) []byte {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, inputEndian)
	return buf.Bytes()
}

func stlAsci() {
	v1 := r3.Vector{0, 0, 0}
	v2 := r3.Vector{3, 2, 0}
	v3 := r3.Vector{0, 1, 3}

	dif1 := v2.Sub(v1)
	dif2 := v3.Sub(v1)

	n := dif1.Cross(dif2)

	stl := "solid landscape\n" +
		"facet normal %f, %f, %f\n" +
		"outer loop\n" +
		"vertex 0, 0, 0\n" +
		"vertex 3, 0, 1\n" +
		"vertex 0, 3, 2\n" +
		"endloop\n" +
		"endfacet\n" +
		"endsolid landscape"

	stl = fmt.Sprintf(stl, n.X, n.Y, n.Z)

	writeStringToFile(stl)

	log.Print("Test")
}

func writeStringToFile(s string) {
	f, err := os.Create("test.stl")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(s)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func writeByteToFile(b []byte, fileName string) {
	completeFileName := fmt.Sprintf("%s.stl", fileName)
	f, err := os.Create(completeFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.Write(b)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	//fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
