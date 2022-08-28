package entity

type Store struct {
	ID              ID
	Email           string
	Name            string
	Password        string
	MeliCredentials struct {
		AccessToken  string
		TokenType    string
		ExpiresIn    int
		Scope        string
		UserID       int
		RefreshToken string
	}
}

func NewStore(email, name, password string) (*Store, error) {
	store := Store{
		ID:       NewID(),
		Email:    email,
		Password: password,
	}
	return &store, nil
}
