package helpers

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testFactor = time.Second

// go test -run TestDurationWithEmptyString
func TestDurationWithEmptyString(t *testing.T) {
	life, err := DurationFromString("", testFactor)

	assert.Error(t, err, "Did not return error!")
	assert.Equal(t, time.Duration(0), life)
	assert.Equal(t, err, ErrIntConvert)
}

// go test -run TestDurationWithInvalidString
func TestDurationWithInvalidString(t *testing.T) {
	life, err := DurationFromString("1m4i", testFactor)

	assert.Error(t, err, "Did not return error!")
	assert.Equal(t, time.Duration(0), life)
	assert.Equal(t, err, ErrIntConvert)
}

// go test -run TestDurationFromString
func TestDurationFromString(t *testing.T) {
	for i := 0; i < 99999; i++ {
		for _, factor := range []time.Duration{time.Nanosecond, time.Microsecond, time.Millisecond, time.Second, time.Minute, time.Hour} {
			life, err := DurationFromString(fmt.Sprint(i+1), factor)

			assert.Nil(t, err, "Returned error!")
			assert.Equal(t, time.Duration(i+1)*factor, life)
		}
	}
}
