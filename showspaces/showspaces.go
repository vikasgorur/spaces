package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

var redHighlight = color.New(color.BgRed).SprintfFunc()

// separate the string and its trailing spaces and return them
func trimTrailing(s string) (string, string) {
	nonspace := strings.TrimRightFunc(s, unicode.IsSpace)
	return nonspace, s[len(nonspace):len(s)]
}

// highlight trailing spaces in the file
func showSpaces(f *os.File) {
	input := bufio.NewScanner(f)

	for input.Scan() {
		nonspace, spaces := trimTrailing(input.Text())
		fmt.Printf("%s%s\n", nonspace, redHighlight(spaces))
	}
}

// exit nonzero if any line in the file has trailing spaces
func hasSpaces(f *os.File) {
	input := bufio.NewScanner(f)

	for input.Scan() {
		_, spaces := trimTrailing(input.Text())
		if spaces != "" {
			os.Exit(1)
		}
	}
}

func main() {
	var check = flag.Bool("c", false, "check mode; exit 0 if there are no trailing spaces, nonzero otherwise")

	flag.Parse()

	var process func(*os.File)
	if *check {
		process = hasSpaces
	} else {
		process = showSpaces
	}

	files := flag.Args()
	if len(files) == 0 {
		process(os.Stdin)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "showspace: %v\n", err)
				continue
			}

			process(f)
			f.Close()
		}
	}

	os.Exit(0)
}
