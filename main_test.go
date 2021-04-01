package main

import (
	"io"
	"strconv"
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
	if got != "hello" {
		t.Errorf("Want \"hello\", got \"%s\"", got)
	}
}

func TestTwoLines(t *testing.T) {
	lines := CaptureReader(strings.NewReader("hello\nworld\n"))

	got := len(lines)
	if got != 2 {
		t.Errorf("Want 2, got %d", got)
	}
	one := lines[0]
	if one != "hello" {
		t.Errorf("Want \"hello\", got %s", one)
	}
	two := lines[1]
	if two != "world" {
		t.Errorf("Want \"world\", got %s", two)
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
	screen := []string{"hello"}
	lines := Print(screen, "world", 0, 0)

	got := strings.Join(lines, "")
	if got != "world" {
		t.Errorf("Want \"world\", got %s", got)
	}
}

func TestPrintOverPartly(t *testing.T) {
	screen := []string{"hello"}
	lines := Print(screen, "world", 4, 0)

	got := strings.Join(lines, "")
	if got != "hellworld" {
		t.Errorf("Want \"hellworld\", got %s", got)
	}
	got = Print(lines, "hi, ", 0, 0)[0]
	if got != "hi, world" {
		t.Errorf("Want \"hi, world\", got %s", got)
	}
	got = Print([]string{"hello world"}, "owdy ", 1, 0)[0]
	if got != "howdy world" {
		t.Errorf("Want \"howdy world\", got %s", got)
	}
	got = Print([]string{"hello"}, "world", 10, 0)[0]
	if got != "hello     world" {
		t.Errorf("Want \"hello     world\", got %s", got)
	}
}

func TestDown(t *testing.T) {
	lines := CaptureReader(strings.NewReader("hello\x1b[Bhi\n"))

	got := strings.Join(lines, ",")
	if got != "hello,     hi" {
		t.Errorf("Want \"hello,     hi\", got \"%s\"", got)
	}
}

func TestUp(t *testing.T) {
	lines := CaptureReader(strings.NewReader("hello\n\x1b[Aansi\n"))

	got := strings.Join(lines, "")
	if got != "ansio" {
		t.Errorf("Want \"ansio\", got \"%s\"", got)
	}
}

func TestUpDown(t *testing.T) {
	lines := CaptureReader(strings.NewReader("one \x1b[2B two \x1b[2A three\n"))

	want := `one       three

     two `
	got := strings.Join(lines, "\n")
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestLeftRight(t *testing.T) {
	lines := CaptureReader(strings.NewReader("\x1b[10C world \x1b[14D hello,\n"))

	got := strings.Join(lines, ":")
	want := "    hello, world "
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestCursorPosition(t *testing.T) {
	for _, code := range []string{"0;0H", ";1H", "1H"} {
		lines := CaptureReader(strings.NewReader("\x1b[" + code + "one\n"))

		got := strings.Join(lines, ":")
		want := "one"
		if got != want {
			t.Errorf("Want:\n%s\ngot:\n%s", want, got)
		}
	}
}

func TestCursorPosition2(t *testing.T) {
	string := ""
	for i := 4; i >= 2; i-- {
		string += "\x1b[" + strconv.Itoa(i) + ";2Ho"
	}
	string += "\n"
	lines := CaptureReader(strings.NewReader(string))

	got := strings.Join(lines, "\n")
	want := `
 o
 o
 o`
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestEraseInLineAll(t *testing.T) {
	for _, str := range []string{"", "Hi \x1b[1K", "Yo \x1b[2K", "\x1b[1K", "\x1b[2K", "\x1b[0K", "\x1b[K"} {
		lines := strings.Join(CaptureReader(strings.NewReader(str+"\n")), "\n")

		got := strings.ReplaceAll(lines, " ", "")
		want := ""
		if got != want {
			t.Errorf("Want:\n%s\ngot:\n%s", want, got)
		}
	}
}

func TestEraseInLine(t *testing.T) {
	str := "Hello, \x1b[1K world!\n"
	lines := strings.Join(CaptureReader(strings.NewReader(str)), "\n")

	got := lines
	//nt := "Hello,  world!"
	want := "        world!"
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestEraseInLineEnd(t *testing.T) {
	str := "Hello, world! \x1b[1;6H\x1b[K\n"
	lines := strings.Join(CaptureReader(strings.NewReader(str)), "\n")

	got := lines
	want := "Hello"
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestEraseInDisplay(t *testing.T) {
	str := "Hello,\n world! \x1b[2J\n"
	lines := strings.Join(CaptureReader(strings.NewReader(str)), "\n")

	got := lines
	want := ""
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestEraseInDisplayToEndEmpty(t *testing.T) {
	str := "\x1b[0J\n"
	lines := strings.Join(CaptureReader(strings.NewReader(str)), "\n")

	got := lines
	want := ""
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestEraseInDisplayToBeginningEmpty(t *testing.T) {
	str := "\x1b[1J\n"
	lines := strings.Join(CaptureReader(strings.NewReader(str)), "\n")

	got := lines
	want := ""
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestEraseInDisplayToEnd(t *testing.T) {
	str := "Howdy, earth\nHello, world \x1b[7D\x1b[A\x1b[0J\n"
	lines := strings.Join(CaptureReader(strings.NewReader(str)), "\n")

	got := lines
	want := "Howdy,"
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestEraseInDisplayToBeginning(t *testing.T) {
	str := "Howdy, earth\nHello, world \x1b[7D\x1b[A\x1b[0J\n"
	lines := strings.Join(CaptureReader(strings.NewReader(str)), "\n")

	got := lines
	want := "Howdy,"
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}
