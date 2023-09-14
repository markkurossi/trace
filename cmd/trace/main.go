//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/markkurossi/scheme"
	"github.com/markkurossi/trace"
)

const (
	path = "/tmp/trace.sock"
)

var (
	scm             *scheme.Scheme
	scmHandleRecord scheme.Value
	scmRedraw       scheme.Value

	cSignal = make(chan os.Signal, 1)
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	var err error

	scm, err = scheme.New()
	if err != nil {
		log.Fatalf("scheme.New(): %s", err)
	}

	for _, arg := range flag.Args() {
		v, err := scm.EvalFile(arg)
		if err != nil {
			log.Fatalf("failed to load init file '%s': %s", arg, err)
		}
		fmt.Printf("%v\n", v)
	}

	// Fetch Scheme callbacks.
	scmHandleRecord, err = scm.Global("handle-record")
	if err != nil {
		log.Fatalf("Init: %s", err)
	}
	scmRedraw, err = scm.Global("redraw")
	if err != nil {
		log.Fatalf("Init: %s", err)
	}

	terminalSize()
	go signalHandler()
	signal.Notify(cSignal, syscall.SIGWINCH)

	server, err := trace.NewServer(path)
	if err != nil {
		log.Fatalf("failed to create server: %s", err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("accept failed: %s", err)
			continue
		}
		log.Printf("new connection")
		go func(conn *trace.Connection) {
			for msg := range conn.C {
				if msg.Err != nil {
					fmt.Printf("error: %s\n", msg.Err)
					return
				}
				err := handleRecord(msg.R)
				if err != nil {
					fmt.Printf("handleRecord: %v\n", err)
					return
				}
			}
			conn.Close()
		}(conn)
	}
}

func terminalSize() {
	rows, cols, err := Size()
	if err != nil {
		fmt.Printf("failed to get terminal size: %s\n", err)
	} else {
		scm.SetGlobal("rows", scheme.Int(rows))
		scm.SetGlobal("cols", scheme.Int(cols))

		_, err = scm.Apply(scmRedraw, nil)
		if err != nil {
			fmt.Printf("redraw: %s\n", err)
		}
	}
}

func signalHandler() {
	for signal := range cSignal {
		switch signal {
		case syscall.SIGWINCH:
			terminalSize()
		default:
			fmt.Printf("signal: %v\n", signal)
		}
	}
}
