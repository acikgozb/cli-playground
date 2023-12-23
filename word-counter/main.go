package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Define flags for CLI
	lineFlag := flag.Bool("l", false, "Count lines of given input")
	byteFlag := flag.Bool("b", false, "Count bytes of given input")
	flag.Parse()

	// Calling the count() function to count the number of words
	// received from the Standard Input and printing it out
	fmt.Println(count(os.Stdin, *lineFlag, *byteFlag))
}

func count(reader io.Reader, lineFlag bool, byteFlag bool) int {
	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(reader)

	// Define the scanner to split type to words (default is split by lines)
	if byteFlag {
		scanner.Split(bufio.ScanBytes)
	} else if !lineFlag {
		scanner.Split(bufio.ScanWords)
	}

	// Defining a counter
	counter := 0

	for scanner.Scan() {
		counter++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return counter
}
