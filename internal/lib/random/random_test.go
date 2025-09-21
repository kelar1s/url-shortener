package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	testCases := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 50",
			size: 50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			string1 := NewRandomString(tc.size)
			string2 := NewRandomString(tc.size)

			assert.Len(t, string1, tc.size)
			assert.Len(t, string2, tc.size)

			assert.NotEqual(t, string1, string2)
		})
	}
}
