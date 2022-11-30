package entity

type Store struct {
	ID    ID
	Email string
	Name  string
}

func NewStore(email, name string) (*Store, error) {
	store := Store{
		ID:    NewID(),
		Email: email,
		Name:  name,
	}
	return &store, nil
}

func NewStoreWithKnownId(id ID, name, email string) (*Store, error) {
	panic("unimplemented")
}
