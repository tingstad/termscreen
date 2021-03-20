package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
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
	var bufioReader *bufio.Reader = bufio.NewReader(reader)
	var myReader MyReader = bufioReader
	return Capture(myReader)
}

type State struct {
	x, y   int
	screen [][]string
}

func CaptureShim(reader MyReader) []string {
	lines := Capture(reader)
	return lines
}
func Capture(reader MyReader) []string {
	screen := make([]string, 0)
	x, y := 0, 0
	for {
		line, err := reader.ReadString('\n')
		if err == nil {
			fmt.Printf("Read %d %d  %s\n", y, x, line)
			esc := "\x1b"
			re := regexp.MustCompile(esc + "\\[([0-9]*)([ABCDEFGJK]|;[0-9]*H)")
			indicies := re.FindStringSubmatchIndex(line)
			if indicies != nil && len(indicies) > 4 {
				fmt.Printf("indices %d\n", indicies)
				countStart := indicies[2]
				countEnd := indicies[3]
				start := indicies[4]
				runes := []rune(line)
				code := string(runes[start : start+1])
				count := string(runes[countStart:countEnd])
				switch code {
				case "A":
					fmt.Printf("AAAA " + count)
				default:
					fmt.Printf("substr %s %s", code, count)
				}
				fmt.Printf("index2 %d %d\n", indicies[0], indicies[1])
			}
			screen = Print(screen, line, x, y)
			y += 1
		} else {
			if err != io.EOF {
				fmt.Printf("Rrrot %s", err)
			}
			break
		}
	}
	for _, line := range screen {
		if line[len(line)-1:] != "\n" {
			line = line + "\n"
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
		if len(screen[y]) <= x+len(text) {
			screen[y] = text
		}
	}
	return screen
}

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
