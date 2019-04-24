package types

type Address struct {
	ID     uint   `json:"-"`
	UserID uint   `json:"-"`
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip"`
}
