package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	src := os.Stdin
	dst := os.Stdout

	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Connection error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Fprintf(dst, "Connected to %s\n", address)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer cancel()

		reader := bufio.NewReader(src)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Connection closed by client")
				return
			}

			_, err = conn.Write([]byte(line))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Write error:", err)
				return
			}
		}
	}()

	go func() {
		defer cancel()

		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Connection closed by server")
				return
			}

			fmt.Fprint(dst, line)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Fprintln(dst)
		cancel()
	}()
	<-ctx.Done()
	conn.Close()
}
