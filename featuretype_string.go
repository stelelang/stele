// Code generated by "stringer -type FeatureType"; DO NOT EDIT.

package stele

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InvalidFeature-0]
	_ = x[LetFeature-1]
	_ = x[FuncFeature-2]
	_ = x[MemLayoutFeature-3]
}

const _FeatureType_name = "InvalidFeatureLetFeatureFuncFeatureMemLayoutFeature"

var _FeatureType_index = [...]uint8{0, 14, 24, 35, 51}

func (i FeatureType) String() string {
	if i < 0 || i >= FeatureType(len(_FeatureType_index)-1) {
		return "FeatureType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FeatureType_name[_FeatureType_index[i]:_FeatureType_index[i+1]]
}