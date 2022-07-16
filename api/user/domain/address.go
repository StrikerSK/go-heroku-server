package userDomains

import "fmt"

type Address struct {
	ID       uint   `json:"-"`
	Username string `json:"-"`
	Street   string `json:"street"`
	City     string `json:"city"`
	Zip      string `json:"zip"`
}

func (a Address) toString() {
	fmt.Printf("ID: %d, UserID: %d, Street: %s, City: %s, Zip: %s\n", a.ID, a.Username, a.Street, a.City, a.Zip)
}
