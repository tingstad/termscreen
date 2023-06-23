package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var allAnsiCodes = regexp.MustCompile("\x1b\\[[0-9;]*[A-Za-z]")
var ansiStyleCodes = regexp.MustCompile("\x1b\\[[0-9;]*m")
var ansiResetCode = regexp.MustCompile("\x1b\\[([0-9;]*;[0;]*)?[0;]*m")

func main() {
	lines := CaptureReader(os.Stdin)
	for _, line := range lines {
		fmt.Printf("%s\n", line)
	}
}

type MyReader interface {
	ReadString(delim byte) (string, error)
}

func CaptureReader(reader io.Reader) []string {
	var bufioReader *bufio.Reader = bufio.NewReader(reader)
	var myReader MyReader = bufioReader
	lines := Capture(myReader)
	return lines
}

type Terminal struct {
	screen []string
	x, y   int
	style  string
}

func Capture(reader MyReader) []string {
	terminal := &Terminal{screen: make([]string, 0), x: 0, y: 0, style: ""}
	ansiControlCodes := regexp.MustCompile("\x1b\\[([0-9]*)([ABCDEFGJK]|;?[0-9]*H)")
	for {
		line, err := reader.ReadString('\n')
		if err == nil || (err == io.EOF && len(line) > 0) {
			if len(line) > 0 && line[len(line)-1:] == "\n" {
				line = line[:len(line)-1]
			}
			terminal.HandleLine(ansiControlCodes, line)
		}
		if err != nil && err != io.EOF {
			panic(fmt.Sprintf("Error %s", err))
		}
		if err != nil && err == io.EOF {
			break
		}
	}
	return terminal.screen
}

func (terminal *Terminal) HandleLine(ansiControlCodes *regexp.Regexp, line string) {
	terminal.x = 0
	text := line
	for {
		indices := ansiControlCodes.FindStringSubmatchIndex(text)
		printable := text
		if indices != nil && len(indices) > 4 {
			printable = text[:indices[0]]
		}
		terminal.PrintTerm(printable)
		if indices != nil && len(indices) > 4 {
			countStart := indices[2]
			countEnd := indices[3]
			codeStart := indices[4]
			codeEnd := indices[5]
			codes := text[codeStart:codeEnd]
			code := codes[:1]
			count := 1
			if countEnd > countStart {
				count = Number(text[countStart:countEnd])
			}
			terminal.HandleCode(countStart, countEnd, codeStart, codeEnd, count, codes, code)
			if len(text) > indices[1] {
				text = text[indices[1]:]
			} else {
				break
			}
		} else {
			break
		}
	}
	terminal.y += 1
}

func (terminal *Terminal) HandleCode(countStart, countEnd, codeStart, codeEnd, count int, codes, code string) {
	screen := terminal.screen
	x, y := terminal.x, terminal.y
	switch code {
	case "A": // Up
		y = Max(0, y-count)
	case "B": // Down
		y += count
	case "C": // Forward
		x += count
	case "D": // Back
		x = Max(0, x-count)
	case "E": // Next line
		y += count
		x = 0
	case "F": // Previous line
		y -= count
		x = 0
	case "G": // Column
		if countEnd == countStart {
			x = 1
		} else {
			x = Max(0, count-1)
		}
	case "J": // Erase in Display
		idx := Pos(screen[y], x)
		if count == 0 || countEnd == countStart { // To end
			if Len(screen[y]) > x {
				screen[y] = screen[y][0:idx]
			}
			screen = screen[0 : y+1]
		} else if count == 1 { // To begining
			screen[y] = strings.Repeat(" ", x) + screen[y][idx:]
			for idx := range screen[0:y] {
				screen[idx] = ""
			}
		} else if count > 1 { // All
			screen = screen[:0]
			x = 0
			y = 0
		}
	case "K": // Erase in Line
		idx := Pos(screen[y], x)
		if count == 0 || countEnd == countStart { // To end
			screen[y] = screen[y][0:idx]
		} else if count == 1 { // To beginning
			screen[y] = strings.Repeat(" ", x) + screen[y][idx:]
		} else if count == 2 { // All
			screen[y] = ""
		}
	default:
		if codes[len(codes)-1:] == "H" { // Position
			y = Max(0, count-1)
			if codes[0:1] == ";" {
				codes = codes[1:]
			}
			if len(codes) > 1 {
				x = Max(0, Number(codes[0:len(codes)-1])-1)
			} else {
				x = 0
			}
		}
	}
	terminal.x = x
	terminal.y = y
	terminal.screen = screen
}

func (terminal *Terminal) PrintTerm(text string) {
	screen := terminal.screen
	screen = Print(screen, terminal.style+text, terminal.x, terminal.y)
	terminal.screen = screen
	terminal.x += Len(text)
	styles := ansiStyleCodes.FindAllString(terminal.style+text, -1)
	terminal.style = UpdateStyle(styles)
}

func Print(screen []string, text string, x int, y int) []string {
	for y >= len(screen) {
		screen = append(screen, "")
	}
	if y < len(screen) {
		prefix := ""
		lineLen := Len(screen[y])
		if x < lineLen {
			prefix = screen[y][0:Pos(screen[y], x)]
		} else {
			prefix = screen[y] + strings.Repeat(" ", x-lineLen)
		}
		suffix := ""
		if lineLen > x+Len(text) {
			idx := Pos(screen[y], Max(0, Min(x+1, lineLen-1)))
			styles := UpdateStyle(ansiStyleCodes.FindAllString(screen[y][:idx], -1))
			suffix = styles + screen[y][Pos(screen[y], x+Len(text)):]
		}
		screen[y] = prefix + text + suffix
	}
	return screen
}

func UpdateStyle(styles []string) string {
	for i := len(styles) - 1; i >= 0; i-- {
		if ansiResetCode.MatchString(styles[i]) {
			styles = styles[i:]
			break
		}
	}
	return strings.Join(styles, "")
}

// Returns byte index of letter #i (0-based) in string (incl. leading style code):
func Pos(value string, i int) int {
	if len(value) == 0 {
		return 0
	}
	offset := 0
	columns := 0
	for {
		pos := allAnsiCodes.FindStringIndex(value)
		passed := value
		if pos != nil {
			passed = value[0:pos[0]]
		}
		lenPassed := len(passed)
		for index, w := 0, 0; index < lenPassed; index += w {
			_, w = utf8.DecodeRuneInString(passed[index:])
			if columns >= i {
				return offset + index
			}
			columns++
		}
		offset += lenPassed
		if columns >= i {
			return offset
		}
		if pos != nil {
			value = value[pos[1]:]
			offset += pos[1] - pos[0]
		}
	}
}

func Len(value string) int {
	stripped := string(allAnsiCodes.ReplaceAll([]byte(value), []byte("")))
	return utf8.RuneCountInString(stripped)
}

func Number(value string) int {
	num, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Error converting \"%s\" to int; %s", value, err))
	}
	return num
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(nums ...int) int {
	if len(nums) < 1 {
		panic("No args!")
	}
	smallest := nums[0]
	for _, num := range nums {
		if num < smallest {
			smallest = num
		}
	}
	return smallest
}

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
