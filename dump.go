package main

import (
	"encoding/hex"
	"fmt"
	"io"
)

var dmp *dumper

func newDumpler(w io.Writer) *dumper {
	return &dumper{
		writer: w,
	}
}

type dumper struct {
	writer io.Writer
}

func (d *dumper) Write(p []byte, text string) (n int, err error) {
	fmt.Fprintf(d.writer, "\n%s\n================================\n", text)
	dd := hex.Dumper(d.writer)
	defer dd.Close()
	return dd.Write(p)
}
