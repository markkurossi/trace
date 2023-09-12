//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"log/slog"

	"github.com/markkurossi/scheme"
)

func newAttr(attrs scheme.Value, key string, value scheme.Value) scheme.Value {
	return scheme.NewPair(scheme.NewPair(scheme.String(key), value), attrs)
}

func attrToScheme(a slog.Attr) (scheme.Value, error) {
	switch a.Value.Kind() {
	case slog.KindInt64:
		return scheme.Int(a.Value.Int64()), nil

	case slog.KindGroup:
		var ga scheme.Value
		attrs := a.Value.Group()
		for i := 0; i < len(attrs); i++ {
			v, err := attrToScheme(attrs[i])
			if err != nil {
				return nil, err
			}
			ga = newAttr(ga, attrs[i].Key, v)
		}
		return ga, nil

	default:
		return nil, fmt.Errorf("Attr.Kind %s not supported", a.Value.Kind())
	}
}

func handleRecord(r slog.Record) error {
	var rec scheme.Value

	rec = newAttr(rec, "message", scheme.String(r.Message))
	rec = newAttr(rec, "level", scheme.Int(r.Level))

	if !r.Time.IsZero() {
		rec = newAttr(rec, "time", scheme.Int(r.Time.UnixNano()))
	}

	var attrs scheme.Value
	var err error
	r.Attrs(func(a slog.Attr) bool {
		a.Value = a.Value.Resolve()
		if a.Equal(slog.Attr{}) {
			return true
		}
		var value scheme.Value
		value, err = attrToScheme(a)
		if err != nil {
			return false
		}
		attrs = newAttr(attrs, a.Key, value)
		return true
	})
	if err != nil {
		return err
	}

	rec = newAttr(rec, "attrs", attrs)

	_, err = scm.Apply(scmHandleRecord, []scheme.Value{
		rec,
	})

	return err
}
