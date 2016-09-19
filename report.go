package main

import (
	"path/filepath"
	"os"
	"time"
	"fmt"
)

var logContextReport LoggerContext

var reportLocalContext ReportLocalContext

// ReportLocalContext ...
type ReportLocalContext struct {
	path string
}

const (
	FileOpen  = "opened"
	FileClose = "closed"
	FileRead  = "reading"
	FileWrite = "writing"
)

// ReportContext ...
type ReportContext struct {
	name      string
	uuid      string
	timestamp time.Time
	appnd     bool
	file      *os.File
	status    string
}

type ReportPub struct{}

// ReportInterface ...
type ReportInterface interface {
	New(name, uuid string, timestamp time.Time, appnd bool)
	Update(b []byte)
	UpdateString(s string)
	Report() []byte
}

// New fills a ReportContext struct attributes and creates the log file (as
// well as the parent directory, if not existent)
func (ctx *ReportContext) New(name, uuid string, timestamp time.Time, appnd bool) {
	ctx.name = name
	ctx.uuid = uuid
	ctx.timestamp = timestamp
	ctx.appnd = appnd

	// make dir if it doesn't exist
	dir := filepath.Join(reportLocalContext.path, name)
	_, err := os.Stat(dir)
	if err != nil {
		LogWar(logContextReport, "No dir %s found, making it", dir)
		err = os.Mkdir(dir, 0775)
		if err != nil {
			LogErr(logContextReport, "Unable to make dir %s", dir)
			panic(err)
		}
	}

	// create and open log file
	now := time.Now()
	fname := fmt.Sprintf("%d.%d.%d-%d.%d.%d-%s.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), uuid)
	fpath := filepath.Join(dir, fname)

	var f *os.File
	var perms int
	if appnd {
		perms = os.O_CREATE | os.O_APPEND | os.O_RDWR
	} else {
		perms = os.O_CREATE | os.O_RDWR
	}

	f, err = os.OpenFile(fpath, perms, 0666)
	if err != nil {
		LogErr(logContextReport, "Cannot create file %s", fpath)
		panic(err)
	}

	ctx.file = f
	ctx.status = FileOpen
	return
}

// UpdateString appends a string to the log file
func (ctx *ReportContext) UpdateString(s string) {

	if ctx.status != FileClose{
		for ctx.status != FileOpen{}
	} else {
		LogErr(logContextReport, "Cannot write on closed report")
		return
	}

	ctx.status = FileWrite
	defer func(){ ctx.status = FileOpen}()
	if ctx.appnd {
		_, err := ctx.file.Seek(0, 2)
		if err != nil {
			LogErr(logContextReport, "Cannot seek to end of file. %s", err)
			return
		}
	}

	_, err := ctx.file.WriteString(s)
	if err != nil {
		LogErr(logContextReport, "Cannot write to file %s", ctx.file)
		return
	}
	LogInf(logContextConfig, "%s", s)
}

// ReportInit initializes the Report module
func ReportInit(cm *ConfigModule) {
	logContextReport = LoggerContext{
		level: cm.GetLogLevel("report", 3),
		name:  "REPORT"}

	reportLocalContext = ReportLocalContext{
		path: cm.Get("report", "dir", "logs")}
}