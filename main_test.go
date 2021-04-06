package main

import (
	"strconv"
	"strings"
	"testing"
)

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

func TestNoNewline(t *testing.T) {
	lines := CaptureReader(strings.NewReader("hello"))

	got := strings.Join(lines, "")
	if got != "hello" {
		t.Errorf("Want \"hello\", got \"%s\"", got)
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

func TestPrintBug(t *testing.T) {
	screen := []string{"hello"}
	lines := Print(screen, "world", 0, 0)

	got := strings.Join(lines, "")
	if got != "world" {
		t.Errorf("Want \"world\", got %s", got)
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

func TestCursorPositionAndPrint(t *testing.T) {
	string := "\n o\n o\n o\x1b[3;4Hz\n"
	lines := CaptureReader(strings.NewReader(string))

	got := strings.Join(lines, "\n")
	want := `
 o
 o z
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
	str := "Hello,\nworld\x1b[1J\n"
	lines := strings.Join(CaptureReader(strings.NewReader(str)), "\n")

	got := lines
	want := "\n     "
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}

func TestLenEmpty(t *testing.T) {
	got := Len("")
	want := 0
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func TestLenString(t *testing.T) {
	got := Len("Hello, world!")
	want := 13
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func TestLenColored(t *testing.T) {
	got := Len("One \x1b[0m two")
	want := 8
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func TestLenUnicode(t *testing.T) {
	got := Len("↑")
	want := 1
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func TestPosZero(t *testing.T) {
	for _, str := range []string{"", "foo"} {
		got := Pos(str, 0)
		want := 0
		if got != want {
			t.Errorf("Want:\n%d\ngot:\n%d", want, got)
		}
	}
}

func TestPosSimple(t *testing.T) {
	for _, str := range []string{"foo", "foo\x1b[m"} {
		got := Pos(str, 1)
		want := 1
		if got != want {
			t.Errorf("Want:\n%d\ngot:\n%d", want, got)
		}
	}
}

func TestPos(t *testing.T) {
	got := Pos("\x1b[mabc", 1)
	want := 4
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func TestPosZeroAnsi(t *testing.T) {
	got := Pos("\x1b[mFoo", 0)
	want := 3
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func TestPosBug(t *testing.T) {
	got := Pos("\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2", 0)
	want := 3
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func TestPosUnicode(t *testing.T) {
	got := Pos("↑ ", 1)
	want := 3
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func TestPrintStyle(t *testing.T) {
	lines := CaptureReader(strings.NewReader("\x1b[31mRED\nHello"))

	got := strings.Join(lines, ":")
	want := "\x1b[31mRED:\x1b[31mHello"
	if got != want {
		t.Errorf("Want \"%s\", got \"%s\"", want, got)
	}
}

func TestPrintStyleAccumulate(t *testing.T) {
	lines := CaptureReader(strings.NewReader("\x1b[31mRE\x1b[1mD\nHello"))

	got := lines[1]
	want := "\x1b[31m\x1b[1mHello"
	if got != want {
		t.Errorf("Want \"%s\", got \"%s\"", want, got)
	}
}

func TestPrintStyleReset(t *testing.T) {
	lines := CaptureReader(strings.NewReader("\x1b[31mRED\x1b[0m\nHello"))

	got := strings.Join(lines, ":")
	want := "\x1b[31mRED\x1b[0m:\x1b[0mHello"
	if got != want {
		t.Errorf("Want \"%s\", got \"%s\"", want, got)
	}
}

func TestPrintStyleResetOptimize(t *testing.T) {
	lines := CaptureReader(strings.NewReader("Foo \x1b[31m\x1b[0m \n bar"))

	got := lines[1]
	want := "\x1b[0m bar"
	if got != want {
		t.Errorf("Want \"%s\", got \"%s\"", want, got)
	}
}

func FixTestPrintStyleBug(t *testing.T) {
	lines := CaptureReader(strings.NewReader("\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2\n  \x1b[1;1H>"))

	got := lines[0]
	want := "\x1b[m>  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2"
	if got != want {
		t.Errorf("Want \"%s\", got \"%s\"", want, got)
	}
}
