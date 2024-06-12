package main

import (
	"os"
	"testing"
)

func TestGetSourceImportString(t *testing.T) {
	os.Setenv("GOFILE", "helper")
	os.Setenv("GOPACKAGE", "usecase")

	dir, err := helper.getSourceImportString()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if dir != "gox/cmd/usecase" {
		t.Fatalf("mismatch: %s", dir)
	}
}
