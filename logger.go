package main

import (
	"io"
	"log"
	"os"
)

// Logs lever enum
const (
	LogError   = 0
	LogWarning = 1
	LogInfo    = 2
	LogDebug   = 3
)

// LoggerContext ...
type LoggerContext struct {
	level int
	name  string
}

// LoggerModuleInit initializes the logger module
func LoggerModuleInit(logName string) {
	logName += ".log"
	logfile, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", logfile, ":", err)
	}
	multilog := io.MultiWriter(logfile, os.Stdout)
	log.SetOutput(multilog)
}

// LogErr calls log.Printf to print to the logger. Arguments following the
// first are handled in the manner of fmt.Printf.
func LogErr(c LoggerContext, s string, v ...interface{}) {
	if c.level >= LogError {
		log.Printf("["+c.name+"]"+" ERR: "+s, v...)
	}
}

// LogWar calls log.Printf to print to the logger. Arguments following the
// first are handled in the manner of fmt.Printf.
func LogWar(c LoggerContext, s string, v ...interface{}) {
	if c.level >= LogWarning {
		log.Printf("["+c.name+"]"+" WARN: "+s, v...)
	}
}

// LogInf calls log.Printf to print to the logger. Arguments following the
// first are handled in the manner of fmt.Printf.
func LogInf(c LoggerContext, s string, v ...interface{}) {
	if c.level >= LogInfo {
		log.Printf("["+c.name+"]"+" INFO: "+s, v...)
	}
}

// LogDeb calls log.Printf to print to the logger. Arguments following the
// first are handled in the manner of fmt.Printf.
func LogDeb(c LoggerContext, s string, v ...interface{}) {
	if c.level <= LogDebug {
		log.Printf("["+c.name+"]"+" DEBUG: "+s, v...)
	}
}

// LogFatal writes to logfile and terminates the program when the called
// interface ends
func LogFatal(c LoggerContext, v ...interface{}) {
	log.Fatalf("["+c.name+"]"+" FATAL: %s", v...)
}
