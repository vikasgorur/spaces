package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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
	//TODO handle err properly
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

// walkDir walks every file under the cwd
func walkDir() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "trimspaces: could not get current directory: %v\n", err)
		os.Exit(2)
	}

	filepath.Walk(cwd, transformFile)
}

// extractPath returns a path contained in a line of 'git status' output
// if it's a path we're interested in (added, modified, ...)
// returns nil otherwise
func extractPath(line string) string {
	pieces := strings.Split(line, " ")
	if len(pieces) < 2 {
		return ""
	}

	code := pieces[0]
	if code == "M" || code == "A" || code == "??" {
		//TODO: handle renames
		//TODO: this won't work if the filename contains a space

		return pieces[1]
	}

	return ""
}

// changedFiles returns a slice of file names that have been modified/added
// to the git repository
func changedFiles() ([]string, error) {
	//TODO: make sure this is run only in the top-level dir of the repo
	status := exec.Command("git", "status", "--porcelain")
	output, err := status.Output()
	if err != nil {
		return nil, err
	}

	var paths []string
	lines := bufio.NewScanner(bytes.NewReader(output))
	for lines.Scan() {
		path := extractPath(lines.Text())
		if path != "" {
			paths = append(paths, path)
		}
	}

	return paths, nil
}

// walkChanges walks only changed/added files in the git repository
func walkChanged() {
	paths, err := changedFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
		}
		transformFile(path, info, err)
	}
}

func main() {
	var dir = flag.Bool("dir", true, "operate recursively on all source files in the current directory")
	var changed = flag.Bool("changed", false, "operate only on files that have been changed")

	flag.Parse()

	if *changed {
		walkChanged()
	} else if *dir {
		walkDir()
	}
}
