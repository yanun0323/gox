package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"errors"
)

func getDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get word dir, err: %+v", err)
	}

	dir := filepath.Join(cwd, os.Getenv("GOFILE"))
	return dir, nil
}

func requireNoError(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 || len(msg[0]) == 0 {
			log.Fatal(err)
		}
		log.Fatalf("%s, err: %+v", msg[0], err)
	}
}

func cleanStringQuote(s string) string {
	return strings.Trim(strings.Trim(s, "\""), "'")
}

func currentPackage() Package {
	return Package(os.Getenv("GOPACKAGE"))
}

func findInternalPath(dir string) (string, error) {
	spans := strings.Split(dir, "internal")
	if len(spans) == 1 {
		return "", errors.New("missing internal folder in working path")
	}

	internal := filepath.Join(spans[0], "internal")
	if *_debug {
		println("internal folder path:", internal)
	}

	return internal, nil
}

func findProjectPath(dir string) (string, error) {
	spans := strings.Split(dir, "internal")
	if len(spans) == 1 {
		return "", errors.New("missing internal folder in working path")
	}

	return spans[0], nil
}
