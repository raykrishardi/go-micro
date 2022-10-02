package main

import (
	"os"
	"testing"
)

// Must be called TestMain
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
