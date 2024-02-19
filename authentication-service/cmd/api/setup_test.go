package main

import (
	"authentication/data"
	"os"
	"testing"
)

var testApp Config

func TestMain(m *testing.M) {
	repo := data.NewPostgresrTestRepository(nil)
	testApp.Repo = repo
	os.Exit(m.Run())
}