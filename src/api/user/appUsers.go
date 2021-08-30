package user

import (
	"go-heroku-server/src/api/types"
)

var admin = User{
	FirstName: "Ivan",
	LastName:  "Tester",
	Role:      AdminRole,
	Address: types.Address{
		Street: "Foo",
		City:   "Bar",
		Zip:    "12345",
	},
	Credentials: Credentials{
		Username: "admin",
		Password: "admin",
	},
}

var guestUser = User{
	FirstName: "Karel",
	LastName:  "Tester",
	Role:      GuestRole,
	Address: types.Address{
		Street: "Foo",
		City:   "Bar",
		Zip:    "98765",
	},
	Credentials: Credentials{
		Username: "guest",
		Password: "guest",
	},
}

func InitializeUsers() {
	_ = addUser(admin)
	_ = addUser(guestUser)
}
