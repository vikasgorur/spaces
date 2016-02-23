package main

import "testing"

// cases where highlighting should be done
func TestHighlightNeeded(t *testing.T) {
	text := "some text "
	if "some text"+highlightedSpace != highlightTrailingSpaces(text) {
		t.Error("trailing space not replaced")
	}

	text = "  "
	if highlightedSpace+highlightedSpace != highlightTrailingSpaces(text) {
		t.Error("trailing space not replaced")
	}
}

// cases where highlighting should not be done
func TestHighlightNotNeeded(t *testing.T) {
	text := "nospaces"
	if text != highlightTrailingSpaces(text) {
		t.Error("unexpected highlight")
	}

	text = "  foo"
	if text != highlightTrailingSpaces(text) {
		t.Error("unexpected highlight")
	}
}
