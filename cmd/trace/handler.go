//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"log/slog"

	"github.com/markkurossi/scheme"
)

func handleRecord(r slog.Record) error {
	var rec scheme.Value

	rec = scheme.NewPair(scheme.NewPair(
		scheme.String("message"),
		scheme.String(r.Message)),
		rec)
	rec = scheme.NewPair(scheme.NewPair(
		scheme.String("level"),
		scheme.Int(r.Level)),
		rec)

	if !r.Time.IsZero() {
		rec = scheme.NewPair(scheme.NewPair(
			scheme.String("time"),
			scheme.Int(r.Time.UnixNano())),
			rec)
	}

	var attrs scheme.Value

	empty := slog.Attr{}

	r.Attrs(func(a slog.Attr) bool {
		a.Value = a.Value.Resolve()
		if a.Equal(empty) {
			return true
		}
		switch a.Value.Kind() {
		case slog.KindInt64:
			attrs = scheme.NewPair(scheme.NewPair(
				scheme.String(a.Key),
				scheme.Int(a.Value.Int64())),
				attrs)
		}
		return true
	})

	rec = scheme.NewPair(scheme.NewPair(scheme.String("attrs"), attrs), rec)

	_, err := scm.Apply(scmHandleRecord, []scheme.Value{
		rec,
	})

	return err
}
