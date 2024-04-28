package u_conv

import "strconv"

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}
