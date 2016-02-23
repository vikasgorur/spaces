package main

import (
  "bufio"
  "regexp"
  "fmt"
  "os"
  "github.com/fatih/color"
)

var trailingSpace = regexp.MustCompile("\\s+$")
var redHighlight = color.New(color.BgRed).SprintfFunc()
var highlightedSpace = redHighlight(" ")

// replace trailing spaces with a red highlight
func highlightTrailingSpaces(s string) string {
  return trailingSpace.ReplaceAllString(s, highlightedSpace)    
}

func main() {
  input := bufio.NewScanner(os.Stdin)
  
  for input.Scan() {
    fmt.Println(highlightTrailingSpaces(input.Text()))
  }  
}