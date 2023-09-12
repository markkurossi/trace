//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"log/slog"

	"github.com/markkurossi/trace"
)

func main() {
	logger, err := trace.NewClient("/tmp/trace.sock")
	if err != nil {
		slog.Error("init", "error", err.Error())
		return
	}
	slog.SetDefault(slog.New(logger))

	slog.Info("test", "a", 1, "b", 2)
	slog.Info("step", "done", 1, slog.Group("gates",
		"xor", 10,
		"or", 20,
		"not", 7))
}
