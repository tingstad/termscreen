package main

import (
	"fmt"
	"github.com/tingstad/termscreen"
	"os"
)

func main() {
	lines := termscreen.CaptureReader(os.Stdin)
	for _, line := range lines {
		fmt.Printf("%s\n", line)
	}
}
