package main

import "github.com/matishsiao/goInfo"

func main() {
	var sm StatusModule
	var cm ConfigModule

	// init modules
	LoggerModuleInit("bender-test")
	ConfigInit(&cm, "bender-test")

	logContextMain := LoggerContext{
		name:  "MAIN",
		level: cm.GetLogLevel("general", 3)}

	StatusModuleInit(&sm, &cm)
	WorkerInit(&sm)

	// Start daemon
	gi := goInfo.GetInfo()
	LogInf(logContextMain, "== Bender test ==")
	LogInf(logContextMain, "Run on: %v", gi)
	DaemonInit(&sm, &cm)
}
