//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package trace

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"

	"github.com/markkurossi/trace/tlv"
)

var (
	_ slog.Handler = &Client{}
)

// Client implements trace client and slog.Handler.
type Client struct {
	conn  net.Conn
	level slog.Level
}

// NewClient creates a new trace client.
func NewClient(path string) (*Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

// Enabled implements slog.Handler.Enabled.
func (c *Client) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= c.level
}

// Handle implements slog.Handler.Handle.
func (c *Client) Handle(ctx context.Context, r slog.Record) error {
	data, err := tlv.Marshal(r)
	if err != nil {
		fmt.Printf("tlv.Marshal failed: %s\n", err)
		return err
	}
	fmt.Printf("Data:\n%s", hex.Dump(data))

	l := len(data)
	var buf [4]byte
	tlv.BO.PutUint32(buf[:], uint32(l))

	_, err = c.conn.Write(buf[:])
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	return err
}

// WithAttrs implements slog.Handler.WithAttrs.
func (c *Client) WithAttrs(attrs []slog.Attr) slog.Handler {
	return c
}

// WithGroup implements slog.Handler.WithGroup.
func (c *Client) WithGroup(name string) slog.Handler {
	return c
}
