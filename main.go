package main

import "github.com/matishsiao/goInfo"

var logContextMain LoggerContext

func main() {
	var sm StatusModule
	var cm ConfigModule

	// init modules
	LoggerModuleInit("bender-test")
	ConfigInit(&cm, "bender-test")

	logContextMain = LoggerContext{
		name:  "MAIN",
		level: cm.GetLogLevel("general", 3)}

	StatusModuleInit(&sm, cm.Get("status", "servername", "bender"))

	// Start daemon
	gi := goInfo.GetInfo()
	LogInf(logContextMain, "== Bender test ==")
	LogInf(logContextMain, "Run on: %v", gi)
	//DaemonInit(cfg.daemonLogLevel, "0.0.0.0", "8080")
}
