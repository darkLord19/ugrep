package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) < 3 {
		panic("Invalid input")
	}
	searchTerm := os.Args[1]
	filenames := os.Args[2:]

	for i := range filenames {
		file, err := os.Open(filenames[i])
		check(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)

		//Read each line one by one of file
		for scanner.Scan() {
			line := scanner.Text()
			// Check if line contains given search string
			if strings.Contains(line, searchTerm) {
				fmt.Printf("%v: %v\n", filenames[i], line)
			}
		}
	}
}
