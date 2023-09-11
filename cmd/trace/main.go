//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/markkurossi/scheme"
	"github.com/markkurossi/trace/tlv"
)

const (
	path = "/tmp/trace.sock"
)

var (
	scm             *scheme.Scheme
	scmHandleRecord scheme.Value
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

	os.RemoveAll(path)
	listener, err := net.Listen("unix", path)
	if err != nil {
		log.Fatalf("failed to create listener: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept failed: %s", err)
			continue
		}
		log.Printf("new connection")
		go func(conn net.Conn) {
			var buf [1024]byte
			for {
				// XXX length + data
				n, err := conn.Read(buf[:])
				if err != nil {
					fmt.Printf("read failed: %s\n", err)
					break
				}
				fmt.Printf("Data:\n%s", hex.Dump(buf[:n]))

				r, err := tlv.Unmarshal(buf[:n])
				if err != nil {
					fmt.Printf("Unmarshal failed: %v\n", err)
					break
				}
				err = handleRecord(r)
				if err != nil {
					fmt.Printf("handleRecord: %v\n", err)
					break
				}
			}
			conn.Close()
		}(conn)
	}
}
