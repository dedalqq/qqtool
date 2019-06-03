package main

import (
	"crypto/tls"
	"flag"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	host string

	listening      string
	forward        string
	useTLS         bool
	folderPath     string
	retry          bool
	fileServerPath string
	execCommand    string
	verbose        bool
)

func init() {
	flag.StringVar(&listening, "l", "", "Listening TCP port. Example: :80 or 0.0.0.0:80")
	flag.StringVar(&forward, "f", "", "Forward incoming connection to host. (with -l)")
	flag.BoolVar(&useTLS, "t", false, "Connect with TLS. (With -f or simple host connection)")
	flag.StringVar(&folderPath, "s", "", "Path to folder for save data from incoming connections (with -l or simple host connection)")
	flag.BoolVar(&retry, "r", false, "Repeat the connection after disconnecting (simple host connection)")
	flag.StringVar(&fileServerPath, "F", "", "Run file server on path (with -l)")
	flag.StringVar(&execCommand, "e", "", "Run command for incomming connections (with -l)")
	flag.BoolVar(&verbose, "v", false, "Make the operation more talkative")

	flag.Parse()

	if len(flag.Args()) > 0 {
		host = flag.Args()[0]
	}
}

func makeTLS(conn net.Conn) (net.Conn, error) {

	addr := conn.RemoteAddr()

	c := tls.Client(conn, &tls.Config{
		ServerName:         strings.Split(addr.String(), ":")[0],
		InsecureSkipVerify: true,
	})

	err := c.Handshake()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getFileForOut(folderPath, name string) (*os.File, error) {

	// if folderPath == ":stdout" {
	// 	return os.Stdout, nil
	// }

	// if folderPath == ":stderr" {
	// 	return os.Stderr, nil
	// }

	if folderPath != "" {
		f, err := openNewFile(folderPath, name)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	return nil, nil
}

func listeningTCP(addr string) {
	logState("Start listening TCP PORT [%v]", addr)

	server, err := getListener(addr, useTLS)
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

		go func(conn net.Conn) {
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
				defer c.Close()

				f, err := getFileForOut(folderPath, conn.RemoteAddr().String())
				if err != nil {
					logError(err)
					return
				}

				logSuccess()

				bind(c, c, conn, f)
				logInfo("Close connection from [%v]", conn.RemoteAddr())
				return
			}

			if execCommand != "" {
				logState("Run command [%s] for incoming connection [%s]", execCommand, conn.RemoteAddr())
				cmd := exec.Command(execCommand)

				stdin, err := cmd.StdinPipe()
				if nil != err {
					logError(err)
					return
				}

				stdout, err := cmd.StdoutPipe()
				if nil != err {
					logError(err)
					return
				}

				stderr, err := cmd.StderrPipe()
				if nil != err {
					logError(err)
					return
				}

				err = cmd.Start()
				if nil != err {
					logError(err)
					return
				}

				f, err := getFileForOut(folderPath, conn.RemoteAddr().String())
				if err != nil {
					logError(err)
					return
				}

				logSuccess()

				bind(io.MultiReader(stdout, stderr), stdin, conn, f)

				return
			}

			logInfo("Accept new connection from [%v]", conn.RemoteAddr())
			io.Copy(os.Stdout, conn)
		}(conn)
	}
}

func connectToHost(host string) {
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

	var f *os.File
	if folderPath != "" {
		f, err = openNewFile(folderPath, conn.RemoteAddr().String())
		if f != nil {
			defer f.Close()
		}
		if err != nil {
			logError(err)
		}
	}

	bind(os.Stdin, os.Stdout, conn, f)
	logInfo("Disconnect [%v]", conn.RemoteAddr())
}

func runFileServer(addr, path string) {
	logInfo("Run file server on [%s] in dir: [%s]", addr, path)

	listener, err := getListener(addr, useTLS)
	if err != nil {
		logError(err)
		return
	}

	server := http.Server{
		Addr:    host,
		Handler: http.FileServer(http.Dir(path)),
	}

	server.Serve(listener)
}

func main() {

	if verbose {
		dmp = newDumpler(os.Stderr)
	}

	if listening != "" {
		if fileServerPath != "" {
			runFileServer(listening, fileServerPath)
			return
		}

		listeningTCP(listening)
		return
	}

	logState("Connect to [%s]", host)
	if retry {
		for {
			connectToHost(host)
			logState("Reconnect to [%s]", host)
		}
	} else {
		connectToHost(host)
	}
}
