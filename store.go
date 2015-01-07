package repass

type Store interface {
	Insert(v interface{}) error
	Update(id string, v map[string]interface{}) error
	Get(tokenID string) (*Token, error)
}
