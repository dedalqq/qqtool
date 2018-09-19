package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func bind(r io.Reader, w io.Writer, rw io.ReadWriter, f *os.File) {
	var copy = func(dst io.Writer, src io.Reader, f *os.File, wg *sync.WaitGroup) (int64, error) {
		defer wg.Done()
		buf := make([]byte, 32*1024)
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
					return written, ew
				}
				if nr != nw {
					return written, errors.New("short write")
				}
			}
			if er != nil {
				if er == io.EOF {
					return written, nil
				}
				return written, er
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go copy(w, rw, f, &wg)
	go copy(rw, r, f, &wg)
	wg.Wait()
}

func openNewFile(path string, v ...string) (*os.File, error) {
	f, err := os.Create(filepath.Join(path, fmt.Sprintf("dump_%s.txt", strings.Join(v, "_"))))
	if err != nil {
		return nil, err
	}
	return f, nil
}
