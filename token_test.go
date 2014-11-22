package repass

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsExpire(t *testing.T) {
	token := Token{
		ID:     "123",
		Email:  "test@test.com",
		Expire: time.Date(2000, 10, 10, 10, 10, 0, 0, time.UTC),
	}

	expire := token.IsExpire()
	assert.True(t, expire)
}
