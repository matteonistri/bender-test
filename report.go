package main

import (
	"bufio"
	"bytes"
	"fmt"
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
	file      *os.File
}

type ReportInterface interface {
	New(name, uuid string, timestamp time.Time, appnd bool)
	Update(b []byte)
	Append(s string)
	Report() bufio.Reader
}

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
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		LogErr(logContextReport, "Cannot create file %s", fpath)
		panic(err)
	}

	ctx.file = f
	return
}

func (ctx *ReportContext) Update(b []byte) {
	w := bufio.NewWriter(ctx.file)
	_, err := w.Write(b)
	if err != nil {
		LogErr(logContextReport, "Cannot write to file %s", ctx.file)
		panic(err)
	}
	w.Flush()
}

func (ctx *ReportContext) Append(s string) {
	w := bufio.NewWriter(ctx.file)
	_, err := w.WriteString(s)
	if err != nil {
		LogErr(logContextReport, "Cannot write to file %s", ctx.file)
		panic(err)
	}
	w.Flush()
}

func (ctx *ReportContext) Report() bytes.Buffer {
	// read log file
	// make it a buffer
	// return it
	var b bytes.Buffer
	return b
}

func ReportInit(cm *ConfigModule) {
	logContextReport = LoggerContext{
		level: cm.GetLogLevel("report", 3),
		name:  "REPORT"}

	report_localContext = ReportLocalContext{
		path: cm.Get("report", "dir", "logs")}
}
