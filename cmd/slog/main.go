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
	slog.Info("progress", slog.Int("done", 43), slog.Int("total", 100),
		slog.Group("gates",
			slog.Uint64("xor", 10),
			slog.Uint64("or", 20),
			slog.Uint64("not", 7)))
}
