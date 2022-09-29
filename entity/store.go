package entity

type Store struct {
	ID              ID
	Email           string
	Name            string
	MeliCredentials struct {
		AccessToken  string
		ExpiresIn    int
		UserID       int
		RefreshToken string
	}
}

func NewStore(email, name string) (*Store, error) {
	store := Store{
		ID:    NewID(),
		Email: email,
		Name:  name,
	}
	return &store, nil
}
