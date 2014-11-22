package repass

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoding(t *testing.T) {
	tokenID := encode("jack@jack.com")
	assert.NotEmpty(t, tokenID)
}
