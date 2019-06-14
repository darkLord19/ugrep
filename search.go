package main

import (
	// "bufio"
	"fmt"
	// "strings"
	// "io"
	"io/ioutil"
	"os"
	"regexp"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	args := os.Args[1:]
	if args == nil {
		panic("Please provide file name")
	}
	// fmt.Println(args)
	dat, err := ioutil.ReadFile(args[0])
	check(err)
	data := string(dat)
	// fmt.Print(string(data))
	re := regexp.MustCompile(args[1])
	idx := re.FindAllStringIndex(data, -1)
	for i := range idx {
		fmt.Printf("%v: %v\n", args[0], data[idx[i][0]:idx[i][1]])
	}
}
