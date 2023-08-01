// Copyright (C) 2021-2023 Richard H. Tingstad
// This program is free software: you can redistribute it and/or modify it under the terms of the
// GNU General Public License as published by the Free Software Foundation, version 3.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
// without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU General Public License for more details.

package termscreen

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var allAnsiCodes = regexp.MustCompile("\x1b\\[[0-9;]*[A-Za-z]")
var ansiStyleCodes = regexp.MustCompile("\x1b\\[[0-9;]*m")
var ansiResetCode = regexp.MustCompile("\x1b\\[([0-9;]*;[0;]*)?[0;]*m")

type stringReader interface {
	// ReadString reads until the first occurrence of delim in the input,
	// returning a string containing the data up to and including the delimiter.
	// If ReadString encounters an error before finding a delimiter,
	// it returns the data read before the error and the error itself (often io.EOF).
	// ReadString returns err != nil if and only if the returned data does not end
	// in delim.
	ReadString(delim byte) (string, error)
}

func Capture(reader io.Reader) []string {
	var bufioReader *bufio.Reader = bufio.NewReader(reader)
	var myReader stringReader = bufioReader
	lines := captureStringReader(myReader)
	return lines
}

type terminal struct {
	screen []string
	x, y   int
	style  string
}

func captureStringReader(reader stringReader) []string {
	terminal := &terminal{screen: make([]string, 0), x: 0, y: 0, style: ""}
	ansiControlCodes := regexp.MustCompile("\x1b\\[([0-9]*)([ABCDEFGJK]|;?[0-9]*H)")
	for {
		line, err := reader.ReadString('\n')
		if err == nil || (err == io.EOF && len(line) > 0) {
			if len(line) > 0 && line[len(line)-1:] == "\n" {
				line = line[:len(line)-1]
			}
			terminal.handleLine(ansiControlCodes, line)
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

func (terminal *terminal) handleLine(ansiControlCodes *regexp.Regexp, line string) {
	terminal.x = 0
	text := line
	for {
		indices := ansiControlCodes.FindStringSubmatchIndex(text)
		printable := text
		if indices != nil && len(indices) > 4 {
			printable = text[:indices[0]]
		}
		terminal.printTerm(printable)
		if indices != nil && len(indices) > 4 {
			countStart := indices[2]
			countEnd := indices[3]
			codeStart := indices[4]
			codeEnd := indices[5]
			codes := text[codeStart:codeEnd]
			code := codes[:1]
			count := 1
			if countEnd > countStart {
				count = number(text[countStart:countEnd])
			}
			terminal.handleCode(countStart, countEnd, codeStart, codeEnd, count, codes, code)
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

func (terminal *terminal) handleCode(countStart, countEnd, codeStart, codeEnd, count int, codes, code string) {
	screen := terminal.screen
	x, y := terminal.x, terminal.y
	switch code {
	case "A": // Up
		y = max(0, y-count)
	case "B": // Down
		y += count
	case "C": // Forward
		x += count
	case "D": // Back
		x = max(0, x-count)
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
			x = max(0, count-1)
		}
	case "J": // Erase in Display
		idx := pos(screen[y], x)
		if count == 0 || countEnd == countStart { // To end
			if length(screen[y]) > x {
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
		idx := pos(screen[y], x)
		if count == 0 || countEnd == countStart { // To end
			screen[y] = screen[y][0:idx]
		} else if count == 1 { // To beginning
			screen[y] = strings.Repeat(" ", x) + screen[y][idx:]
		} else if count == 2 { // All
			screen[y] = ""
		}
	default:
		if codes[len(codes)-1:] == "H" { // Position
			y = max(0, count-1)
			if codes[0:1] == ";" {
				codes = codes[1:]
			}
			if len(codes) > 1 {
				x = max(0, number(codes[0:len(codes)-1])-1)
			} else {
				x = 0
			}
		}
	}
	terminal.x = x
	terminal.y = y
	terminal.screen = screen
}

func (terminal *terminal) printTerm(text string) {
	screen := terminal.screen
	screen = print(screen, terminal.style+text, terminal.x, terminal.y)
	terminal.screen = screen
	terminal.x += length(text)
	styles := ansiStyleCodes.FindAllString(terminal.style+text, -1)
	terminal.style = updateStyle(styles)
}

func print(screen []string, text string, x int, y int) []string {
	for y >= len(screen) {
		screen = append(screen, "")
	}
	if y < len(screen) {
		prefix := ""
		lineLen := length(screen[y])
		if x < lineLen {
			prefix = screen[y][0:pos(screen[y], x)]
		} else {
			prefix = screen[y] + strings.Repeat(" ", x-lineLen)
		}
		suffix := ""
		if lineLen > x+length(text) {
			idx := pos(screen[y], max(0, min(x+1, lineLen-1)))
			styles := updateStyle(ansiStyleCodes.FindAllString(screen[y][:idx], -1))
			suffix = styles + screen[y][pos(screen[y], x+length(text)):]
		}
		screen[y] = prefix + text + suffix
	}
	return screen
}

func updateStyle(styles []string) string {
	for i := len(styles) - 1; i >= 0; i-- {
		if ansiResetCode.MatchString(styles[i]) {
			styles = styles[i:]
			break
		}
	}
	return strings.Join(styles, "")
}

// pos returns byte index of letter #i (0-based) in string (incl. leading style code):
func pos(value string, i int) int {
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

func length(value string) int {
	stripped := string(allAnsiCodes.ReplaceAll([]byte(value), []byte("")))
	return utf8.RuneCountInString(stripped)
}

func number(value string) int {
	num, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Error converting \"%s\" to int; %s", value, err))
	}
	return num
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func min(nums ...int) int {
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

func use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
