package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/nyaosorg/go-windows-mbcs"
)

var rxPattern = regexp.MustCompile(`%[^%]+%`)

func replaceFile(macro map[string][]byte, fname string) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fd.Close()
	return replaceReader(macro, fd)
}

func replaceReader(macro map[string][]byte, fd io.Reader) error {
	br := bufio.NewReader(fd)
	for {
		line, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		line = rxPattern.ReplaceAllFunc(line, func(s []byte) []byte {
			name := string(s[1 : len(s)-1])
			if value, ok := macro[name]; ok {
				return value
			} else {
				return s
			}
		})
		os.Stdout.Write(line)
		if err == io.EOF {
			return nil
		}
	}
}

var flagAnsi = flag.Bool("ansi", false, "macro value is not UTF8 (ANSI)")

func mains(args []string) error {
	macro := make(map[string][]byte)
	fileCount := 0
	for _, arg := range args {
		pos := strings.IndexByte(arg, '=')
		if pos >= 0 {
			left := arg[0:pos]
			right := arg[pos+1:]
			if *flagAnsi {
				var err error
				macro[left], err = mbcs.UtoA(right, mbcs.ACP)
				if err != nil {
					return fmt.Errorf("%s: %w", arg, err)
				}
			} else {
				macro[left] = []byte(right)
			}
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
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
