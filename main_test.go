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
	lines := CaptureReader(strings.NewReader("hello\nworld\n"))

	got := len(lines)
	if got != 2 || lines[0] != "hello\n" {
		t.Errorf("Want 2, got %d", got)
	}
	one := lines[0]
	if one != "hello\n" {
		t.Errorf("Want \"hello\", got %s", one)
	}
	two := lines[1]
	if two != "world\n" {
		t.Errorf("Want \"world\", got %s", two)
	}
	if strings.ReplaceAll(strings.Join(lines, " "), "\n", "") != "hello world" {
		t.Errorf("Want \"hello world\", got other")
	}
	all := strings.Join(lines, "")
	if all != "hello\nworld\n" {
		t.Errorf("Want \"world\", got %s", all)
	}
}

func TestPrint(t *testing.T) {
	screen := make([]string, 0)
	lines := Print(screen, "hello", 0, 0)

	got := strings.Join(lines, "")
	if got != "hello" {
		t.Errorf("Want \"hello\", got %s", got)
	}
}

func TestPrintDown(t *testing.T) {
	screen := make([]string, 0)
	lines := Print(screen, "hello", 0, 2)

	got := strings.Join(lines, ",")
	if got != ",,hello" {
		t.Errorf("Want \",,hello\", got %s", got)
	}
}

func TestPrintOver(t *testing.T) {
	screen := make([]string, 0)
	lines := Print(screen, "hello", 0, 0)

	got := strings.Join(lines, "")
	if got != "hello" {
		t.Errorf("Want \"hello\", got %s", got)
	}
}
