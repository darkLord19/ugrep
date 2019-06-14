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
	dat, err := ioutil.ReadFile(args[0])
	check(err)
	data := string(dat)
	// fmt.Print(string(data))
	indice := strings.Index(data, args[1])
	fmt.Println(strings.Split(data[indice:], "\n")[0])
}