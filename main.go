package main

import (
	"strings"
	"sync"
	"io"
	"os"
	"net"
	"flag"
	"crypto/tls"
)

var (
	host string

	listening string
	forward string
	useTLS bool
)


func init() {
	flag.StringVar(&listening, "l", "", "Listening TCP port. Example: :80 or 0.0.0.0:80")
	flag.StringVar(&forward, "f", "", "Forward incoming connection to host. (with -l)")
	flag.BoolVar(&useTLS, "t", false, "Connect with TLS. (With -f or simple host connection)")

	flag.Parse()

	if len(flag.Args()) > 0 {
		host = flag.Args()[0]
	}
}

func makeTLS(conn net.Conn) (net.Conn, error) {

	addr := conn.RemoteAddr()

	c := tls.Client(conn, &tls.Config{
		ServerName: strings.Split(addr.String(), ":")[0],
		InsecureSkipVerify: true,
	})

	err := c.Handshake()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func bind(a, b io.ReadWriter) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		io.Copy(a, b)
		wg.Done()
	}()
	go func() {
		io.Copy(b, a)
		wg.Done()
	}()
	wg.Wait()
}

func listeningTCP(addr string) {
	logState("Start listening TCP PORT [%v]", addr)

	server, err := net.Listen("tcp", addr)
	if err != nil {
		logError(err)
		return
	}

	logSuccess()
	defer func() {
		err := server.Close()
		if err != nil {
			logError(err)
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			logError(err)
			break
		}

		go func() {
			defer conn.Close()

			if forward != "" {
				logState("Forwarding incoming connection [%v] to [%v]", conn.RemoteAddr(), forward)
				c, err := net.Dial("tcp", forward)
				if err != nil {
					logError(err)
					return
				}

				if useTLS {
					c, err = makeTLS(c)
					if err != nil {
						logError(err)
						return
					}
				}

				logSuccess()
				defer c.Close()
				
				bind(c, conn)
				logInfo("Close connection from [%v]", conn.RemoteAddr())
				return
			}

			logInfo("Accept new connection from [%v]", conn.RemoteAddr())
			io.Copy(os.Stdout, conn)
		}()
	}
}

func main() {

	if listening != "" {
		listeningTCP(listening)
		return
	}

	logState("Connect to [%s]", host)

	conn, err := net.Dial("tcp", host)
	if err != nil {
		logError(err)
		return
	}

	if useTLS {
		conn, err = makeTLS(conn)
		if err != nil {
			logError(err)
			return
		}
	}

	logSuccess()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		io.Copy(os.Stdout, conn)
		wg.Done()
	}()
	go func() {
		io.Copy(conn, os.Stdin)
		wg.Done()
	}()
	wg.Wait()
}