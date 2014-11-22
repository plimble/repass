package repass

import (
	"github.com/plimble/crud"
)

type Store interface {
	crud.CRUD
	Get(tokenID string) (*Token, error)
}
