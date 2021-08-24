package main

import (
	"bufio"
	"strconv"
	"strings"
	"testing"
)

func TestOneLine(t *testing.T) {
	lines := Capture(StrReader("hello\n"))

	got := strings.Join(lines, "")
	AssertEqualsStr(t, "hello", got)
}

func TestTwoLines(t *testing.T) {
	lines := Capture(StrReader("hello\nworld\n"))

	AssertEquals(t, 2, len(lines))
	AssertEqualsStr(t, "hello", lines[0])
	AssertEqualsStr(t, "world", lines[1])
}

func TestNoNewline(t *testing.T) {
	lines := Capture(StrReader("hello"))
	got := strings.Join(lines, "")

	AssertEqualsStr(t, "hello", got)
}

func TestPrint(t *testing.T) {
	screen := make([]string, 0)
	lines := Print(screen, "hello", 0, 0)

	got := strings.Join(lines, "")
	AssertEqualsStr(t, "hello", got)
}

func TestPrintDown(t *testing.T) {
	screen := make([]string, 0)
	lines := Print(screen, "hello", 0, 2)

	got := strings.Join(lines, ",")
	AssertEqualsStr(t, ",,hello", got)
}

func TestPrintOver(t *testing.T) {
	screen := []string{"hello"}
	lines := Print(screen, "world", 0, 0)

	got := strings.Join(lines, "")
	AssertEqualsStr(t, "world", got)
}

func TestPrintOverPartly(t *testing.T) {
	screen := []string{"hello"}
	lines := Print(screen, "world", 4, 0)

	got := strings.Join(lines, "")
	AssertEqualsStr(t, "hellworld", got)
	got = Print(lines, "hi, ", 0, 0)[0]
	AssertEqualsStr(t, "hi, world", got)
	got = Print([]string{"hello world"}, "owdy ", 1, 0)[0]
	AssertEqualsStr(t, "howdy world", got)
	got = Print([]string{"hello"}, "world", 10, 0)[0]
	AssertEqualsStr(t, "hello     world", got)
}

func FixTestPrintBug(t *testing.T) {
	screen := []string{"\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2"}
	lines := Print(screen, ">", 0, 0)

	got := strings.Join(lines, "")
	want := ">\x1b[m * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2"
	AssertEqualsStr(t, want, got)
}

func TestDown(t *testing.T) {
	lines := Capture(StrReader("hello\x1b[Bhi\n"))

	got := strings.Join(lines, ",")
	AssertEqualsStr(t, "hello,     hi", got)
}

func TestUp(t *testing.T) {
	lines := Capture(StrReader("hello\n\x1b[Aansi\n"))

	got := strings.Join(lines, "")
	AssertEqualsStr(t, "ansio", got)
}

func TestUpDown(t *testing.T) {
	lines := Capture(StrReader("one \x1b[2B two \x1b[2A three\n"))

	want := `one       three

     two `
	got := strings.Join(lines, "\n")
	AssertEqualsStr(t, want, got)
}

func TestLeftRight(t *testing.T) {
	lines := Capture(StrReader("\x1b[10C world \x1b[14D hello,\n"))

	got := strings.Join(lines, ":")
	want := "    hello, world "
	AssertEqualsStr(t, want, got)
}

func TestCursorPosition(t *testing.T) {
	for _, code := range []string{"0;0H", ";1H", "1H"} {
		lines := Capture(StrReader("\x1b[" + code + "one\n"))

		got := strings.Join(lines, ":")
		want := "one"
		AssertEqualsStr(t, want, got)
	}
}

func TestCursorPosition2(t *testing.T) {
	string := ""
	for i := 4; i >= 2; i-- {
		string += "\x1b[" + strconv.Itoa(i) + ";2Ho"
	}
	string += "\n"
	lines := Capture(StrReader(string))

	got := strings.Join(lines, "\n")
	want := `
 o
 o
 o`
	AssertEqualsStr(t, want, got)
}

func TestCursorPositionAndPrint(t *testing.T) {
	string := "\n o\n o\n o\x1b[3;4Hz\n"
	lines := Capture(StrReader(string))

	got := strings.Join(lines, "\n")
	want := `
 o
 o z
 o`
	AssertEqualsStr(t, want, got)
}

func TestEraseInLineAll(t *testing.T) {
	for _, str := range []string{"", "Hi \x1b[1K", "Yo \x1b[2K", "\x1b[1K", "\x1b[2K", "\x1b[0K", "\x1b[K"} {
		lines := strings.Join(Capture(StrReader(str+"\n")), "\n")

		got := strings.ReplaceAll(lines, " ", "")
		want := ""
		AssertEqualsStr(t, want, got)
	}
}

func TestEraseInLine(t *testing.T) {
	str := "Hello, \x1b[1K world!\n"
	lines := strings.Join(Capture(StrReader(str)), "\n")

	got := lines
	//nt := "Hello,  world!"
	want := "        world!"
	AssertEqualsStr(t, want, got)
}

func TestEraseInLineEnd(t *testing.T) {
	str := "Hello, world! \x1b[1;6H\x1b[K\n"
	lines := strings.Join(Capture(StrReader(str)), "\n")

	AssertEqualsStr(t, "Hello", lines)
}

func TestEraseInDisplay(t *testing.T) {
	str := "Hello,\n world! \x1b[2J\n"
	lines := strings.Join(Capture(StrReader(str)), "\n")

	AssertEqualsStr(t, "", lines)
}

func TestEraseInDisplayToEndEmpty(t *testing.T) {
	str := "\x1b[0J\n"
	lines := strings.Join(Capture(StrReader(str)), "\n")

	AssertEqualsStr(t, "", lines)
}

func TestEraseInDisplayToBeginningEmpty(t *testing.T) {
	str := "\x1b[1J\n"
	lines := strings.Join(Capture(StrReader(str)), "\n")

	AssertEqualsStr(t, "", lines)
}

func TestEraseInDisplayToEnd(t *testing.T) {
	str := "Howdy, earth\nHello, world \x1b[7D\x1b[A\x1b[0J\n"
	lines := strings.Join(Capture(StrReader(str)), "\n")

	AssertEqualsStr(t, "Howdy,", lines)
}

func TestEraseInDisplayToBeginning(t *testing.T) {
	str := "Hello,\nworld\x1b[1J\n"
	lines := strings.Join(Capture(StrReader(str)), "\n")

	AssertEqualsStr(t, "\n     ", lines)
}

func TestLenEmpty(t *testing.T) {
	AssertEquals(t, 0, Len(""))
}

func TestLenString(t *testing.T) {
	AssertEquals(t, 13, Len("Hello, world!"))
}

func TestLenColored(t *testing.T) {
	AssertEquals(t, 8, Len("One \x1b[0m two"))
}

func TestLenUnicode(t *testing.T) {
	AssertEquals(t, 1, Len("↑"))
}

func TestLenColored2(t *testing.T) {
	AssertEquals(t, 8, Len("\x1b[31mOne \x1b[0m two"))
}

func TestLenColoredBug(t *testing.T) {
	got := Len("\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2")
	AssertEquals(t, 43, got)
}

func TestPosZero(t *testing.T) {
	for _, str := range []string{"", "foo"} {
		AssertEquals(t, 0, Pos(str, 0))
	}
}

func TestPosSimple(t *testing.T) {
	for _, str := range []string{"foo", "foo\x1b[m"} {
		AssertEquals(t, 1, Pos(str, 1))
		AssertEquals(t, 2, Pos(str, 2))
		AssertEquals(t, 3, Pos(str, 3))
	}
}

func TestPos(t *testing.T) {
	str := "\x1b[mABC"
	AssertEquals(t, 0, Pos(str, 0))
	AssertEquals(t, 4, Pos(str, 1))
	AssertEquals(t, 5, Pos(str, 2))
}

func TestPosComplex(t *testing.T) {
	//byte index:           1            2
	//      0   1234567   890123456789   0123
	str := "\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2"
	//col:  0     01234       45678901     123456789012
	//                              1               2
	AssertEquals(t, 0, Pos(str, 0))
	AssertEquals(t, 4, Pos(str, 1))
	AssertEquals(t, 6, Pos(str, 3))
	AssertEquals(t, 7, Pos(str, 4))
	AssertEquals(t, 13, Pos(str, 5))
	AssertEquals(t, 14, Pos(str, 6))
	AssertEquals(t, 18, Pos(str, 10))
	AssertEquals(t, 19, Pos(str, 11))
	AssertEquals(t, 23, Pos(str, 12))
}

func TestPosUnicode(t *testing.T) {
	AssertEquals(t, 3, Pos("↑ ", 1))
}

func TestPrintStyle(t *testing.T) {
	lines := Capture(StrReader("\x1b[31mRED\nHello"))

	got := strings.Join(lines, ":")
	want := "\x1b[31mRED:\x1b[31mHello"
	AssertEqualsStr(t, want, got)
}

func TestPrintStyleAccumulate(t *testing.T) {
	lines := Capture(StrReader("\x1b[31mRE\x1b[1mD\nHello"))

	got := lines[1]
	want := "\x1b[31m\x1b[1mHello"
	AssertEqualsStr(t, want, got)
}

func TestPrintStyleReset(t *testing.T) {
	lines := Capture(StrReader("\x1b[31mRED\x1b[0m\nHello"))

	got := strings.Join(lines, ":")
	want := "\x1b[31mRED\x1b[0m:\x1b[0mHello"
	AssertEqualsStr(t, want, got)
}

func TestPrintStyleResetOptimize(t *testing.T) {
	lines := Capture(StrReader("Foo \x1b[31m\x1b[0m \n bar"))

	AssertEqualsStr(t, "\x1b[0m bar", lines[1])
}

func FixTestPrintStyleBug(t *testing.T) {
	lines := Capture(StrReader("\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2\n  \x1b[1;1H>"))

	got := lines[0]
	want := "\x1b[m>  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2"
	AssertEqualsStr(t, want, got)
}

func TestUpdateStyle(t *testing.T) {
	AssertEqualsStr(t, "", UpdateStyle([]string{}))
	AssertEqualsStr(t, "", UpdateStyle([]string{""}))
	AssertEqualsStr(t, "", UpdateStyle([]string{"", ""}))
	AssertEqualsStr(t, "\x1b[33m", UpdateStyle([]string{"\x1b[33m"}))
	AssertEqualsStr(t, "\x1b[33m", UpdateStyle([]string{"", "\x1b[33m"}))
	AssertEqualsStr(t, "\x1b[m", UpdateStyle([]string{"\x1b[m"}))
	AssertEqualsStr(t, "\x1b[m", UpdateStyle([]string{"\x1b[33m", "\x1b[m"}))
	AssertEqualsStr(t, "\x1b[0m", UpdateStyle([]string{"\x1b[33m", "\x1b[0m"}))
	AssertEquals(t, 0, len(ansiStyleCodes.FindAllString("", -1)))
	AssertEquals(t, 1, len(ansiStyleCodes.FindAllString("\x1b[0m", -1)))
	AssertEquals(t, 2, len(ansiStyleCodes.FindAllString("\x1b[m\x1b[m", -1)))
}

func TestResetCode(t *testing.T) {
	AssertTrue(t, ansiResetCode.MatchString("\x1b[0m"))
	var assertResetCode = func(expect bool, text string) {
		if ansiResetCode.MatchString(text) != expect {
			t.Errorf("Expected '%s' reset code match to be %t", text, expect)
		}
	}
	assertResetCode(true, "\x1b[m")
	assertResetCode(true, "\x1b[0m")
	assertResetCode(true, "\x1b[;m")
	assertResetCode(true, "\x1b[1;m")
	assertResetCode(true, "\x1b[;0;m")
	assertResetCode(true, "\x1b[1;31;0;m")
	assertResetCode(false, "")
	assertResetCode(false, "foo")
	assertResetCode(false, "[0m")
	assertResetCode(false, "\x1b[1m")
	assertResetCode(false, "\x1b[1;3m")
}

func StrReader(str string) MyReader {
	return bufio.NewReader(strings.NewReader(str))
}

func AssertTrue(t *testing.T, want bool) {
	if !want {
		t.Errorf("Expected true")
	}
}

func AssertEquals(t *testing.T, want int, got int) {
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func AssertEqualsStr(t *testing.T, want string, got string) {
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}
