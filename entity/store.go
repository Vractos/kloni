package entity

type Store struct {
	ID              ID
	Email           string
	Name            string
	MeliCredentials struct {
		AccessToken  string
		TokenType    string
		ExpiresIn    int
		Scope        string
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
