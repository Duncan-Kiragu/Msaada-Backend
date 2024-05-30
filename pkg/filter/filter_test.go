package filter

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -run TestNewFilter
func TestNewFilter(t *testing.T) {
	_ = os.Setenv("API_DEFAULT_SORT", "updated_at")
	_ = os.Setenv("API_DEFAULT_ORDER", "desc")

	filter := NewFilter()

	assert.NotNil(t, filter)
	assert.Equal(t, "", filter.Search)
	assert.Equal(t, 0, filter.Limit)
	assert.Equal(t, 0, filter.Page)
	assert.Equal(t, os.Getenv("API_DEFAULT_SORT"), filter.Sort)
	assert.Equal(t, os.Getenv("API_DEFAULT_ORDER"), filter.Order)
}
