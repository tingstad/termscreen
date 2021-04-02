package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

func Capture(reader MyReader) []string {
	screen := make([]string, 0)
	esc := "\x1b"
	re := regexp.MustCompile(esc + "\\[([0-9]*)([ABCDEFGJK]|;?[0-9]*H)")
	x, y := 0, 0
	for {
		line, err := reader.ReadString('\n')
		if err == nil {
			x = 0
			if len(line) > 0 && line[len(line)-1:] == "\n" {
				line = line[:len(line)-1]
			}
			text := line
			for {
				indices := re.FindStringSubmatchIndex(text)
				printable := text
				if indices != nil && len(indices) > 4 {
					printable = text[:indices[0]]
				}
				screen = Print(screen, printable, x, y)
				x += len(printable)
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
						if count == 0 || countEnd == countStart { // To end
							if len(screen[y]) > x {
								screen[y] = screen[y][0:x]
							}
							screen = screen[0 : y+1]
						} else if count == 1 { // To begining
							screen[y] = strings.Repeat(" ", x) + screen[y][x:]
							for idx := range screen[0:y] {
								screen[idx] = ""
							}
						} else if count > 1 { // All
							screen = screen[:0]
							x = 0
							y = 0
						}
					case "K": // Erase in Line
						if count == 0 || countEnd == countStart { // To end
							screen[y] = screen[y][0:x]
						} else if count == 1 { // To beginning
							screen[y] = strings.Repeat(" ", x) + screen[y][x:]
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
					if len(text) > indices[1] {
						text = text[indices[1]:]
					} else {
						break
					}
				} else {
					break
				}
			}
			y += 1
		} else {
			if err != io.EOF {
				panic(fmt.Sprintf("Error %s", err))
			}
			break
		}
	}
	return screen
}

func Print(screen []string, text string, x int, y int) []string {
	for y >= len(screen) {
		screen = append(screen, "")
	}
	if y < len(screen) {
		suffix := ""
		if Len(screen[y]) > x+Len(text) {
			suffix = screen[y][x+Len(text):]
		}
		prefix := ""
		if x > 0 {
			if x < Len(screen[y]) {
				prefix = screen[y][0:x]
			} else {
				prefix = screen[y] + strings.Repeat(" ", x-Len(screen[y]))
			}
		}
		screen[y] = prefix + text + suffix
	}
	return screen
}

func Len(value string) int {
	return len(value)
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

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
