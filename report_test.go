package main

import (
	"reflect"
	"testing"
	"time"
)

func TestWriteString(t *testing.T) {
	var cm ConfigModule
	LoggerModuleInit("dummy")
	ConfigInit(&cm, "dummy")

	expectedS := "foobarbiz"

	r := ReportContext{}
	r.New("foo", "0000", time.Now(), true)
	r.UpdateString("foo")
	r.UpdateString("bar")
	r.UpdateString("biz")

	actualB := r.Report()
	actualS := string(actualB)
	if actualS != expectedS {
		t.Error("Expected %s, got %s", expectedS, actualS)
	}
}

func TestWriteBytes(t *testing.T) {
	var cm ConfigModule
	LoggerModuleInit("dummy")
	ConfigInit(&cm, "dummy")

	expectedS := "foobarbiz"
	expectedB := []byte(expectedS)

	r := ReportContext{}
	r.New("foo", "0000", time.Now(), true)
	r.Update([]byte("foo"))
	r.Update([]byte("bar"))
	r.Update([]byte("biz"))

	actualB := r.Report()
	actualS := string(actualB)
	if reflect.DeepEqual(expectedB, actualB) {
		t.Error("Expected %s, got %s", expectedS, actualS)
	}
}
