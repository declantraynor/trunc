package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomStringGenerator(t *testing.T) {
	outputs := []string{}
	for i := 0; i < 1000; i++ {
		s := RandomString(8)
		if ok := assert.NotContains(t, outputs, s); !ok {
			return
		}
		outputs = append(outputs, s)
	}
}
