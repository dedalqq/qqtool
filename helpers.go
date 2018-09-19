package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func bind(r io.Reader, w io.Writer, rw io.ReadWriter, f *os.File) {
	var copy = func(dst io.Writer, src io.Reader, f *os.File, ech chan error) int64 {
		buf := make([]byte, 4*1024)
		var written int64
		for {
			nr, er := src.Read(buf)
			if nr > 0 {
				if f != nil {
					_, ewf := f.Write(buf[0:nr])
					if ewf != nil {
						f = nil
					}
				}
				nw, ew := dst.Write(buf[0:nr])
				if nw > 0 {
					written += int64(nw)
				}
				if ew != nil {
					ech <- ew
					return written
				}
				if nr != nw {
					ech <- errors.New("short write")
					return written
				}
			}
			if er != nil {
				if er == io.EOF {
					ech <- nil
					return written
				}
				ech <- er
				return written
			}
		}
	}

	ech := make(chan error)
	go copy(w, rw, f, ech)
	go copy(rw, r, f, ech)
	<-ech
}

func openNewFile(path string, v ...string) (*os.File, error) {
	f, err := os.Create(filepath.Join(path, fmt.Sprintf("dump_%s.txt", strings.Join(v, "_"))))
	if err != nil {
		return nil, err
	}
	return f, nil
}
