package main

import "fmt"
import "os"
import "time"

func main() {
	//argsWithProg := os.Args
    //fmt.Println(argsWithProg)
	
	for i := 0; i <= 100; i+=10 {
		fmt.Printf("%d;0;0\n", i)
		time.Sleep(100 * time.Millisecond)
	}
	for i := 0; i <= 100; i+=10 {
		fmt.Printf("100;%d;0\n", i)
		time.Sleep(100 * time.Millisecond)
	}
	for i := 0; i <= 100; i+=10 {
		fmt.Printf("100;100;%d\n", i)
		time.Sleep(100 * time.Millisecond)
	}
	
	os.Exit(0)
}