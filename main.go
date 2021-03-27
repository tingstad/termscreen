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
	Use(strconv.Atoi("1"))
	fmt.Printf("Started")
	lines := CaptureReader(os.Stdin)
	fmt.Printf("\n\n")
	for _, line := range lines {
		fmt.Printf("Line %s", line)
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
			fmt.Printf("Read %d %d  %s\n", y, x, line)
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
					fmt.Printf("indices %d\n", indices)
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
						y -= count
						if y < 0 {
							y = 0
						}
					case "B": // Down
						y += count
					case "C": // Forward
						x += count
					case "D": // Back
						x -= count
						if x < 0 {
							x = 0
						}
					case "E": // Next line
						y += count
						x = 0
					case "F": // Previous line
						y -= count
						x = 0
					case "G": // Column
						x = count
					default:
						fmt.Printf("substr %s %d", code, count)
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
					fmt.Printf("index2 %d %d\n", indices[0], indices[1])
					if len(text) > indices[1]+1 {
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
	fmt.Printf("DEBUG Printf %d %d %s\n", x, y, text)
	for y >= len(screen) {
		screen = append(screen, "")
	}
	if y < len(screen) {
		suffix := ""
		if len(screen[y]) > x+len(text) {
			suffix = screen[y][x+len(text):]
		}
		prefix := ""
		if x > 0 {
			if x < len(screen[y]) {
				prefix = screen[y][0:x]
			} else {
				prefix = screen[y] + strings.Repeat(" ", x-len(screen[y]))
			}
		}
		screen[y] = prefix + text + suffix
	}
	return screen
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

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
