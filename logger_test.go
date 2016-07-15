package main

import (
    "testing"
    "fmt"
    "path/filepath"
    "time"
    "os"
)

var test_path string
var file_name string
var job Job

func TestMain(m *testing.M) {
    now := time.Now()
    job = Job{Script:"script", Uuid: "uuid"}
    test_path, _ = filepath.Abs(filepath.Join("log", "script"))
    file_name = fmt.Sprintf("%d.%d.%d-%d.%d.%d-%s.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), job.Uuid)
    code := m.Run()
    os.RemoveAll(test_path)
    os.Exit(code)
}

func TestWriteLog(t *testing.T){
    go WriteLog()
    jobDone <- job
    file_path:= filepath.Join(test_path, file_name)

    time.Sleep(20 * time.Millisecond)
    if _, err := os.Stat(file_path); os.IsNotExist(err){
        t.Error("Test failed")
    }
}

func TestReadLog(t *testing.T){
    expected := "{\"Script\":\"script\",\"Path\":\"\",\"Args\":null,\"Uuid\":\"uuid\"," +
                "\"Output\":\"\",\"Exit\":\"\",\"Request\":\"0001-01-01T00:00:00Z\"," +
                "\"Start\":\"0001-01-01T00:00:00Z\",\"Finish\":\"0001-01-01T00:00:00Z\"," +
                "\"Status\":\"\"}"

    file_path := filepath.Join(test_path, file_name)
    actual, _ := ReadLog(file_path)

    if actual != expected{
        t.Error("Test failed")
    }
}

func TestReadLogWrongFile(t *testing.T){
    file_path := filepath.Join(test_path, "inexistant_file.log")
    _, err := ReadLog(file_path)

    if err == nil{
        t.Error("Test failed")
    }
}

func TestReadLogWrongDir(t *testing.T){
    file_path := filepath.Join("log", "inexistant_dir", file_name)
    _, err := ReadLog(file_path)

    if err == nil{
        t.Error("Test failed")
    }
}