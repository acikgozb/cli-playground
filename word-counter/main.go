package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	lineFlag := flag.Bool("l", false, "Count lines of given input")
	byteFlag := flag.Bool("b", false, "Count bytes of given input")
	flag.Parse()

	fmt.Println(count(os.Stdin, *lineFlag, *byteFlag))
}

func count(reader io.Reader, lineFlag bool, byteFlag bool) int {
	scanner := bufio.NewScanner(reader)

	if byteFlag {
		scanner.Split(bufio.ScanBytes)
	} else if !lineFlag {
		scanner.Split(bufio.ScanWords)
	}

	counter := 0

	for scanner.Scan() {
		counter++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return counter
}
