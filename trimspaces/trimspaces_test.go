package main

import "testing"

func TestExtractPath(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{" M test.go", "test.go"},
		{"?? index.html", "index.html"},
		{"RM index.html -> index.htmls", ""},
	}

	for _, c := range cases {
		if c.output != extractPath(c.input) {
			t.Errorf("extractPath failed, expected: %v, actual: %v\n", c.output, c.input)
		}
	}
}
