package main

import "github.com/matishsiao/goInfo"

var sm StatusModule
var logContextMain LoggerContext

func main() {
	// Put init here..
	LoggerModuleInit("bender-test")
	logContextMain = LoggerContext{
		name:  "MAIN",
		level: 3}
	sm = StatusModuleInit("bender-test")

	// Start daemon
	gi := goInfo.GetInfo()
	LogInf(logContextMain, "== Bender test ==")
	LogInf(logContextMain, "Run on: %v", gi)
	LogInf(logContextMain, "Start")
	DaemonInit("localhost", "8080")
}
