package repass

import (
	"github.com/plimble/crud"
)

//go:generate mockgen -destination=mock_store.go --self_package=github.com/plimble/repass -package=repass github.com/plimble/repass Store

type Store interface {
	crud.CRUD
	Get(tokenID string) (*Token, error)
}
