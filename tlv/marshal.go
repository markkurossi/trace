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

	BO.PutUint32(buf[:], uint32(i))
	(*bytes.Buffer)(enc).Write(buf[:])
}

func (enc *encoder) writeInt64(i int64) {
	var buf [8]byte

	BO.PutUint64(buf[:], uint64(i))
	(*bytes.Buffer)(enc).Write(buf[:])
}

func (enc *encoder) writeString(s string) {
	data := []byte(s)

	enc.writeInt32(int32(len(data)))
	(*bytes.Buffer)(enc).Write(data)
}

func (enc *encoder) writeAttr(a slog.Attr) error {
	// Resolve attributes.
	a.Value = a.Value.Resolve()

	// Ignore empty attributes.
	if a.Equal(slog.Attr{}) {
		return nil
	}
	switch a.Value.Kind() {
	case slog.KindInt64:
		enc.writeTag(TypeAttr, VTypeInt64)
		enc.writeString(a.Key)
		enc.writeInt64(a.Value.Int64())

	case slog.KindUint64:
		enc.writeTag(TypeAttr, VTypeUint64)
		enc.writeString(a.Key)
		enc.writeInt64(int64(a.Value.Uint64()))

	case slog.KindGroup:
		attrs := a.Value.Group()
		enc.writeTag(TypeAttr, VTypeGroup)
		enc.writeString(a.Key)
		enc.writeInt32(int32(len(attrs)))

		for i := 0; i < len(attrs); i++ {
			if err := enc.writeAttr(attrs[i]); err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("a.Value.Kind=%v not implemented yet", a.Value.Kind())
	}
	return nil
}

func (enc *encoder) Bytes() []byte {
	return (*bytes.Buffer)(enc).Bytes()
}

// Marshal encodes a slog.Record into binary data.
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

	var err error
	r.Attrs(func(a slog.Attr) bool {
		err = enc.writeAttr(a)
		if err != nil {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return enc.Bytes(), nil
}
