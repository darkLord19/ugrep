package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"flag"
)

const (
	patternColor    = "\033[31m"
	filenameColor   = "\033[35m"
	lineNumberColor = "\033[34m"
	resetColor      = "\033[0m"
)

var (
	showLineNum bool
	showColoredOut bool
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printOut(filename string, matchedLine string, lnum int) bool {
	if showLineNum {
		fmt.Printf("%s:%d: %s\n", filename, lnum, matchedLine)
	}else{
		fmt.Printf("%s:%d: %s\n", filename, lnum, matchedLine)
	}
	return false
}

func init() {
	flag.BoolVar(&showLineNum, "n", false, "Flag to specify if you want to print line numbers or not")
	flag.BoolVar(&showColoredOut, "colored", false, "Flag to specify if you want colored output or not")
	flag.BoolVar(&showColoredOut, "c", false, "Flag to specify if you want colored output or not (shorthand)")
	flag.Parse()
}

func main() {

	if len(os.Args) < 3 {
		panic("Invalid input")
	}

	args := flag.Args()
	searchTerm := args[0]
	filenames := args[1:]

	for i := range filenames {
		ln := 0
		file, err := os.Open(filenames[i])
		check(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)

		//Read each line one by one of file
		for scanner.Scan() {
			line := scanner.Text()
			// Check if line contains given search string
			if strings.Contains(line, searchTerm) {
				printOut(filenames[i], line, ln)
			}
			ln++
		}
	}
}