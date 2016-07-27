package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	UpdateString(s string)
	Report() []byte
}

type ReportPubInterface interface {
	List(name string) ([][]string, error)
	Read(name, uuid string, size, offset int64) ([]byte, error)
}

type ReportPub struct{}

// New fills a ReportContext struct attributes and creates the log file (as
// well as the parent directory, if not existent)
func (ctx *ReportContext) New(name, uuid string, timestamp time.Time, appnd bool) {
	ctx.name = name
	ctx.uuid = uuid
	ctx.timestamp = timestamp
	ctx.appnd = appnd

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
	return
}

// Update appends bytes to the log file
func (ctx *ReportContext) Update(b []byte) {
	if ctx.appnd {
		_, err := ctx.file.Seek(0, 2)
		if err != nil {
			LogErr(logContextReport, "Cannot seek to end of file. %s", err)
			panic(err)
		}
	}

	_, err := ctx.file.Write(b)
	if err != nil {
		LogErr(logContextReport, "Cannot write to file %s", ctx.file)
		panic(err)
	}
	LogInf(logContextConfig, "%s", b)
}

// UpdateString appends a string to the log file
func (ctx *ReportContext) UpdateString(s string) {
	if ctx.appnd {
		_, err := ctx.file.Seek(0, 2)
		if err != nil {
			LogErr(logContextReport, "Cannot seek to end of file. %s", err)
			panic(err)
		}
	}

	_, err := ctx.file.WriteString(s)
	if err != nil {
		LogErr(logContextReport, "Cannot write to file %s", ctx.file)
		panic(err)
	}
	LogInf(logContextConfig, "%s", s)
}

// Report returns the content of the log file as bytes
func (ctx *ReportContext) Report() []byte {
	_, err := ctx.file.Seek(0, 0)
	if err != nil {
		LogErr(logContextReport, "Cannot seek to start of file. %s", err)
		panic(err)
	}

	var out []byte
	buff := make([]byte, 1024)
	for {
		n, err := ctx.file.Read(buff)
		if err != nil && err != io.EOF {
			LogErr(logContextReport, "Cannot read file. %s", err)
			panic(err)
		}
		if n == 0 {
			break
		}
		out = append(out, buff[:n]...)
	}

	return out
}

// List returns a list of available log files for the specified test name.
// Available log files are specified with their uuid and timestamp.
func (rp *ReportPub) List(name string) ([][]string, error) {
	var out [][]string

	// look for dir 'name' in logs dir
	dir := filepath.Join(report_localContext.path, name)
	_, err := os.Stat(dir)
	if err != nil {
		LogWar(logContextReport, "No logs available for script %s", name)
		err := errors.New("No logs available for script " + name)
		return nil, err
	}

	// get a list of files
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		LogWar(logContextReport, "Cannot stat %s", dir)
		err := errors.New("Canno stat " + name)
		return nil, err
	}

	// build list from file name timestamp-uuid
	for _, file := range files {
		//		f := make([]string, 2)
		LogDeb(logContextReport, "Found file: %s", file.Name())
		x := strings.Split(file.Name(), "-")
		tr := x[:2]
		id := x[3:]

		timestamp := strings.Join(tr, "-")
		uuid := strings.Join(id, "-")
		uuid = string(strings.Split(uuid, ".")[0])

		LogDeb(logContextReport, "  -timestamp: %s", timestamp)
		LogDeb(logContextReport, "  -uuid: %s", uuid)

		t := make([]string, 2)
		t[0] = uuid
		t[1] = timestamp
		out = append(out, t)
	}

	return out, nil
}

// Read reads <size> byte, starting from <offset> for the specified test name
// and uuid
func (rp *ReportPub) Read(name, uuid string, size, offset int64) ([]byte, error) {
	var out []byte
	// attempt to open file

	// attempt to read from it
	return out, nil
}

// ReportInit initializes the Report module
func ReportInit(cm *ConfigModule) {
	logContextReport = LoggerContext{
		level: cm.GetLogLevel("report", 3),
		name:  "REPORT"}

	report_localContext = ReportLocalContext{
		path: cm.Get("report", "dir", "logs")}
}
