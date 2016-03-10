package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"github.com/vikasgorur/spaces"
)

var redHighlight = color.New(color.BgRed).SprintfFunc()

// Separate the string and its trailing spaces and return them
func trimTrailing(s string) (string, string) {
	nonspace := strings.TrimRightFunc(s, unicode.IsSpace)
	return nonspace, s[len(nonspace):len(s)]
}

// Highlight trailing spaces in the file
func showSpaces(f *os.File) {
	input := bufio.NewScanner(f)

	for input.Scan() {
		nonspace, spaces := trimTrailing(input.Text())
		fmt.Printf("%s%s\n", nonspace, redHighlight(spaces))
	}
}

// Exit nonzero if any line in the file has trailing spaces
func hasSpaces(f *os.File) {
	input := bufio.NewScanner(f)

	for input.Scan() {
		_, spaces := trimTrailing(input.Text())
		if spaces != "" {
			os.Exit(1)
		}
	}
}

var process func(*os.File)

func walk(path string, info os.FileInfo, err error) (bool, error) {
	if spaces.IsSourceFile(path, info) && !spaces.IsIgnored(path, info) {
		f, err := os.Open(path)
		if err != nil {
			return false, err
		}
		defer f.Close()

		process(f)
		return true, nil
	}
	return false, nil
}

func main() {
	var check = flag.Bool("check", false, "check mode; exit 0 if there are no trailing spaces, nonzero otherwise.")
	var dir = flag.Bool("dir", false, "operate recursively on all source files in the current directory.")
	var changed = flag.Bool("changed", false, "operate only on files that have been changed (only works in git repos).")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage: %s [-check] [-dir|-changed] [file1 ...]

Highlight trailing spaces in input files (or stdin).

`, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *check {
		process = hasSpaces
	} else {
		process = showSpaces
	}

	files := flag.Args()
	if *dir {
		spaces.WalkDir(walk)
	} else if *changed {
		spaces.WalkChanged(walk)
	} else if len(files) != 0 {
		spaces.WalkList(files, walk)
	} else {
		process(os.Stdin)
	}

	os.Exit(0)
}
