package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func trimSpaces(f *os.File) {
	input := bufio.NewScanner(f)

	for input.Scan() {
		nonspace := strings.TrimRightFunc(input.Text(), unicode.IsSpace)
		fmt.Println(nonspace)
	}
}

// FileList is a list of filenames.
type FileList interface {
	// Next() returns a filename, and the second value is true when there
	// are no more files left.
	Next() (string, bool)
}

// Args is a FileList that gets its names from a given slice.
type Args struct {
	args []string
	i    int
}

// NewArgs returns a new Args
func NewArgs(args []string) *Args {
	return &Args{args, 0}
}

// Next returns the next argument
func (a *Args) Next() (string, bool) {
	if a.i == len(a.args) {
		return "", true
	}

	name := a.args[a.i]
	a.i++
	return name, false
}

func main() {
	if len(os.Args[1:]) == 0 {
		trimSpaces(os.Stdin)
	} else {
		files := NewArgs(os.Args[1:])
		arg, end := files.Next()
		for !end {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "trimspaces: %v\n", err)
				continue
			}

			trimSpaces(f)
			f.Close()
			arg, end = files.Next()
		}
	}
}
