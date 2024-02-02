package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
	{{ .Body }}
	</body>
</html>`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	templateFileName := flag.String("t", "", "Alternate HTML template name")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, *templateFileName, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileName, templateFileName string, out io.Writer, skipPreview bool) error {
	// Read all the data from the input file and check for errors
	input, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	htmlData, parseErr := parseContent(input, templateFileName)
	if parseErr != nil {
		return err
	}

	// Create temporary file and check for errors
	temp, tempErr := os.CreateTemp("", "mdp*.html")
	if tempErr != nil {
		return tempErr
	}

	if tempCloseErr := temp.Close(); tempCloseErr != nil {
		return tempCloseErr
	}

	outputName := temp.Name()
	defer os.Remove(outputName)

	fmt.Fprintln(out, outputName)

	if saveHTMLErr := saveHTML(outputName, htmlData); saveHTMLErr != nil {
		return saveHTMLErr
	}

	if skipPreview {
		return nil
	}

	return preview(outputName)
}

func saveHTML(outputName string, htmlData []byte) error {
	// Write the bytes to the file
	return os.WriteFile(outputName, htmlData, fs.FileMode(0644))
}

func parseContent(input []byte, templateFileName string) ([]byte, error) {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	if templateFileName != "" {
		t, err = template.ParseFiles(templateFileName)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(body),
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser some time to open the file before deleting it (not an ideal solution)
	time.Sleep(2 * time.Second)
	return err
}
