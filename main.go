package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
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
	return CaptureReaderNew(reader)
}
func CaptureReaderNew(reader io.Reader) []string {
	var bufioReader *bufio.Reader = bufio.NewReader(reader)
	var myReader MyReader = bufioReader
	lines := Capture(myReader)
	return lines
}

func Capture(reader MyReader) []string {
	screen := make([]string, 0)
	esc := "\x1b"
	re := regexp.MustCompile(esc + "\\[([0-9]*)([ABCDEFGJK]|;[0-9]*H)")
	x, y := 0, 0
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 && line[len(line)-1:] == "\n" {
			line = line[:len(line)-1]
		}
		if err == nil {
			fmt.Printf("Read %d %d  %s\n", y, x, line)
			text := line
			for {
				indices := re.FindStringSubmatchIndex(text)
				if indices != nil && len(indices) > 4 {
					text = text[:indices[0]]
				}
				screen = Print(screen, text, x, y)
				if indices != nil && len(indices) > 4 {
					fmt.Printf("indices %d\n", indices)
					countStart := indices[2]
					countEnd := indices[3]
					start := indices[4]
					code := string(line[start : start+1])
					count := string(line[countStart:countEnd])
					switch code {
					case "A": // Up
						fmt.Printf("AAAA " + count)
					case "B": // Down
						y += 1
					default:
						fmt.Printf("substr %s %s", code, count)
					}
					fmt.Printf("index2 %d %d\n", indices[0], indices[1])
					text = text[indices[1]:]
				} else {
					break
				}
			}
			y += 1
		} else {
			if err != io.EOF {
				fmt.Printf("Rrrot %s", err)
			}
			break
		}
	}
	return screen
}

func Print(screen []string, text string, x int, y int) []string {
	for y > len(screen) {
		screen = append(screen, "")
	}
	if y == len(screen) {
		screen = append(screen, text)
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

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
