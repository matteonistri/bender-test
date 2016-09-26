package main

import "github.com/matishsiao/goInfo"

func main() {
	var sm StatusModule
	var cm ConfigModule

	ConfigInit(&cm, "bender-test")

	LoggerModuleInit("bender-test")
	logContextMain := LoggerContext{
		name:  "MAIN",
		level: cm.GetLogLevel("general", 3)}

	LogInf(logContextMain, "== Bender test ==")
	LogInf(logContextMain, "Run on: %v", goInfo.GetInfo())

	// init modules
	StatusModuleInit(&sm, &cm)
	RunnerInit(&cm)
	WorkerInit(&sm)
	ReportInit(&cm)
	WebsocketInit(&sm)

	// Start daemon
	DaemonInit(&sm, &cm)
}
