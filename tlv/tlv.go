//
// tlv.go
//
// Copyright (c) 2019-2023 Markku Rossi
//
// All rights reserved.
//

package tlv

import (
	"encoding/binary"
)

var (
	// BO defines the byte order used to marshal records.
	BO = binary.BigEndian
)

// Type specifies value type.
//
//go:generate stringer -type=Type -trimprefix=Type
type Type uint8

// Serialization types.
const (
	TypeTime Type = iota
	TypeMessage
	TypeLevel
	TypeSource
	TypeAttr
)

// VType specifies attribute value types.
//
//go:generate stringer -type=VType -trimprefix=VType
type VType uint32

// Value types.
const (
	VTypeAny VType = iota
	VTypeBool
	VTypeDuration
	VTypeFloat64
	VTypeGroup
	VTypeInt32
	VTypeInt64
	VTypeString
	VTypeTime
	VTypeUint64
)
