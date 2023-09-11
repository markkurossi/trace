// Code generated by "stringer -type=Type -trimprefix=Type"; DO NOT EDIT.

package tlv

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeTime-0]
	_ = x[TypeMessage-1]
	_ = x[TypeLevel-2]
	_ = x[TypeSource-3]
	_ = x[TypeAttr-4]
}

const _Type_name = "TimeMessageLevelSourceAttr"

var _Type_index = [...]uint8{0, 4, 11, 16, 22, 26}

func (i Type) String() string {
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
