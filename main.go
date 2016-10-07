package main

import ( "github.com/matishsiao/goInfo"
		 "github.com/cvanderschuere/avahi-go"
	   )

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

	//Start autodiscover
	avahi.PublishService("Bender-test", "_http._tcp", 8080)
	LogInf(logContextMain, "Start autodiscover")

	// Start daemon
	DaemonInit(&sm, &cm)

}
