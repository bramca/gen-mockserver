package genmock

import (
	"net"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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

func Test_generateExampleData_ReturnsCorrectStringFormat(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		propSchema    *base.Schema
		formatChecker func(string) bool
	}{
		"enum": {
			propSchema: &base.Schema{
				Enum: []*yaml.Node{
					{
						Value: "test 1",
					},
					{
						Value: "test 2",
					},
				},
			},
			formatChecker: func(s string) bool {
				return slices.Contains([]string{"test 1", "test 2"}, s)
			},
		},
		"date-time": {
			propSchema: &base.Schema{
				Format: "date-time",
			},
			formatChecker: func(s string) bool {
				_, err := time.Parse("2006-01-02 15:04:05.000000000 +0000 UTC", s)

				return assert.NoError(t, err)
			},
		},
		"uuid": {
			propSchema: &base.Schema{
				Format: "uuid",
			},
			formatChecker: func(s string) bool {
				_, err := uuid.Parse(s)

				return assert.NoError(t, err)
			},
		},
		"ip": {
			propSchema: &base.Schema{
				Format: "ip",
			},
			formatChecker: func(s string) bool {
				ip := net.ParseIP(s)

				return assert.NotNil(t, ip)
			},
		},
		"ip-cidr-block": {
			propSchema: &base.Schema{
				Format: "ip-cidr-block",
			},
			formatChecker: func(s string) bool {
				_, _, err := net.ParseCIDR(s)

				return assert.NoError(t, err)
			},
		},
		"mac-address": {
			propSchema: &base.Schema{
				Format: "mac-address",
			},
			formatChecker: func(s string) bool {
				_, err := net.ParseMAC(s)

				return assert.NoError(t, err)
			},
		},
		"address-or-block-or-range": {
			propSchema: &base.Schema{
				Format: "address-or-block-or-range",
			},
			formatChecker: func(s string) bool {
				ip := net.ParseIP(s)

				return assert.NotNil(t, ip)
			},
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Act
			result := generateExampleData(data.propSchema)

			// Assert
			assert.NotEmpty(t, result)
			assert.True(t, data.formatChecker(result))
		})
	}
}
