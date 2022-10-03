package main

import (
	"authentication/data"
	"os"
	"testing"
)

var testApp Config

// Must be called TestMain
func TestMain(m *testing.M) {
	// We are just testing the handlers, NOT the DB, therefore, we just pass nil to the DB param
	repo := data.NewTestPostgresRepository(nil)
	testApp.Repo = repo

	// Run all of the tests
	os.Exit(m.Run())
}
