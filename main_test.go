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

func TestFoo(t *testing.T) {
	lines := CaptureReader(strings.NewReader("hello\n"))

	got := strings.Join(lines, "")
	if got != "hello\n" {
		t.Errorf("Want \"hello\", got %s", got)
	}
}
