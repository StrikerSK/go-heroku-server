package user

import "go-heroku-server/api/types"

var admin = User{
	FirstName: "admin",
	LastName:  "admin",
	Role:      AdminRole,
	Address: types.Address{
		Street: "Admin",
		City:   "Admin",
		Zip:    "Admin",
	},
	Credentials: Credentials{
		Username: "admin",
		Password: "admin",
	},
}

var guestUser = User{
	FirstName: "tester",
	LastName:  "tester",
	Role:      UserRole,
	Address: types.Address{
		Street: "Tester",
		City:   "Tester",
		Zip:    "Tester",
	},
	Credentials: Credentials{
		Username: "admin",
		Password: "admin",
	},
}

func InitializeUsers() {
	_ = addUser(admin)
	_ = addUser(guestUser)
}
