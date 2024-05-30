package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type structTest struct {
	Name  string `validate:"required,min=5,max=10"`
	Age   int    `validate:"required,min=12,max=18"`
	Email string `validate:"required,email"`
}

// go test -run TestValidatorWithoutDatas
func TestValidatorWithoutDatas(t *testing.T) {
	element := &structTest{}

	err := StructValidator.Validate(element)
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &ErrValidator)
}

// go test -run TestValidatorWitInvalidDatas
func TestValidatorWitInvalidDatas(t *testing.T) {
	element := &structTest{"1234", 22, "toError@email"}

	err := StructValidator.Validate(element)
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &ErrValidator)

	element.Name = "01234567890"
	element.Age = 33
	element.Email = "email.com"

	err = StructValidator.Validate(element)
	assert.NotNil(t, err)
	assert.IsType(t, "", err.Error())
	assert.ErrorAs(t, err, &ErrValidator)
}

// go test -run TestValidatorWitValidDatas
func TestValidatorWitValidDatas(t *testing.T) {
	element := &structTest{"123456", 15, "example@example.com"}

	err := StructValidator.Validate(element)
	assert.Nil(t, err)

	element.Name = "12345678"
	element.Age = 13
	element.Email = "email@email.com"

	err = StructValidator.Validate(element)
	assert.Nil(t, err)
}
