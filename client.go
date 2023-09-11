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

type Client struct {
	conn  net.Conn
	level slog.Level
}

func NewClient(path string) (*Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= c.level
}

func (c *Client) Handle(ctx context.Context, r slog.Record) error {
	data, err := tlv.Marshal(r)
	if err != nil {
		fmt.Printf("tlv.Marshal failed: %s\n", err)
		return err
	}
	fmt.Printf("Data:\n%s", hex.Dump(data))

	n, err := c.conn.Write(data)
	fmt.Printf("write: n=%v, err=%v\n", n, err)
	return err
}

func (c *Client) WithAttrs(attrs []slog.Attr) slog.Handler {
	return c
}

func (c *Client) WithGroup(name string) slog.Handler {
	return c
}
