package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddError(t *testing.T) {
	t.Log("test add error")
	err := InvalidParams()
	err.AddError(InvalidField("id", "empty", "empty id"))
	err.AddError(InvalidParams())
	assert.Equal(t, true, err.HasError(), "")
}
