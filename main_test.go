package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert.Equal(t, 1+1, 2)
}

func TestSubtract(t *testing.T) {
	assert.Equal(t, 1-1, 0)
}
