package main

import (
	"fmt"
	"time"

	"github.com/matishsiao/goInfo"
)

func main() {
	gi := goInfo.GetInfo()
	LogAppendLine(fmt.Sprintf("== Bender test =="))
	LogAppendLine(fmt.Sprintf("Run on: %v", gi))
	LogAppendLine(fmt.Sprintf("START  %s", time.Now()))

	// Put init here..
	DaemonInit("", "8080")
}
