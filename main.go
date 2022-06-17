package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var rxPattern = regexp.MustCompile(`%[^%]+%`)

func replaceFile(macro map[string]string, fname string) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fd.Close()
	replaceReader(macro, fd)
	return nil
}

func replaceReader(macro map[string]string, fd io.Reader) {
	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		text := sc.Text()
		text = rxPattern.ReplaceAllStringFunc(text, func(s string) string {
			name := s[1 : len(s)-1]
			if value, ok := macro[name]; ok {
				return value
			} else {
				return s
			}
		})
		fmt.Println(text)
	}
}

func mains(args []string) error {
	macro := make(map[string]string)
	fileCount := 0
	for _, arg := range args {
		pos := strings.IndexByte(arg, '=')
		if pos >= 0 {
			left := arg[0:pos]
			right := arg[pos+1:]
			macro[left] = right
		} else {
			if err := replaceFile(macro, arg); err != nil {
				return err
			}
			fileCount++
		}
	}
	if fileCount <= 0 {
		replaceReader(macro, os.Stdin)
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
