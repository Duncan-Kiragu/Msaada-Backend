package helpers

import (
	"errors"
	"fmt"
	"testing"
)

func restore() {
	if r := recover(); r != nil {
		fmt.Println("Recovered in f", r)
	}
}

// go test -run TestErrPanicIfErr
func TestErrPanicIfErr(t *testing.T) {
	defer restore()

	for i := 0; i < 99999; i++ {
		PanicIfErr(errors.New("test error"))
		t.Errorf("The code did not panic")
	}
}

// go test -run TestNilPanicIfErr
func TestNilPanicIfErr(t *testing.T) {
	defer restore()

	for i := 0; i < 99999; i++ {
		PanicIfErr(nil)
	}
}
