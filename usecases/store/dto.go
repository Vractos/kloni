package store

type RegisterStoreDtoInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type RegisterMeliCredentialsDtoInput struct {
	Code string `json:"code"`
	User string `json:"user"`
}
