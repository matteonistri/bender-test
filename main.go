package main

import (
	"fmt"
	"time"
)

func main() {
	LogAppendLine(fmt.Sprintf("== Bender test =="))
	LogAppendLine(fmt.Sprintf("START  %s", time.Now()))

	// Put init here..
	DaemonInit("", "8080")
}
