// Code generated by "stringer -tags dev -type OutputStyle"; DO NOT EDIT.

package libsass

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NestedStyle-0]
	_ = x[ExpandedStyle-1]
	_ = x[CompactStyle-2]
	_ = x[CompressedStyle-3]
}

const _OutputStyle_name = "NestedStyleExpandedStyleCompactStyleCompressedStyle"

var _OutputStyle_index = [...]uint8{0, 11, 24, 36, 51}

func (i OutputStyle) String() string {
	if i < 0 || i >= OutputStyle(len(_OutputStyle_index)-1) {
		return "OutputStyle(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OutputStyle_name[_OutputStyle_index[i]:_OutputStyle_index[i+1]]
}
