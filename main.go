package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type function struct {
	file  string
	line  int64
	locks []string
}

type stackTrace struct {
	goroutine int64      // goroutine number
	calls     []function // function calls (actual stack)
	hasLocks  bool
}

func main() {
	arg := os.Args[1]
	var file *os.File
	var err error
	if arg == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(arg)
		defer file.Close()
	}
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	scanner := bufio.NewScanner(file)

	goRegexp := regexp.MustCompile(`^goroutine ([0-9]+) `)
	lineRegexp := regexp.MustCompile(`([^ \t]*\.go):([0-9]+) `)
	var curr *stackTrace
	var routines []*stackTrace

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 && curr != nil {
			routines = append(routines, curr)
			continue
		}

		submatch := goRegexp.FindStringSubmatch(line)
		if len(submatch) > 0 {
			num, _ := strconv.ParseInt(submatch[1], 10, 64)
			curr = &stackTrace{goroutine: num}
			continue
		}

		submatch = lineRegexp.FindStringSubmatch(line)
		if len(submatch) > 0 {
			file := submatch[1]
			ln, _ := strconv.ParseInt(submatch[2], 10, 64)

			curr.calls = append(curr.calls, function{file, ln, nil})
		}
	}

	for _, r := range routines {
		fillTrace(r)
		if !r.hasLocks {
			continue
		}
		fmt.Println("goroutine", r.goroutine)
		for _, c := range r.calls {
			if strings.HasPrefix(c.file, "/usr/") {
				//continue
			}
			fmt.Println("  ", c.file, c.line)
			for _, l := range c.locks {
				fmt.Println("    ", l)
			}
		}
	}
}

func fillTrace(trace *stackTrace) {
	for i := range trace.calls {
		checkCall(&trace.calls[i])
		trace.hasLocks = trace.hasLocks || len(trace.calls[i].locks) > 0
	}
}

func checkCall(fn *function) error {
	lines, err := getLines(fn.file)
	if err != nil {
		return err
	}

	funcRegexp := regexp.MustCompile(`^func `)
	lockRegexp := regexp.MustCompile(`[L]ock\(`)
	for i := fn.line; i > 0; i-- {
		line := []byte(lines[i])

		if funcRegexp.Match(line) {
			break
		}

		if lockRegexp.Match(line) {
			fn.locks = append(fn.locks, fmt.Sprintf("%s:%d: %s", fn.file, i+1, line))
		}
	}

	return nil
}

func getLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
