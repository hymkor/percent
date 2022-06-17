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
	return replaceReader(macro, fd)
}

func replaceReader(macro map[string]string, fd io.Reader) error {
	br := bufio.NewReader(fd)
	for {
		text, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		text = rxPattern.ReplaceAllStringFunc(text, func(s string) string {
			name := s[1 : len(s)-1]
			if value, ok := macro[name]; ok {
				return value
			} else {
				return s
			}
		})
		io.WriteString(os.Stdout, text)
		if err == io.EOF {
			return nil
		}
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
		if err := replaceReader(macro, os.Stdin); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
