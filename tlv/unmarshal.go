//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package tlv

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"time"
)

type decoder bytes.Buffer

func (dec *decoder) readTag() (Type, VType, error) {
	b, err := (*bytes.Buffer)(dec).ReadByte()
	if err != nil {
		return 0, 0, err
	}
	return Type(b >> 5), VType(b & 0x1f), nil
}

func (dec *decoder) readInt32() (int32, error) {
	var buf [4]byte

	_, err := (*bytes.Buffer)(dec).Read(buf[:])
	if err != nil {
		return 0, err
	}
	return int32(bo.Uint32(buf[:])), nil
}

func (dec *decoder) readInt64() (int64, error) {
	var buf [8]byte

	_, err := (*bytes.Buffer)(dec).Read(buf[:])
	if err != nil {
		return 0, err
	}
	return int64(bo.Uint64(buf[:])), nil
}

func (dec *decoder) readString() (string, error) {
	l, err := dec.readInt32()
	if err != nil {
		return "", err
	}

	buf := make([]byte, l)
	_, err = (*bytes.Buffer)(dec).Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func Unmarshal(data []byte) (r slog.Record, err error) {
	dec := (*decoder)(bytes.NewBuffer(data))

	for {
		t, vt, err := dec.readTag()
		if err != nil {
			if err == io.EOF {
				return r, nil
			}
			return r, err
		}

		switch t {
		case TypeTime:
			i, err := dec.readInt64()
			if err != nil {
				return r, err
			}
			r.Time = time.Unix(0, i)

		case TypeMessage:
			str, err := dec.readString()
			if err != nil {
				return r, err
			}
			r.Message = str
			fmt.Printf("r.Message: %v\n", str)

		case TypeLevel:
			i, err := dec.readInt32()
			if err != nil {
				return r, err
			}
			r.Level = slog.Level(i)

		case TypeAttr:
			key, err := dec.readString()
			if err != nil {
				return r, err
			}
			switch vt {
			case VTypeInt64:
				i, err := dec.readInt64()
				if err != nil {
					return r, err
				}
				r.Add(slog.Int64(key, i))

			default:
				return r, fmt.Errorf("VType %v not implemented yet", vt)
			}

		default:
			return r, fmt.Errorf("Type %v not implemented yet", t)
		}
	}
}
