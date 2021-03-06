package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func getListener(addr string, useTLS bool) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	if useTLS {
		c, err := getTLSConfig()
		if err != nil {
			return nil, err
		}
		listener = tls.NewListener(listener, c)
	}

	return listener, nil
}

func copy(dst io.Writer, src io.Reader, f *os.File, ech chan error, direction string) int64 {
	buf := make([]byte, 4*1024)
	var written int64
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			if dmp != nil {
				dmp.Write(buf[:nr], direction)
			}
			if f != nil {
				_, ewf := f.Write(buf[:nr])
				if ewf != nil {
					f = nil
				}
			}
			nw, ew := dst.Write(buf[:nr])
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

func bind(r io.Reader, w io.Writer, rw io.ReadWriter, f *os.File) {
	ech := make(chan error)
	go copy(w, rw, f, ech, "<<<")
	go copy(rw, r, f, ech, ">>>")
	<-ech
}

func openNewFile(path string, v ...string) (*os.File, error) {
	f, err := os.Create(filepath.Join(path, fmt.Sprintf("dump_%s.txt", strings.Join(v, "_"))))
	if err != nil {
		return nil, err
	}
	return f, nil
}
