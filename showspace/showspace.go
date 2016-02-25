package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

var redHighlight = color.New(color.BgRed).SprintfFunc()

func trimTrailing(s string) (string, string) {
	nonspace := strings.TrimRightFunc(s, unicode.IsSpace)
	return nonspace, s[len(nonspace):len(s)]
}

func showSpaces(f *os.File) {
	input := bufio.NewScanner(f)

	for input.Scan() {
		nonspace, spaces := trimTrailing(input.Text())
		fmt.Printf("%s%s\n", nonspace, redHighlight(spaces))
	}
}

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		showSpaces(os.Stdin)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "showspace: %v\n", err)
				continue
			}

			showSpaces(f)
			f.Close()
		}
	}
}
