// Code generated by "stringer -type=funcEnum"; DO NOT EDIT.

package handler

import "strconv"

const _funcEnum_name = "BFTBlockAnnounceBlockRequestBlockResponseRequestStateRequestStatusTransactions"

var _funcEnum_index = [...]uint8{0, 3, 16, 28, 41, 48, 60, 66, 78}

func (i funcEnum) String() string {
	if i < 0 || i >= funcEnum(len(_funcEnum_index)-1) {
		return "funcEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _funcEnum_name[_funcEnum_index[i]:_funcEnum_index[i+1]]
}