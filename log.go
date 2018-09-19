package main

import (
	"fmt"
	"os"
)

var state bool

func logInfo(s string, i ...interface{}) {
	if state {
		fmt.Fprintf(os.Stderr, "[ unknown ]\n")
	}
	fmt.Fprintf(os.Stderr, " > %s\n", fmt.Sprintf(s, i...))
}

func logState(s string, i ...interface{}) {
	fmt.Fprintf(os.Stderr, " > %s ... ", fmt.Sprintf(s, i...))
	state = true
}

func logSuccess() {
	if state {
		fmt.Fprintf(os.Stderr, "[ success ]\n")
		state = false
	}
}

func logError(err error) {
	if state {
		fmt.Fprintf(os.Stderr, "[ error ] [%v]\n", err)
		state = false
	} else {
		fmt.Fprintf(os.Stderr, "> error [%v]\n", err)
	}
}
