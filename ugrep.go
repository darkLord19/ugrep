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
	regex                *regexp.Regexp
	stdOutWriter         *bufio.Writer
)

func check(e error) {
	if e != nil {
		fmt.Fprintf(stdOutWriter, "%s\n", e.Error())
		stdOutWriter.Flush()
		os.Exit(EXIT_FAILURE)
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
	if fileCount > 1 {
		if showColoredOut {
			filename = getColoredString(filename, filenameColor)
		}
		fmt.Fprintf(stdOutWriter, "%s:", filename)
	}
}

func printLineNum(lnum string) {
	if showLineNum {
		if showColoredOut {
			lnum = getColoredString(lnum, lineNumberColor)
		}
		fmt.Fprintf(stdOutWriter, "%s: ", lnum)
	}
}

func printMatches(matchedLine string, matchedIndices [][]int) {
	lastIdx := 0
	for i := range matchedIndices {
		start := matchedIndices[i][0]
		end := matchedIndices[i][1]
		fmt.Fprintf(stdOutWriter, "%s", matchedLine[lastIdx:start])
		fmt.Fprintf(stdOutWriter, "%s", getColoredString(matchedLine[start:end], patternColor))
		lastIdx = end
	}
	fmt.Fprint(stdOutWriter, matchedLine[lastIdx:], "\n")
}

func printCountOut(filename string, count int) {
	printFilename(filename)
	fmt.Fprintf(stdOutWriter, "%d\n", count)
}

func printMatchedFiles(filename string) {
	printFilename(filename)
	fmt.Fprintf(stdOutWriter, "\n")
}

func printOut(filename string, matchedLine string, lnum string, matchedIndices [][]int) {
	printFilename(filename)
	printLineNum(lnum)
	printMatches(matchedLine, matchedIndices)
}

func searchInFile(filename string) {
	ln := 0
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	matched := false
	matchedLines := 0

	//Read each line one by one of file
	for scanner.Scan() {
		line := scanner.Text()
		// Check if line contains given search string
		indices := regex.FindAllStringIndex(line, -1)
		if showMatchedLineCount {
			if len(indices) > 0 {
				matchedLines++
				continue
			}
		}
		if indices != nil {
			matched = true
			if showMatchedFiles {
				printFilename(filename)
				break
			}
			if showNoMatchFiles {
				break
			}
			printOut(filename, line, strconv.Itoa(ln), indices)
		}
		ln++
	}
	if showNoMatchFiles && !matched {
		printMatchedFiles(filename)
	}
	if showMatchedLineCount {
		printCountOut(filename, matchedLines)
	}
}

func init() {
	flag.BoolVar(&showLineNum, "n", false, "Flag to specify if you want to print line numbers or not")
	flag.BoolVar(&showColoredOut, "-colored", false, "Flag to specify if you want colored output or not")
	flag.BoolVar(&showMatchedLineCount, "c", false, "Count of selected lines is written to standard output (shorthand)")
	flag.BoolVar(&showMatchedLineCount, "-count", false, "Count of selected lines is written to standard output")
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

	var err error
	regex, err = regexp.Compile(searchTerm)
	if err != nil {
		os.Exit(EXIT_FAILURE)
	}

	for i := range filenames {
		searchInFile(filenames[i])
	}
	stdOutWriter.Flush()
}
