// Code generated by "stringer -type=VType -trimprefix=VType"; DO NOT EDIT.

package tlv

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[VTypeAny-0]
	_ = x[VTypeBool-1]
	_ = x[VTypeDuration-2]
	_ = x[VTypeFloat64-3]
	_ = x[VTypeGroup-4]
	_ = x[VTypeInt32-5]
	_ = x[VTypeInt64-6]
	_ = x[VTypeString-7]
	_ = x[VTypeTime-8]
	_ = x[VTypeUint64-9]
}

const _VType_name = "AnyBoolDurationFloat64GroupInt32Int64StringTimeUint64"

var _VType_index = [...]uint8{0, 3, 7, 15, 22, 27, 32, 37, 43, 47, 53}

func (i VType) String() string {
	if i >= VType(len(_VType_index)-1) {
		return "VType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _VType_name[_VType_index[i]:_VType_index[i+1]]
}
