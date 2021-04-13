package domain

type SessionUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewSessionUser(id, name, picture string) SessionUser {
	return SessionUser{id, name, picture}
}
