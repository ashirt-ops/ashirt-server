package isthere

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNo(t *testing.T) {
	assert.True(t, No(nil))
	assert.False(t, No(fmt.Errorf("yep")))
}
