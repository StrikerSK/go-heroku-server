package types

type User struct {
	ID        uint    `json:"-"`
	Username  string  `json:"username"`
	Password  string  `json:"-"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Address   Address `json:"address"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
