package main

import "github.com/matishsiao/goInfo"

var sm StatusModule
var logContextMain LoggerContext

func main() {
	// init modules
	LoggerModuleInit("bender-test")
	cfg := ConfigInit("bender-test")

	logContextMain = LoggerContext{
		name:  "MAIN",
		level: cfg.generalLogLevel}
	sm = StatusModuleInit(cfg.statusName)

	// Start daemon
	gi := goInfo.GetInfo()
	LogInf(logContextMain, "== Bender test ==")
	LogInf(logContextMain, "Run on: %v", gi)
	DaemonInit(cfg.daemonLogLevel, "0.0.0.0", "8080")
}
