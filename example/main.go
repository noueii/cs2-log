package main

import (
	"bufio"
	"fmt"
	"os"

	cs2log "github.com/noueii/cs2-log"
)

// Usage:
//
// From file:
// go run main.go ../debug/good_logs/sv1/combined_logs.log
//
// From STDIN:
// cat ../debug/good_logs/sv1/combined_logs.log | go run main.go
//
// To File:
// go run main.go ../debug/good_logs/sv1/combined_logs.log > out.txt
//
// Omit errors:
// go run main.go ../debug/good_logs/sv1/combined_logs.log 2>/dev/null

func main() {

	var file *os.File
	var err error

	if len(os.Args) < 2 {
		file = os.Stdin
	} else {
		file, err = os.Open(os.Args[1])
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := bufio.NewReader(file)

	// read first line
	l, _, err := r.ReadLine()

	for err == nil {

		// parse
		m, errParse := cs2log.Parse(string(l))

		if errParse != nil {
			// print parse errors to stderr
			fmt.Fprintf(os.Stderr, "ERROR: %s", cs2log.ToJSON(m))
		} else {
			// print to stdout
			fmt.Fprintf(os.Stdout, "%s", cs2log.ToJSON(m))
		}

		// next line
		l, _, err = r.ReadLine()
	}
}
