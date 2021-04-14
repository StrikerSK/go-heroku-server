package types

import "fmt"

type Address struct {
	ID     uint   `json:"-"`
	UserID uint   `json:"-"`
	Street string `json:"street" example:"Test Street 999"`
	City   string `json:"city" example:"Test City"`
	Zip    string `json:"zip" example:"99999"`
}

func (a Address) toString() {
	fmt.Printf("ID: %d, UserID: %d, Street: %s, City: %s, Zip: %s\n", a.ID, a.UserID, a.Street, a.City, a.Zip)
}
