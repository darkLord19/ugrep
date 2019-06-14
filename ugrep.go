package main

import (
	// "bufio"
	"fmt"
	"strings"
	// "io"
	"io/ioutil"
	"os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main(){
	args := os.Args[1:]
	if args==nil {
		panic("Please provide file name")
	}
	// fmt.Println(args)
	data, err := ioutil.ReadFile(args[0])
    check(err)
	// fmt.Print(string(data))
	if strings.Contains(string(data), args[1]) {
		fmt.Printf("Found %v\n", args[1])
	}
	
}