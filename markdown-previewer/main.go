package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"runtime"

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
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileName string, out io.Writer, skipPreview bool) error {
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

	if saveHTMLErr := saveHTML(outputName, htmlData); saveHTMLErr != nil {
		return saveHTMLErr
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outputName)

	return preview(outputName)
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

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	// Append filename to parameters slice
	cParams = append(cParams, fname)

	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file using default program
	return exec.Command(cPath, cParams...).Run()
}
