package repass

import (
	"time"
)

type Token struct {
	ID     string    `json:"id" bson:"_id"`
	UserID string    `json:"user_id" bson:"user_id"`
	Email  string    `json:"email" bson:"email"`
	Expire time.Time `json:"expire" bson:"expire"`
	Used   bool      `json:"used bson:"used"`
}

func (t *Token) IsExpire() bool {
	if t.Expire.Before(time.Now().UTC()) {
		return true
	}

	return false
}
