package termscreen

import (
	"bufio"
	"strconv"
	"strings"
	"testing"
)

func TestOneLine(t *testing.T) {
	lines := captureStringReader(strReader("hello\n"))

	got := strings.Join(lines, "")
	assertEqualsStr(t, "hello", got)
}

func TestTwoLines(t *testing.T) {
	lines := captureStringReader(strReader("hello\nworld\n"))

	assertEquals(t, 2, len(lines))
	assertEqualsStr(t, "hello", lines[0])
	assertEqualsStr(t, "world", lines[1])
}

func TestNoNewline(t *testing.T) {
	lines := captureStringReader(strReader("hello"))
	got := strings.Join(lines, "")

	assertEqualsStr(t, "hello", got)
}

func TestPrint(t *testing.T) {
	screen := make([]string, 0)
	lines := print(screen, "hello", 0, 0)

	got := strings.Join(lines, "")
	assertEqualsStr(t, "hello", got)
}

func TestPrintDown(t *testing.T) {
	screen := make([]string, 0)
	lines := print(screen, "hello", 0, 2)

	got := strings.Join(lines, ",")
	assertEqualsStr(t, ",,hello", got)
}

func TestPrintOver(t *testing.T) {
	screen := []string{"hello"}
	lines := print(screen, "world", 0, 0)

	got := strings.Join(lines, "")
	assertEqualsStr(t, "world", got)
}

func TestPrintOverPartly(t *testing.T) {
	screen := []string{"hello"}
	lines := print(screen, "world", 4, 0)

	got := strings.Join(lines, "")
	assertEqualsStr(t, "hellworld", got)
	got = print(lines, "hi, ", 0, 0)[0]
	assertEqualsStr(t, "hi, world", got)
	got = print([]string{"hello world"}, "owdy ", 1, 0)[0]
	assertEqualsStr(t, "howdy world", got)
	got = print([]string{"hello"}, "world", 10, 0)[0]
	assertEqualsStr(t, "hello     world", got)
}

func TestPrintBug(t *testing.T) {
	screen := []string{"\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2"}
	lines := print(screen, ">", 0, 0)

	got := strings.Join(lines, "")
	want := ">\x1b[m * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2"
	assertEqualsStr(t, want, got)
}

func TestDown(t *testing.T) {
	lines := captureStringReader(strReader("hello\x1b[Bhi\n"))

	got := strings.Join(lines, ",")
	assertEqualsStr(t, "hello,     hi", got)
}

func TestUp(t *testing.T) {
	lines := captureStringReader(strReader("hello\n\x1b[Aansi\n"))

	got := strings.Join(lines, "")
	assertEqualsStr(t, "ansio", got)
}

func TestUpDown(t *testing.T) {
	lines := captureStringReader(strReader("one \x1b[2B two \x1b[2A three\n"))

	want := `one       three

     two `
	got := strings.Join(lines, "\n")
	assertEqualsStr(t, want, got)
}

func TestLeftRight(t *testing.T) {
	lines := captureStringReader(strReader("\x1b[10C world \x1b[14D hello,\n"))

	got := strings.Join(lines, ":")
	want := "    hello, world "
	assertEqualsStr(t, want, got)
}

func TestCursorPosition(t *testing.T) {
	for _, code := range []string{"0;0H", ";1H", "1H"} {
		lines := captureStringReader(strReader("\x1b[" + code + "one\n"))

		got := strings.Join(lines, ":")
		want := "one"
		assertEqualsStr(t, want, got)
	}
}

func TestCursorPosition2(t *testing.T) {
	string := ""
	for i := 4; i >= 2; i-- {
		string += "\x1b[" + strconv.Itoa(i) + ";2Ho"
	}
	string += "\n"
	lines := captureStringReader(strReader(string))

	got := strings.Join(lines, "\n")
	want := `
 o
 o
 o`
	assertEqualsStr(t, want, got)
}

func TestCursorPositionAndPrint(t *testing.T) {
	string := "\n o\n o\n o\x1b[3;4Hz\n"
	lines := captureStringReader(strReader(string))

	got := strings.Join(lines, "\n")
	want := `
 o
 o z
 o`
	assertEqualsStr(t, want, got)
}

func TestEraseInLineAll(t *testing.T) {
	for _, str := range []string{"", "Hi \x1b[1K", "Yo \x1b[2K", "\x1b[1K", "\x1b[2K", "\x1b[0K", "\x1b[K"} {
		lines := strings.Join(captureStringReader(strReader(str+"\n")), "\n")

		got := strings.ReplaceAll(lines, " ", "")
		want := ""
		assertEqualsStr(t, want, got)
	}
}

func TestEraseInLine(t *testing.T) {
	str := "Hello, \x1b[1K world!\n"
	lines := strings.Join(captureStringReader(strReader(str)), "\n")

	got := lines
	//nt := "Hello,  world!"
	want := "        world!"
	assertEqualsStr(t, want, got)
}

func TestEraseInLineEnd(t *testing.T) {
	str := "Hello, world! \x1b[1;6H\x1b[K\n"
	lines := strings.Join(captureStringReader(strReader(str)), "\n")

	assertEqualsStr(t, "Hello", lines)
}

func TestEraseInDisplay(t *testing.T) {
	str := "Hello,\n world! \x1b[2J\n"
	lines := strings.Join(captureStringReader(strReader(str)), "\n")

	assertEqualsStr(t, "", lines)
}

func TestEraseInDisplayToEndEmpty(t *testing.T) {
	str := "\x1b[0J\n"
	lines := strings.Join(captureStringReader(strReader(str)), "\n")

	assertEqualsStr(t, "", lines)
}

func TestEraseInDisplayToBeginningEmpty(t *testing.T) {
	str := "\x1b[1J\n"
	lines := strings.Join(captureStringReader(strReader(str)), "\n")

	assertEqualsStr(t, "", lines)
}

func TestEraseInDisplayToEnd(t *testing.T) {
	str := "Howdy, earth\nHello, world \x1b[7D\x1b[A\x1b[0J\n"
	lines := strings.Join(captureStringReader(strReader(str)), "\n")

	assertEqualsStr(t, "Howdy,", lines)
}

func TestEraseInDisplayToBeginning(t *testing.T) {
	str := "Hello,\nworld\x1b[1J\n"
	lines := strings.Join(captureStringReader(strReader(str)), "\n")

	assertEqualsStr(t, "\n     ", lines)
}

func TestLenEmpty(t *testing.T) {
	assertEquals(t, 0, length(""))
}

func TestLenString(t *testing.T) {
	assertEquals(t, 13, length("Hello, world!"))
}

func TestLenColored(t *testing.T) {
	assertEquals(t, 8, length("One \x1b[0m two"))
}

func TestLenUnicode(t *testing.T) {
	assertEquals(t, 1, length("↑"))
}

func TestLenColored2(t *testing.T) {
	assertEquals(t, 8, length("\x1b[31mOne \x1b[0m two"))
}

func TestLenColoredBug(t *testing.T) {
	got := length("\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2")
	assertEquals(t, 43, got)
}

func TestPosZero(t *testing.T) {
	for _, str := range []string{"", "foo"} {
		assertEquals(t, 0, pos(str, 0))
	}
}

func TestPosSimple(t *testing.T) {
	for _, str := range []string{"foo", "foo\x1b[m"} {
		assertEquals(t, 1, pos(str, 1))
		assertEquals(t, 2, pos(str, 2))
		assertEquals(t, 3, pos(str, 3))
	}
}

func TestPos(t *testing.T) {
	str := "\x1b[mABC"
	assertEquals(t, 0, pos(str, 0))
	assertEquals(t, 4, pos(str, 1))
	assertEquals(t, 5, pos(str, 2))
}

func TestPosComplex(t *testing.T) {
	//byte index:           1            2
	//      0   1234567   890123456789   0123
	str := "\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[m\x1b[1;32musability2"
	//col:  0     01234       45678901     123456789012
	//                              1               2
	assertEquals(t, 0, pos(str, 0))
	assertEquals(t, 4, pos(str, 1))
	assertEquals(t, 6, pos(str, 3))
	assertEquals(t, 7, pos(str, 4))
	assertEquals(t, 13, pos(str, 5))
	assertEquals(t, 14, pos(str, 6))
	assertEquals(t, 18, pos(str, 10))
	assertEquals(t, 19, pos(str, 11))
	assertEquals(t, 23, pos(str, 12))
}

func TestPosUnicode(t *testing.T) {
	assertEquals(t, 3, pos("↑ ", 1))
}

func TestPrintStyle(t *testing.T) {
	lines := captureStringReader(strReader("\x1b[31mRED\nHello"))

	got := strings.Join(lines, ":")
	want := "\x1b[31mRED:\x1b[31mHello"
	assertEqualsStr(t, want, got)
}

func TestPrintStyleAccumulate(t *testing.T) {
	lines := captureStringReader(strReader("\x1b[31mRE\x1b[1mD\nHello"))

	got := lines[1]
	want := "\x1b[31m\x1b[1mHello"
	assertEqualsStr(t, want, got)
}

func TestPrintStyleReset(t *testing.T) {
	lines := captureStringReader(strReader("\x1b[31mRED\x1b[0m\nHello"))

	got := strings.Join(lines, ":")
	want := "\x1b[31mRED\x1b[0m:\x1b[0mHello"
	assertEqualsStr(t, want, got)
}

func TestPrintStyleResetOptimize(t *testing.T) {
	lines := captureStringReader(strReader("Foo \x1b[31m\x1b[0m \n bar"))

	assertEqualsStr(t, "\x1b[0m bar", lines[1])
}

func TestPrintStyleBug(t *testing.T) {
	lines := captureStringReader(strReader("\x1b[m  * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[;m\x1b[1;32musability2\n  \x1b[1;1H>"))

	got := lines[0]
	want := "\x1b[;m\x1b[1;32m>\x1b[m * \x1b[33m0793964\x1b[m 2021-04-03 \x1b[33m (\x1b[m\x1b[1;36mHEAD -> \x1b[;m\x1b[1;32musability2"
	assertEqualsStr(t, want, got)
}

func TestUpdateStyle(t *testing.T) {
	assertEqualsStr(t, "", updateStyle([]string{}))
	assertEqualsStr(t, "", updateStyle([]string{""}))
	assertEqualsStr(t, "", updateStyle([]string{"", ""}))
	assertEqualsStr(t, "\x1b[33m", updateStyle([]string{"\x1b[33m"}))
	assertEqualsStr(t, "\x1b[33m", updateStyle([]string{"", "\x1b[33m"}))
	assertEqualsStr(t, "\x1b[m", updateStyle([]string{"\x1b[m"}))
	assertEqualsStr(t, "\x1b[m", updateStyle([]string{"\x1b[33m", "\x1b[m"}))
	assertEqualsStr(t, "\x1b[0m", updateStyle([]string{"\x1b[33m", "\x1b[0m"}))
	assertEquals(t, 0, len(ansiStyleCodes.FindAllString("", -1)))
	assertEquals(t, 1, len(ansiStyleCodes.FindAllString("\x1b[0m", -1)))
	assertEquals(t, 2, len(ansiStyleCodes.FindAllString("\x1b[m\x1b[m", -1)))
}

func TestResetCode(t *testing.T) {
	assertTrue(t, ansiResetCode.MatchString("\x1b[0m"))
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
	assertResetCode(true, "\x1b[;;;m")
	assertResetCode(true, "\x1b[1;31;m")
	assertResetCode(true, "\x1b[1;31;0m")
	assertResetCode(true, "\x1b[1;31;0;m")
	assertResetCode(false, "")
	assertResetCode(false, "foo")
	assertResetCode(false, "[0m")
	assertResetCode(false, "\x1b[1m")
	assertResetCode(false, "\x1b[1;3m")
	assertResetCode(false, "\x1b[0;1m")
	assertResetCode(false, "\x1b[40;1m")
}

func strReader(str string) stringReader {
	return bufio.NewReader(strings.NewReader(str))
}

func assertTrue(t *testing.T, want bool) {
	if !want {
		t.Errorf("Expected true")
	}
}

func assertEquals(t *testing.T, want int, got int) {
	if got != want {
		t.Errorf("Want:\n%d\ngot:\n%d", want, got)
	}
}

func assertEqualsStr(t *testing.T, want string, got string) {
	if got != want {
		t.Errorf("Want:\n%s\ngot:\n%s", want, got)
	}
}
