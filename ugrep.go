package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const (
	patternColor    = "\033[31m"
	filenameColor   = "\033[35m"
	lineNumberColor = "\033[34m"
	resetColor      = "\033[0m"

	EXIT_FAILURE = 1
)

var (
	showLineNum      bool
	showColoredOut   bool
	showMatchedFiles bool
	showNoMatchFiles bool
	fileCount        int
	stdOutWriter     *bufio.Writer
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getColoredString(str string, color string) string {
	return color + str + resetColor
}

func printUsage() {
	val := "usage: grep [-n] [-c/--colored] [-h/--help] [pattern] [file ...]"
	fmt.Fprintf(stdOutWriter, "%s\n", val)
}

func printFilename(filename string) {
	if showColoredOut {
		filename = getColoredString(filename, filenameColor)
	}
	fmt.Fprintf(stdOutWriter, "%s\n", filename)
	stdOutWriter.Flush()
}

func printOut(filename string, matchedLine string, lnum string) {
	if showColoredOut {
		filename = getColoredString(filename, filenameColor)
		lnum = getColoredString(filename, lineNumberColor)
	}
	if fileCount > 1 {
		fmt.Fprintf(stdOutWriter, "%s:", filename)
	}
	if showLineNum {
		fmt.Fprintf(stdOutWriter, "%s: %s\n", lnum, matchedLine)
	} else {
		fmt.Fprintf(stdOutWriter, "%s\n", matchedLine)
	}
	stdOutWriter.Flush()
}

func init() {
	flag.BoolVar(&showLineNum, "n", false, "Flag to specify if you want to print line numbers or not")
	flag.BoolVar(&showColoredOut, "colored", false, "Flag to specify if you want colored output or not")
	flag.BoolVar(&showColoredOut, "c", false, "Flag to specify if you want colored output or not (shorthand)")
	flag.BoolVar(&showMatchedFiles, "l", false, "Flag to get list of files containing search pattern")
	flag.BoolVar(&showNoMatchFiles, "L", false, "Flag to get list of files not containing search pattern")
	flag.Parse()
	stdOutWriter = bufio.NewWriter(os.Stdout)
}

func main() {

	args := flag.Args()

	if len(args) < 2 {
		printUsage()
		os.Exit(EXIT_FAILURE)
	}

	searchTerm := args[0]
	filenames := args[1:]
	fileCount = len(filenames)

	re := regexp.MustCompile(searchTerm)

	for i := range filenames {
		ln := 0
		file, err := os.Open(filenames[i])
		check(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)

		matched := false
		//Read each line one by one of file
		for scanner.Scan() {
			line := scanner.Text()
			// Check if line contains given search string
			indices := re.FindAllStringIndex(line, -1)
			if indices != nil {
				matched = true
				if showMatchedFiles {
					printFilename(filenames[i])
					break
				}
				printOut(filenames[i], line, strconv.Itoa(ln))
			}
			ln++
		}
		if showNoMatchFiles && !matched {
			printFilename(filenames[i])
		}
	}
}
