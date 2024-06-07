package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
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

func requireNoError(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 || len(msg[0]) == 0 {
			log.Fatal(err)
		}
		log.Fatalf("%s, err: %+v", msg[0], err)
	}
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
	return s[0] >= 'A' && s[0] <= 'Z'
}

func SaveAst(fset *token.FileSet, f *ast.File, targetFullPath string) error {
	var buf strings.Builder
	if err := format.Node(&buf, fset, f); err != nil {
		return fmt.Errorf("format node, err: %w", err)
	}

	dirPath := filepath.Dir(targetFullPath)
	if err := os.MkdirAll(dirPath, 0o777); err != nil {
		return fmt.Errorf("make directory: %s, err: %w", dirPath, err)
	}

	file, err := os.Create(targetFullPath)
	if err != nil {
		return fmt.Errorf("create file: %s, err: %w", targetFullPath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(buf.String()); err != nil {
		return fmt.Errorf("write ast buffer into file: %s , err: %w", targetFullPath, err)
	}

	return nil
}
