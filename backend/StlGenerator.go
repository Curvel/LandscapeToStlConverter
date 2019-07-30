package main

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/geo/r3"
	"fmt"
	//"golang.org/x/image/math/f32"
	//"encoding/binary"
	. "github.com/go-gl/mathgl/mgl32"
	//"github.com/golang/geo/r3"
	"log"
	"os"
)

type triangle struct{
	v1 Vec3
	v2 Vec3
	v3 Vec3
}

func main() {
	v1 := Vec3{0, 0, 0}
	v2 := Vec3{3, 2, 0}
	v3 := Vec3{0, 1, 3}
	v4 := Vec3{3, 2, 1}



	/*stl := []byte("solid landscape\n" +
	"facet normal %f, %f, %f\n" +
	"outer loop\n" +
	"vertex 0, 0, 0\n" +
	"vertex 3, 0, 1\n" +
	"vertex 0, 3, 2\n" +
	"endloop\n" +
	"endfacet\n" +
	"endsolid landscape")*/

	var header [80]byte

	var byteStl []byte
	byteStl = header[:80]
	/*byteStl = []byte{0,0,0,0,0,0,0,0,0,0,
					 0,0,0,0,0,0,0,0,0,0,
					 0,0,0,0,0,0,0,0,0,0,
					 0,0,0,0,0,0,0,0,0,0,
					 0,0,0,0,0,0,0,0,0,0,
					 0,0,0,0,0,0,0,0,0,0,
					 0,0,0,0,0,0,0,0,0,0,
					 0,0,0,0,0,0,0,0,0,0,
					}*/
	triangleCount := convertLittleEndianInt(2)

	triangleByte := triangleToByte(v1,v2,v3)

	byteStl = append(byteStl, triangleCount...)
	byteStl = append(byteStl, triangleByte...)

	triangleByte = triangleToByte(v2,v3,v4)
	byteStl = append (byteStl, triangleByte...)

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
func convertLittleEndianInt(inputEndian int32) []byte {

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
