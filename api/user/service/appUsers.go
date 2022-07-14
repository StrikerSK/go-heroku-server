package userServices

import (
	"go-heroku-server/api/types"
	userDomains "go-heroku-server/api/user/domain"
)

var admin = userDomains.User{
	FirstName: "admin",
	LastName:  "admin",
	Role:      userDomains.AdminRole,
	Address: types.Address{
		Street: "Admin",
		City:   "Admin",
		Zip:    "Admin",
	},
	Credentials: userDomains.Credentials{
		Username: "admin",
		Password: "admin",
	},
}

var guestUser = userDomains.User{
	FirstName: "tester",
	LastName:  "tester",
	Role:      userDomains.UserRole,
	Address: types.Address{
		Street: "Tester",
		City:   "Tester",
		Zip:    "Tester",
	},
	Credentials: userDomains.Credentials{
		Username: "admin",
		Password: "admin",
	},
}

func (s UserService) AddDemoUsers() {
	_ = s.repository.CreateUser(admin)
	_ = s.repository.CreateUser(guestUser)
}
