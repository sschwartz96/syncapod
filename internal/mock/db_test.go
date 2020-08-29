package mock

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/sschwartz96/syncapod/internal/database"
)

type testStruct struct {
	Name  string
	Value int
}

func TestDB(t *testing.T) {
	logger := log.New(os.Stdout, "", 0)
	logger.SetOutput(ioutil.Discard)
	logger.SetFlags(0)
	db := DB{l: logger}
	db.Open(context.Background())

	// test insert
	err := db.Insert("test", testStruct{Name: "testName", Value: 123})
	if err != nil {
		t.Errorf("db.Insert(): %v", err)
	}
	err = db.Insert("test", testStruct{Name: "testName", Value: 0})
	if err != nil {
		t.Errorf("db.Insert(): %v", err)
	}

	// test find one
	ts := &testStruct{}
	err = db.FindOne("test", ts, &database.Filter{"name": "testName", "value": 123}, database.CreateOptions())
	if err != nil {
		t.Errorf("db.FindOne(): %v", err)
	}

	if ts.Name != "testName" {
		t.Errorf("name is not correct, got: %v", ts.Name)
	}

	if ts.Value != 123 {
		t.Errorf("value is not correct, got: %v", ts.Value)
	}

	// test find all
	var tSlice []testStruct
	err = db.FindAll("test", &tSlice, &database.Filter{"name": "testName"}, database.CreateOptions())

	if len(tSlice) != 2 {
		t.Errorf("length of array is %d, expecting 2", len(tSlice))
	} else if tSlice[0].Value != 123 {
		t.Errorf("tSlice[0] value is: %d, expecting 123", tSlice[0].Value)
	} else if tSlice[1].Value != 1 {
		t.Errorf("tSlice[1] value is: %d, expecting 0", tSlice[1].Value)
	}
}
