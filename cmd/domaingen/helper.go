package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func (helperInstance) requireTag() {
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
	cwd = filepath.Clean(cwd)
	file := filepath.Join(cwd, os.Getenv("GOFILE"))
	return cwd, file, nil
}

func (helperInstance) findProjectDir() (string, error) {
	_, filePath, err := helperInstance{}.getDir()
	if err != nil {
		return "", err
	}

	filePathSpan := strings.SplitAfter(filePath, string(os.PathSeparator))
	for len(filePathSpan) != 0 {
		filePathSpan = filePathSpan[:len(filePathSpan)-1]
		d := strings.Join(filePathSpan, "")
		_, err := os.Stat(d + "go.mod")
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return "", err
		}

		return filepath.Join(filePathSpan...), nil
	}

	return "", errors.New("project not found")
}

func (h helperInstance) getSourceImportString() (alias, importPath string, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	moduleName, err := h.getModuleName()
	if err != nil {
		return "", "", err
	}

	projectDir, err := h.findProjectDir()
	if err != nil {
		return "", "", err
	}

	relativePath := strings.TrimPrefix(currentDir, projectDir)
	relativePathSpan := strings.Split(relativePath, string(os.PathSeparator))

	ali := os.Getenv("GOPACKAGE")
	if len(relativePathSpan) != 0 && relativePathSpan[len(relativePathSpan)-1] == ali {
		ali = ""
	}

	return ali, strings.Join([]string{moduleName, strings.Join(relativePathSpan, "/")}, ""), nil

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

func (helperInstance) getGoModulePath() (string, error) {
	env, err := exec.Command("go", "env").Output()
	if err != nil {
		return "", err
	}

	// find GOMOD keyword
	rows := strings.Split(string(env), "\n")
	for _, row := range rows {
		span := strings.Split(row, "=")
		if len(span) != 2 || span[0] != "GOMOD" {
			continue
		}

		mod := strings.Trim(span[1], "'")
		mod = strings.Trim(mod, "\"")
		if len(mod) == 0 {
			return "", errors.New("go.mod not found, please run this program in the root folder of a go module project")
		}

		return mod, nil
	}

	return "", errors.New("go.mod not found, please run this program in the root folder of a go module project")
}

func (h helperInstance) getModuleName() (string, error) {
	path, err := h.getGoModulePath()
	if err != nil {
		return "", err
	}

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("go mod file not found, err: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			name := strings.TrimPrefix(line, "module ")
			return strings.TrimSpace(name), nil
		}
	}

	return "", errors.New("module not found")
}
