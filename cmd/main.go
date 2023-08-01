// Copyright (C) 2021-2023 Richard H. Tingstad
// This program is free software: you can redistribute it and/or modify it under the terms of the
// GNU General Public License as published by the Free Software Foundation, version 3.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
// without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU General Public License for more details.

package main

import (
	"fmt"
	"github.com/tingstad/termscreen"
	"os"
)

func main() {
	lines := termscreen.Capture(os.Stdin)
	for _, line := range lines {
		fmt.Printf("%s\n", line)
	}
}
