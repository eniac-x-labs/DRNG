package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateNumber(t *testing.T) {
	n, err := GenerateRandomNum()
	assert.NoError(t, err)
	t.Log(n)
}
