package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html"; charset=utf-8>
	<title>Markdown Preview Tool</title>
	</head>
	<body>`
	footer = `</body></html>`
)

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileName string, out io.Writer) error {
	// Read all the data from the input file and check for errors
	input, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	// Create temporary file and check for errors
	temp, tempErr := os.CreateTemp("", "mdp*.html")
	if tempErr != nil {
		return tempErr
	}

	if tempCloseErr := temp.Close(); tempCloseErr != nil {
		return tempCloseErr
	}

	outputName := temp.Name()
	fmt.Fprintln(out, outputName)

	return saveHTML(outputName, htmlData)
}

func saveHTML(outputName string, htmlData []byte) error {
	// Write the bytes to the file
	return os.WriteFile(outputName, htmlData, fs.FileMode(0644))
}

func parseContent(input []byte) []byte {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var outputHTML bytes.Buffer

	outputHTML.WriteString(header)
	outputHTML.Write(body)
	outputHTML.WriteString(footer)

	return outputHTML.Bytes()
}
	resultFile = "test1.md.html"
