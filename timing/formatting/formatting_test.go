package formatting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeFormat(t *testing.T) {
	assert.Equal(t, TimeFormatMilliseconds(200000), "03:20.000",
		"If input less than an hour, don't show hours")
	assert.Equal(t, TimeFormatMilliseconds(4000000), "01:06:40.000",
		"If input more than an hour, show hours")
}

func TestDeltaFormat(t *testing.T) {
	assert.Equal(t, DeltaFormatMilliseconds(-100), "-0.100",
		"Case: (negative) seconds")
	assert.Equal(t, DeltaFormatMilliseconds(-100000), "-1:40.000",
		"Case: (negative) minutes")
	assert.Equal(t, DeltaFormatMilliseconds(-4000000),
		"-1:06:40.000", "Case: (negative) hours")
	assert.Equal(t, DeltaFormatMilliseconds(100000), "+1:40.000",
		"Case: positive (minutes)")
	assert.Equal(t, DeltaFormatMilliseconds(0), "=0.000",
		"Case: millisecond tie")
}
