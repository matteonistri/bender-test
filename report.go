package main

import (
	"path/filepath"
	"os"
	"time"
	"fmt"
	"io/ioutil"
	"io"
	"errors"
	"strings"
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


// ReportInterface ...
type ReportInterface interface {
	New(name, uuid string, timestamp time.Time, appnd bool)
	Update(b []byte)
	UpdateString(s string)
	Report() []byte
}

// New fills a ReportContext struct attributes and creates the log file (as
// well as the parent directory, if not existent)
func (ctx *ReportContext) New(name, uuid string, timestamp time.Time, appnd bool) error{
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
			return errors.New("Unable to make report dir")
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
		return errors.New("Cannot create report file")
	}

	ctx.file = f
	ctx.status = FileOpen
	return nil
}

// UpdateString appends a string to the log file
func (ctx *ReportContext) UpdateString(s string) {
	//waiting for file if not available
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

// Read reads <size> byte, starting from <offset> for the specified test name
// and uuid
func (ctx *ReportContext) Read (offset, size int64) []byte{
	//waiting for file if not available
	if ctx.status != FileClose{
		for ctx.status != FileOpen{}
	} else {
		LogErr(logContextReport, "Cannot open this report")
		return nil
	}

	_, err := ctx.file.Seek(offset, 0)
	if err != nil {
		LogErr(logContextReport, "Cannot seek to specified offset. %s", err)
		return nil
	}

	var buf []byte
	if size < 0 {
		//read until EOF
		buf, err = ioutil.ReadAll(ctx.file)
		if err != nil {
			LogErr(logContextReport, "Cannot read from file. %s", err)
			return nil
		}

	} else {
		//read <size> bytes
		buf = make([]byte, size)
		_, err = io.ReadFull(ctx.file, buf)
		if err != nil {
			LogErr(logContextReport, "Error reading from file. %s", err)
			return nil
		}
	}

	return buf
}

//Close closes the report file
func (ctx *ReportContext) Close() {
	err := ctx.file.Close()
	if err != nil {
		LogErr(logContextReport, "Error while closing file. %s", err)
		return
	}
	ctx.status = FileClose
}

// Report returns the content of the log file as bytes
func Report(name, uuid string) ([]byte, error){
	dir := filepath.Join(reportLocalContext.path, name)
	_, err := os.Stat(dir)
	if err != nil {
		LogWar(logContextReport, "No logs available for script %s", name)
		err := errors.New("No logs available for script " + name)
		return nil, err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		LogWar(logContextReport, "Cannot stat %s", dir)
		err := errors.New("Cannot stat " + name)
		return nil, err
	}

	var out []byte
	for _, file := range files {
		if strings.Contains(file.Name(), uuid) {
			fpath := filepath.Join(dir, file.Name())

			out, err = ioutil.ReadFile(fpath)
			if err != nil {
				LogWar(logContextReport, "Cannot open log file %s", fpath)
				return nil, err
			}
		}
	}

	return out, nil
}

// List returns a list of available log files for the specified test name.
// Available log files are specified with their uuid and timestamp.
func ReportList(name string) ([][]string, error) {
	var out [][]string

	// look for dir 'name' in logs dir
	dir := filepath.Join(reportLocalContext.path, name)
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

// ReportInit initializes the Report module
func ReportInit(cm *ConfigModule) {
	logContextReport = LoggerContext{
		level: cm.GetLogLevel("report", 3),
		name:  "REPORT"}

	reportLocalContext = ReportLocalContext{
		path: cm.Get("report", "dir", "logs")}
}