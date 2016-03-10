package spaces

import "testing"

func TestExtractPath(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{" M test.go", "test.go"},
		{"?? index.html", "index.html"},
		{"?? name with spaces", "name with spaces"},
		{"RM README.md -> file name.md", "file name.md"},
	}

	for _, c := range cases {
		actual := extractPath(c.input)
		if c.output != extractPath(c.input) {
			t.Errorf("extractPath failed, expected: %v, actual: %v\n", c.output, actual)
		}
	}
}
