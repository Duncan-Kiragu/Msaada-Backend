package helpers

import (
	"errors"
	"strconv"
	"time"
)

var ErrIntConvert = errors.New("invalid string number")

func DurationFromString(str string, factor time.Duration) (time.Duration, error) {
	converted, err := strconv.Atoi(str)
	if err != nil {
		return time.Duration(0), ErrIntConvert
	}

	return time.Duration(converted) * factor, nil
}
