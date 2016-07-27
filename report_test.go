package main

import (
	"bytes"
	"testing"
	"time"
)

func TestWriteString(t *testing.T) {
	var cm ConfigModule
	LoggerModuleInit("dummy")
	ConfigInit(&cm, "dummy")

	expectedS := "foobarbiz"
	expectedB := []byte(expectedS)

	r := ReportContext{}
	r.New("foo", "0000", time.Now(), false)
	r.UpdateString("foo")
	r.UpdateString("bar")
	r.UpdateString("biz")

	actualB := r.Report()
	actualS := string(actualB)

	if !bytes.Equal(expectedB, actualB) {
		t.Errorf("Expected %s, got %s\n", expectedB, actualB)
	}
	if expectedS != actualS {
		t.Errorf("Expected %s, got %s\n", expectedS, actualS)
	}
}

func TestWriteBytes(t *testing.T) {
	var cm ConfigModule
	LoggerModuleInit("dummy")
	ConfigInit(&cm, "dummy")

	expectedS := "foobarbiz"
	expectedB := []byte(expectedS)

	r := ReportContext{}
	r.New("foo", "0001", time.Now(), true)
	r.Update([]byte("foo"))
	r.Update([]byte("bar"))
	r.Update([]byte("biz"))

	actualB := r.Report()
	actualS := string(actualB)

	if !bytes.Equal(expectedB, actualB) {
		t.Errorf("Expected %s, got %s\n", expectedB, actualB)
	}
	if expectedS != actualS {
		t.Errorf("Expected %s, got %s\n", expectedS, actualS)
	}
}
