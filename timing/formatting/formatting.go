package formatting

import (
	"fmt"
)

func TimeFormatMilliseconds(milliseconds int64) (out string) {
	// minutes, seconds, milliseconds
	out = fmt.Sprintf("%02d:%02d.%03d", milliseconds/60000%60, milliseconds/1000%60, milliseconds%1000)
	if milliseconds >= 3600000 {
		// prepend hours if necessary
		out = fmt.Sprintf("%02d:%s", milliseconds/3600000, out)
	}
	return out
}

func DeltaFormatMilliseconds(milliseconds int64) (out string) {
	sign := (int64)(1)
	if milliseconds < 0 {
		sign = -1
	}
	milliseconds *= sign

	if milliseconds >= 3600000 {
		// hours (single digit), minutes, seconds
		out = fmt.Sprintf("%d:%02d:%02d", sign*milliseconds/3600000, milliseconds/60000%60, milliseconds/1000%60)
	} else if milliseconds >= 60000 {
		// minutes (single digit), seconds
		out = fmt.Sprintf("%d:%02d", sign*milliseconds/60000%60, milliseconds/1000%60)
	} else {
		// seconds (always present, but can only be single digit in this case)
		out = fmt.Sprintf("%d", sign*milliseconds/1000%60)
	}
	// milliseconds (always present)
	out += fmt.Sprintf(".%03d", milliseconds%1000)

	return out
}
