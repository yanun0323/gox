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

func NoError(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 || len(msg[0]) == 0 {
			println(err)
			log.Fatal(err)
		}

		println(fmt.Errorf("%s, err: %w", msg[0], err))
		log.Fatalf("%s, err: %+v", msg[0], err)
	}
}

var helper = helperInstance{}

type helperInstance struct{}

func (helperInstance) environmentPrint() {
	for _, ev := range []string{"GOARCH", "GOOS", "GOFILE", "GOLINE", "GOPACKAGE", "DOLLAR"} {
		fmt.Println("\t", ev, "=", os.Getenv(ev))
	}
}

func (helperInstance) firstLowerCase(s string) string {
	if s[0] <= 'Z' && s[0] >= 'A' {
		buf := []byte(s)
		gap := byte('a' - 'A')
		buf[0] = buf[0] + gap
		return string(buf)
	}
	return s
}

func (helperInstance) firstUpperCase(s string) string {
	if s[0] <= 'z' && s[0] >= 'a' {
		buf := []byte(s)
		gap := byte('a' - 'A')
		buf[0] = buf[0] - gap
		return string(buf)
	}
	return s
}

func (helperInstance) isFirstUpperCase(s string, ignoreChars ...byte) bool {
	tidied := helper.tidyString(s, ignoreChars...)

	return tidied[0] >= 'A' && tidied[0] <= 'Z'
}

func (helperInstance) setupLog() {
	log.SetFlags(0)
	log.SetPrefix(_commandName + ": ")
	flag.Usage = Usage
	flag.Parse()
}

func (helperInstance) requireDestination() {
	if len(*_destination) == 0 {
		flag.Usage()
		NoError(errors.New("entity/use/repo at least one param provide"))
	}

	if len(*_package) == 0 {
		flag.Usage()
		NoError(errors.New("package not define"))
	}
}

func (helperInstance) debugPrint() {
	if *_debug {
		println()
		println("\t", "replace", "=", *_replace)
		println("\t", "name", "=", *_destination)
		println()
	}
}

func (helperInstance) getDir() (currentDirectory string, currentFile string, e error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	file := filepath.Join(cwd, os.Getenv("GOFILE"))
	return cwd, file, nil
}

func (helperInstance) getSourceImportString() (string, error) {
	cwd, file, err := helperInstance{}.getDir()
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

func (helperInstance) EqualFold(a, b string, ignoreChars ...byte) bool {
	a = helper.tidyString(a, ignoreChars...)
	b = helper.tidyString(b, ignoreChars...)

	if len(a) == 0 || len(b) == 0 {
		return false
	}

	return strings.EqualFold(a, b)
}

func (helperInstance) tidyString(s string, removeChars ...byte) string {
	tidied := s
	for _, char := range removeChars {
		tidied = strings.ReplaceAll(tidied, string(char), "")
	}

	return tidied
}

func (helperInstance) insertString(s, prefix, insert string) string {
	if strings.HasPrefix(s, prefix) {
		return prefix + insert + strings.TrimPrefix(s, prefix)
	}

	return insert + s
}
