package genmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RandStringBytesRmndr_ReturnsRandomString(t *testing.T) {
	t.Parallel()

	// Arrange
	n := 10

	// Act
	result := RandStringBytesRmndr(n)

	// Assert
	assert.NotEmpty(t, result)
	assert.Len(t, result, n)
}
