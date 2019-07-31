package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/geo/r3"
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
const HeightFaktor = 12.0


type triangle struct{
	v1 Vec3
	v2 Vec3
	v3 Vec3
}

func GenerateSTLMapFromSideMap(sidemap []float32, sizeInMM uint32){
	var step float32
	size :=float32(sizeInMM)
	step = size / float32(len(sidemap)-1)
	heightstep := step / HeightFaktor
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
}

func GenerateSTLMapFromHeightMap(heightMap [][]float32, sizeInMM uint32){
	var step float32
	size := float32(sizeInMM)
	c1 := Vec3{0,0,0}//bottom up left
	c2 := Vec3{size,0,0}//bottom up right
	c3 := Vec3{size,size,0} // bottom down right
	c4 := Vec3{0,size,0} //bottom down left
	c5 := Vec3{0,0,MapHeight} // top up left
	c6 := Vec3{size,0,MapHeight} //top up right
	c7 := Vec3{size,size,MapHeight} //top down right
	c8 := Vec3{0,size,MapHeight} //top down left

	ct1 := triangle{c1,c2,c3}//bottom
	ct2 := triangle{c1,c3,c4}//bottom
	ct3 := triangle{c1,c5,c6}//up
	ct4 := triangle{c1,c6,c2}//up
	ct5 := triangle{c2,c6,c7}//right
	ct6 := triangle{c2,c7,c3}//right
	ct7 := triangle{c3,c7,c8}//down
	ct8 := triangle{c3,c8,c4}//down
	ct9 := triangle{c4,c8,c5}//left
	ct10 := triangle{c4,c5,c1}//left

	step = size / float32(len(heightMap)-1)
	heightstep := step / HeightFaktor
	var triangles []triangle
	triangles = append(triangles, ct1,ct2,ct3,ct4,ct5,ct6,ct7,ct8,ct9,ct10)
	for i := 0; i< len(heightMap); i++ {
		for j := 0; j<len(heightMap[0]) ; j++ {
			v1 := Vec3{float32(i) * step, float32(j) * step, heightMap[i][j]*heightstep + MapHeight}
			if (i < len(heightMap)-1 && j < len(heightMap[0])-1) {
				v2 := Vec3{float32(i+1) * step, float32(j) * step, heightMap[i+1][j]*heightstep + MapHeight}
				v3 := Vec3{float32(i) * step, float32(j+1) * step, heightMap[i][j+1]*heightstep + MapHeight}
				v4 := Vec3{float32(i+1) * step, float32(j+1) * step, heightMap[i+1][j+1]*heightstep + MapHeight}
				t1 := triangle{v1, v2, v4}
				t2 := triangle{v1, v4, v3}
				triangles = append(triangles, t1, t2)
			}
			if ( i == 0 && heightMap[i][j]*step != 0){
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if (j < len(heightMap[0])-1){
					vr := Vec3{v1.X(), v1.Y()+step, MapHeight}
					triangles = append(triangles, triangle{v1, vd,vr})
				}
				if (j> 0){
					vl := Vec3{v1.X(), v1.Y()-step, heightMap[i][j-1]*heightstep + MapHeight}
					triangles = append(triangles, triangle{v1, vl,vd})
				}
			}
			if (i == len(heightMap)-1 && heightMap[i][j]*step != 0){
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if (j < len(heightMap[0])-1){ 
					vr := Vec3{v1.X(), v1.Y()+step, MapHeight}
					triangles = append(triangles, triangle{v1, vr,vd})
				}
				if (j> 0){
					vl := Vec3{v1.X(), v1.Y()-step,heightMap[i][j-1]*heightstep + MapHeight}//falsch
					triangles = append(triangles, triangle{v1, vd,vl})
				}
			}
			if (j == 0 && heightMap[i][j]*step != 0){
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if (i < len(heightMap)-1){
					vr := Vec3{v1.X()+step, v1.Y(), MapHeight}
					triangles = append(triangles, triangle{v1, vd,vr})
				}
				if (i> 0){
					vl := Vec3{v1.X()-step, v1.Y(), heightMap[i-1][j]*heightstep + MapHeight}
					triangles = append(triangles, triangle{v1, vl,vd})
				}
			}

			if (j == len(heightMap[0])-1 && heightMap[i][j]*step != 0){
				vd := Vec3{v1.X(), v1.Y(), MapHeight}
				if (i < len(heightMap)-1){
					vr := Vec3{v1.X()+step, v1.Y(), MapHeight}
					triangles = append(triangles, triangle{v1, vr,vd})
				}
				if (i > 0){
					vl := Vec3{v1.X()-step, v1.Y(), heightMap[i-1][j]*heightstep + MapHeight}//falsch
					triangles = append(triangles, triangle{v1, vd,vl})
				}
			}

		}
	}
	generateSTLMapFromTriangles(triangles)
}

func generateSTLMapFromTriangles(triangles []triangle){
	var header [80]byte
	var byteStl []byte
	byteStl = header[:80]

	triangleCount := convertLittleEndianInt(uint32(len(triangles)))
	byteStl = append(byteStl, triangleCount...)
	for i := 0; i< len(triangles) ; i++{

		triangleByte := triangleToByte(triangles[i].v1,triangles[i].v2,triangles[i].v3)
		byteStl = append(byteStl, triangleByte...)

	}

	writeByteToFile(byteStl)

	log.Print("Test")
}

func triangleToByte(vector1 Vec3, vector2 Vec3, vector3 Vec3) []byte {

	var returnBytes []byte
	var colorBytes []byte
	var vectorN Vec3

	dif1 := vector2.Sub(vector1)
	dif2 := vector3.Sub(vector1)

	vectorN = dif1.Cross(dif2)

	colorBytes = []byte{1,1}

	returnBytes =append(returnBytes, convertLittleEndianFloat(vectorN.X())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vectorN.Y())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vectorN.Z())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vector1.X())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vector1.Y())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vector1.Z())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vector2.X())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vector2.Y())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vector2.Z())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vector3.X())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vector3.Y())...)
	returnBytes =append(returnBytes, convertLittleEndianFloat(vector3.Z())...)
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

func writeByteToFile(b []byte) {
	f, err := os.Create("test.stl")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.Write(b)
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
