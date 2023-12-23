package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	input := bytes.NewBufferString("word1 word2 word3 word4\n")
	expectedOutput := 4

	result := count(input, false, false)

	if result != expectedOutput {
		t.Errorf("Expected %d words, got %d instead.", expectedOutput, result)
	}
}

func TestCountLines(t *testing.T) {
	input := bytes.NewBufferString("line1\nline2 word1\nline3\nline4\nline5")
	expectedOutput := 5

	result := count(input, true, false)

	if result != expectedOutput {
		t.Errorf("Expected %d lines, got %d instead.", expectedOutput, result)
	}
}

func TestCountBytes(t *testing.T) {
	input := bytes.NewBufferString("pls count me as bytes")
	expectedOutput := 21

	result := count(input, false, true)

	if result != expectedOutput {
		t.Errorf("Expected %d bytes, got %d instead.", expectedOutput, result)
	}
}
