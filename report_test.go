package main

import (
	"bytes"
	"testing"
	"time"
	"os"
)

func TestMain(m *testing.M) {
    retCode := m.Run()
    os.RemoveAll("foo")
    os.Exit(retCode)
}

func TestWriteString(t *testing.T) {

	expectedS := "foobarbiz"
	expectedB := []byte(expectedS)

	r := ReportContext{}
	r.New("foo", "0000", time.Now(), false)
	r.UpdateString("foo")
	r.UpdateString("bar")
	r.UpdateString("biz")

	actualB, _ := Report("foo", "0000")
	actualS := string(actualB)

	if !bytes.Equal(expectedB, actualB) {
		t.Errorf("Expected %s, got %s\n", expectedB, actualB)
	}
	if expectedS != actualS {
		t.Errorf("Expected %s, got %s\n", expectedS, actualS)
	}
}