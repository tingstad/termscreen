package main

import (
	"io"
	"strings"
	"testing"
)

type TestReader struct {
	text  string
	count int
}

func (r *TestReader) ReadString(delim byte) (string, error) {
	if r.count < 1 {
		r.count++
		return r.text, nil
	}
	return "", io.EOF
}

func TestOneLine(t *testing.T) {
	lines := CaptureReader(strings.NewReader("hello\n"))

	got := strings.Join(lines, "")
	if got != "hello\n" {
		t.Errorf("Want \"hello\", got %s", got)
	}
}

func TestTwoLines(t *testing.T) {
	lines := CaptureReader(strings.NewReader("hello\n"))

	got := lines[0]
	if got != "hello\n" {
		t.Errorf("Want \"hello\", got %s", got)
	}
}
