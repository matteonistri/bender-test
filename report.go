package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var logContextReport LoggerContext
var report_localContext ReportLocalContext

type ReportLocalContext struct {
	path string
}

type ReportContext struct {
	name      string
	uuid      string
	timestamp time.Time
	appnd     bool
	file      string
}

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
	ctx.appnd = true

	// make dir if it doesn't exist
	dir := filepath.Join(report_localContext.path, name)
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
	if appnd {
		f, err = os.OpenFile(fpath, os.O_CREATE|os.O_APPEND, 0666)
	} else {
		f, err = os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0666)
	}
	if err != nil {
		LogErr(logContextReport, "Cannot create file %s", fpath)
		panic(err)
	}

	ctx.file = fpath
	defer f.Close()
	return
}

// Update appends bytes to the log file
func (ctx *ReportContext) Update(b []byte) {
	f, err := os.OpenFile(ctx.file, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		LogErr(logContextReport, "Cannot open file %s", ctx.file)
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write(b)
	if err != nil {
		LogErr(logContextReport, "Cannot write to file %s", ctx.file)
		panic(err)
	}

	w.Flush()
}

// UpdateString appends a string to the log file
func (ctx *ReportContext) UpdateString(s string) {
	f, err := os.OpenFile(ctx.file, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		LogErr(logContextReport, "Cannot open file %s", ctx.file)
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(s)
	if err != nil {
		LogErr(logContextReport, "Cannot write to file %s", ctx.file)
		panic(err)
	}

	w.Flush()
}

// Report returns the content of the log file as bytes
func (ctx *ReportContext) Report() []byte {
	out, err := ioutil.ReadFile(ctx.file)
	if err != nil {
		LogErr(logContextReport, "IO error while reading file %s", ctx.file)
		panic(err)
	}

	return out
}

// ReportInit initializes the Report module
func ReportInit(cm *ConfigModule) {
	logContextReport = LoggerContext{
		level: cm.GetLogLevel("report", 3),
		name:  "REPORT"}

	report_localContext = ReportLocalContext{
		path: cm.Get("report", "dir", "logs")}
}
