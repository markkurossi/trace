//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package tlv

import (
	"bytes"
	"fmt"
	"log/slog"
)

type encoder bytes.Buffer

func (enc *encoder) writeTag(t Type, vt VType) {
	(*bytes.Buffer)(enc).WriteByte(byte(t<<5) | byte(vt&0x1f))
}

func (enc *encoder) writeInt32(i int32) {
	var buf [4]byte

	bo.PutUint32(buf[:], uint32(i))
	(*bytes.Buffer)(enc).Write(buf[:])
}

func (enc *encoder) writeInt64(i int64) {
	var buf [8]byte

	bo.PutUint64(buf[:], uint64(i))
	(*bytes.Buffer)(enc).Write(buf[:])
}

func (enc *encoder) writeString(s string) {
	data := []byte(s)

	enc.writeInt32(int32(len(data)))
	(*bytes.Buffer)(enc).Write(data)
}

func (enc *encoder) Bytes() []byte {
	return (*bytes.Buffer)(enc).Bytes()
}

func Marshal(r slog.Record) ([]byte, error) {
	enc := new(encoder)

	enc.writeTag(TypeMessage, VTypeString)
	enc.writeString(r.Message)

	enc.writeTag(TypeLevel, VTypeInt32)
	enc.writeInt32(int32(r.Level))

	if !r.Time.IsZero() {
		enc.writeTag(TypeTime, VTypeInt64)
		enc.writeInt64(r.Time.UnixNano())
	}

	empty := slog.Attr{}

	r.Attrs(func(a slog.Attr) bool {
		// Resolve attributes.
		a.Value = a.Value.Resolve()

		// Ignore empty attributes.
		if a.Equal(empty) {
			return true
		}
		switch a.Value.Kind() {
		case slog.KindInt64:
			enc.writeTag(TypeAttr, VTypeInt64)
			enc.writeString(a.Key)
			enc.writeInt64(a.Value.Int64())

		default:
			fmt.Printf("a.Value.Kind=%v not implemented yet\n", a.Value.Kind())
		}
		return true
	})

	return enc.Bytes(), nil
}
