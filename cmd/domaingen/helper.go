package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func environmentPrint() {
	for _, ev := range []string{"GOARCH", "GOOS", "GOFILE", "GOLINE", "GOPACKAGE", "DOLLAR"} {
		fmt.Println("\t", ev, "=", os.Getenv(ev))
	}
}

func NoError(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 || len(msg[0]) == 0 {
			log.Fatal(err)
		}
		log.Fatalf("%s, err: %+v", msg[0], err)
	}
}

func firstLowerCase(s string) string {
	if s[0] <= 'Z' && s[0] >= 'A' {
		buf := []byte(s)
		gap := byte('a' - 'A')
		buf[0] = buf[0] + gap
		return string(buf)
	}
	return s
}

func firstUpperCase(s string) string {
	if s[0] <= 'z' && s[0] >= 'a' {
		buf := []byte(s)
		gap := byte('a' - 'A')
		buf[0] = buf[0] - gap
		return string(buf)
	}
	return s
}

func isFirstUpperCase(s string) bool {
	if s[0] == '*' && len(s) >= 2 {
		return s[1] >= 'A' && s[1] <= 'Z'
	}
	return s[0] >= 'A' && s[0] <= 'Z'
}

func setupLog() {
	log.SetFlags(0)
	log.SetPrefix(_commandName + ": ")
	flag.Usage = Usage
	flag.Parse()
}

func requireDestination() {
	if len(*_destination) == 0 {
		flag.Usage()
		NoError(errors.New("entity/use/repo at least one param provide"))
	}

	if len(*_package) == 0 {
		flag.Usage()
		NoError(errors.New("package not define"))
	}
}

func debugPrint() {
	if *_debug {
		println()
		println("\t", "replace", "=", *_replace)
		println("\t", "name", "=", *_destination)
		println()
	}
}

func getDir() (currentDirectory string, currentFile string, e error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	file := filepath.Join(cwd, os.Getenv("GOFILE"))
	return cwd, file, nil
}

func getSourceImportString() (string, error) {
	cwd, file, err := getDir()
	if err != nil {
		return "", err
	}

	cwdSplit := strings.SplitAfter(cwd, string(os.PathSeparator))
	cwdSplit = cwdSplit[:len(cwdSplit)-1]
	cwd = strings.Join(cwdSplit, "")

	ss := strings.SplitAfter(file, string(os.PathSeparator))

	for len(ss) != 0 {
		ss = ss[:len(ss)-1]
		d := strings.Join(ss, "")
		_, err := os.Stat(d + "go.mod")
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return "", err
		}

		prefix := strings.Join(ss[:len(ss)-1], "")
		return strings.TrimPrefix(cwd, prefix) + os.Getenv("GOPACKAGE"), nil
	}

	return "", errors.New("project not found")
}
