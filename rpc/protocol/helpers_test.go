package rpcproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestI2Uint32(t *testing.T) {
	assert.Equal(t, uint32(42), I2Uint32(float64(42)))
}
