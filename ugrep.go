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
	showLineNum          bool
	showColoredOut       bool
	showMatchedFiles     bool
	showNoMatchFiles     bool
	showMatchedLineCount bool
	fileCount            int
	stdOutWriter         *bufio.Writer
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
	val := "usage: grep [-n] [-c/--colored] [-h/--help] [-l] [-L] [pattern] [file ...]"
	fmt.Fprintf(stdOutWriter, "%s\n", val)
	stdOutWriter.Flush()
}

func printFilename(filename string) {
	if showColoredOut {
		filename = getColoredString(filename, filenameColor)
	}
	fmt.Fprintf(stdOutWriter, "%s\n", filename)
	stdOutWriter.Flush()
}

func printCountOut(filename string, count int) {
	if fileCount > 1 {
		fmt.Fprintf(stdOutWriter, "%s:", filename)
	}
	fmt.Fprintf(stdOutWriter, "%d\n", count)
	stdOutWriter.Flush()
}

func printOut(filename string, matchedLine string, lnum string, matchedIndices [][]int) {
	if showColoredOut {
		filename = getColoredString(filename, filenameColor)
		lnum = getColoredString(lnum, lineNumberColor)
	}
	if fileCount > 1 {
		fmt.Fprintf(stdOutWriter, "%s:", filename)
	}
	if showLineNum {
		fmt.Fprintf(stdOutWriter, "%s: ", lnum)
	}
	lastIdx := 0
	for i := range matchedIndices {
		start := matchedIndices[i][0]
		end := matchedIndices[i][1]
		fmt.Fprintf(stdOutWriter, "%s", matchedLine[lastIdx:start])
		fmt.Fprintf(stdOutWriter, "%s", getColoredString(matchedLine[start:end], patternColor))
		lastIdx = end
	}
	fmt.Fprint(stdOutWriter, matchedLine[lastIdx:], "\n")
	stdOutWriter.Flush()
}

func init() {
	flag.BoolVar(&showLineNum, "n", false, "Flag to specify if you want to print line numbers or not")
	flag.BoolVar(&showColoredOut, "-colored", false, "Flag to specify if you want colored output or not")
	flag.BoolVar(&showMatchedLineCount, "c", false, "Count of selected lines is written to standard output")
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

	re, _ := regexp.Compile(searchTerm)

	for i := range filenames {
		ln := 0
		file, err := os.Open(filenames[i])
		check(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)

		matched := false
		matchedLines := 0
		//Read each line one by one of file
		for scanner.Scan() {
			line := scanner.Text()
			// Check if line contains given search string
			indices := re.FindAllStringIndex(line, -1)
			if showMatchedLineCount {
				if len(indices) > 0 {
					matchedLines++
					continue
				}
			}
			if indices != nil {
				matched = true
				if showMatchedFiles {
					printFilename(filenames[i])
					break
				}
				if showNoMatchFiles {
					break
				}
				printOut(filenames[i], line, strconv.Itoa(ln), indices)
			}
			ln++
		}
		if showNoMatchFiles && !matched {
			printFilename(filenames[i])
		}
		if showMatchedLineCount {
			printCountOut(filenames[i], matchedLines)
		}
	}
}
