package main

/*func main() {
	v1 := r3.Vector{0, 0, 0}
	v2 := r3.Vector{3, 2, 0}
	v3 := r3.Vector{0, 1, 3}

	dif1 := v2.Sub(v1)
	dif2 := v3.Sub(v1)

	n := dif1.Cross(dif2)

	/*stl := []byte("solid landscape\n" +
	"facet normal %f, %f, %f\n" +
	"outer loop\n" +
	"vertex 0, 0, 0\n" +
	"vertex 3, 0, 1\n" +
	"vertex 0, 3, 2\n" +
	"endloop\n" +
	"endfacet\n" +
	"endsolid landscape")

	var header [80]byte //Header
	var triangleCount = binary.LittleEndian.Uint32()

	stl = fmt.Sprintf(stl, n.X, n.Y, n.Z)

	writeByteToFile(stl)

	log.Print("Test")
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
}*/
