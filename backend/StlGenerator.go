package main

import (
	"fmt"
	"github.com/golang/geo/r3"
	"log"
	"os"
)

func main() {
	v1 := r3.Vector{0, 0, 0}
	v2 := r3.Vector{3, 2, 0}
	v3 := r3.Vector{0, 1, 3}

	dif1 := v2.Sub(v1)
	dif2 := v3.Sub(v1)

	n := dif1.Cross(dif2)

	slt := "solid landscape\n" +
		"facet normal %f, %f, %f\n" +
		"outer loop\n" +
		"vertex 0, 0, 0\n" +
		"vertex 3, 0, 1\n" +
		"vertex 0, 3, 2\n" +
		"endloop\n" +
		"endfacet\n" +
		"endsolid landscape"

	slt = fmt.Sprintf(slt, n.X, n.Y, n.Z)

	writeToFile(slt)

	log.Print("Test")

}

func writeToFile(s string) {
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
