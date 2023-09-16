//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package trace

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"

	"github.com/markkurossi/trace/tlv"
)

// Server implements trace server.
type Server struct {
	listener net.Listener
}

// NewServer creates a new trace server.
func NewServer(path string) (*Server, error) {
	os.RemoveAll(path)
	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: listener,
	}, nil
}

// Accept accepts a new client connection.
func (s *Server) Accept() (*Connection, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	c := &Connection{
		conn: conn,
		C:    make(chan msg),
	}
	go c.msgLoop()
	return c, nil
}

type msg struct {
	Err error
	R   slog.Record
}

// Connection implements a client connection.
type Connection struct {
	conn net.Conn
	C    chan msg
}

// Close closes the client connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) msgLoop() {
	err := c.processMessages()
	if err != nil {
		c.C <- msg{
			Err: err,
		}
	}
	c.conn.Close()
	close(c.C)
}

func (c *Connection) processMessages() error {
	var buf [5]byte

	// Read magic.
	_, err := c.conn.Read(buf[:4])
	if err != nil {
		return err
	}
	magic := tlv.BO.Uint32(buf[:4])
	if magic != Magic {
		return fmt.Errorf("invalid client magic: %08x", magic)
	}

	var data []byte

	for {
		_, err := c.conn.Read(buf[:5])
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		t := MsgType(buf[0])
		l := int(tlv.BO.Uint32(buf[1:5]))
		if l > len(data) {
			data = make([]byte, l)
		}
		_, err = c.conn.Read(data[:l])
		if err != nil {
			return err
		}
		switch t {
		case MsgSlog:
			r, err := tlv.Unmarshal(data[:l])
			if err != nil {
				return err
			}
			c.C <- msg{
				R: r,
			}

		default:
			return fmt.Errorf("unknown message: %v", t)
		}
	}
}
