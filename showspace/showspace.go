package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
)

var trailingSpace = regexp.MustCompile("\\s+$")
var redHighlight = color.New(color.BgRed).SprintfFunc()
var highlightedSpace = redHighlight(" ")

// replace trailing spaces with a red highlight
func replaceTrailingSpaces(s string, r string) string {
	return trailingSpace.ReplaceAllString(s, r)
}

func highlightTrailingSpaces(s string) string {
	return replaceTrailingSpaces(s, highlightedSpace)
}

func main() {
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		fmt.Println(replaceTrailingSpaces(input.Text(), highlightedSpace))
	}
}
