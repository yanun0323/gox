package main

import (
	"os"
	"testing"
)

func TestGetSourceImportString(t *testing.T) {
	{
		os.Setenv("GOFILE", "helper")
		os.Setenv("GOPACKAGE", "usecase")

		alias, dir, err := helper.getSourceImportString()
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if alias != "usecase" {
			t.Fatalf("alias mismatch: %s", alias)
		}

		if dir != "github.com/yanun0323/gox/cmd/domaingen" {
			t.Fatalf("import path mismatch: %s", dir)
		}
	}

	{
		os.Setenv("GOFILE", "helper")
		os.Setenv("GOPACKAGE", "domaingen")

		alias, dir, err := helper.getSourceImportString()
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if alias != "" {
			t.Fatalf("alias mismatch: %s", alias)
		}

		if dir != "github.com/yanun0323/gox/cmd/domaingen" {
			t.Fatalf("import path mismatch: %s", dir)
		}
	}
}

func TestGetGoModulePath(t *testing.T) {
	name, err := helper.getGoModulePath()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	t.Log("module path:", name)
}

func TestGetGoModuleName(t *testing.T) {
	name, err := helper.getModuleName()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	t.Log("module name:", name)

	if name != "github.com/yanun0323/gox" {
		t.Fatalf("mismatch: %s", name)
	}
}
