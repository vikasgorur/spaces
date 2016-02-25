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

	"gopkg.in/fatih/set.v0"
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

var srcExtensions set.Interface = set.New(
	"html",
	"go",
)

// isSourceFile returns true if the path is a source file
func isSourceFile(path string, info os.FileInfo) bool {
	ext := filepath.Ext(path)
	if info.Mode().IsRegular() && strings.HasPrefix(ext, ".") && srcExtensions.Has(filepath.Ext(path)[1:]) {
		return true
	}

	return false
}

// transformFile reads a single file, fixes trailing spaces, and writes it back.
func transformFile(path string, info os.FileInfo, err error) error {
	if isSourceFile(path, info) {
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		out, err := ioutil.TempFile(".", "ts")
		if err != nil {
			return err
		}

		outPath := out.Name()
		trimSpaces(f, out)

		if err := f.Close(); err != nil {
			return err
		}
		if err := out.Close(); err != nil {
			return err
		}

		if err := os.Rename(outPath, path); err != nil {
			os.Remove(outPath)
			return err
		}
	}
	return nil
}

func walkDir() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "trimspaces: could not get current directory: %v\n", err)
		os.Exit(2)
	}

	filepath.Walk(cwd, transformFile)
}

func main() {
	var dir = flag.Bool("dir", true, "operate recursively on all source files in the current directory")
	//var changes = flag.Bool("changes", false, "operate only on files that have been changed")

	if *dir {
		walkDir()
	}
}
