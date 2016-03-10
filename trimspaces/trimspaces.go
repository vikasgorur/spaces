package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/vikasgorur/spaces"
)

// trimSpaces reads lines from `in` and writes them to `out` with
// trailing spaces removed
func trimSpaces(in io.Reader, out io.Writer) {
	lines := bufio.NewScanner(in)

	for lines.Scan() {
		nonspace := strings.TrimRightFunc(lines.Text(), unicode.IsSpace)
		out.Write([]byte(nonspace + "\n"))
	}
}

// transformFile reads a single file, fixes trailing spaces, and writes it back.
// satisfies filepath.WalkFunc
func transformFile(path string, info os.FileInfo, err error) (bool, error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not fix %v", path)
		return false, nil
	}

	if info.IsDir() && spaces.IsIgnored(path, info) {
		return false, filepath.SkipDir
	}

	if spaces.IsSourceFile(path, info) && !spaces.IsIgnored(path, info) {
		f, err := os.Open(path)
		if err != nil {
			return false, err
		}

		out, err := ioutil.TempFile(".", "ts")
		if err != nil {
			return false, err
		}

		outPath := out.Name()
		trimSpaces(f, out)

		if err := f.Close(); err != nil {
			return false, err
		}
		if err := out.Close(); err != nil {
			return false, err
		}

		if err := os.Rename(outPath, path); err != nil {
			os.Remove(outPath)
			return false, err
		}

		return true, nil
	}
	return false, nil
}

var listFixed = flag.Bool("list-fixed", false, "list files that were fixed in a format suitable for shell arguments")

func main() {
	var dir = flag.Bool("dir", false, "operate recursively on all source files in the current directory.")
	var changed = flag.Bool("changed", false, "operate only on files that have been changed (only works in git repos).")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage: %s [-list-fixed] [-dir|-changed] [file1 ...]

Fix trailing spaces in input files (or stdin).

`, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	files := flag.Args()
	var affected []string

	if *dir {
		affected = spaces.WalkDir(transformFile)
	} else if *changed {
		affected = spaces.WalkChanged(transformFile)
	} else if len(files) != 0 {
		affected = spaces.WalkList(files, transformFile)
	} else {
		trimSpaces(os.Stdin, os.Stdout)
	}

	if *listFixed {
		fmt.Fprintf(os.Stdout, "%v\n", strings.Join(affected, " "))
	}
	os.Exit(0)
}
