package userServices

import (
	userDomains "go-heroku-server/api/user/domain"
)

var admin = userDomains.User{
	FirstName: "admin",
	LastName:  "admin",
	Role:      userDomains.AdminRole,
	Address: userDomains.Address{
		Street: "Admin",
		City:   "Admin",
		Zip:    "Admin",
	},
	UserCredentials: userDomains.UserCredentials{
		Username: "admin",
		Password: "admin",
	},
}

var guestUser = userDomains.User{
	FirstName: "tester",
	LastName:  "tester",
	Role:      userDomains.UserRole,
	Address: userDomains.Address{
		Street: "Tester",
		City:   "Tester",
		Zip:    "Tester",
	},
	UserCredentials: userDomains.UserCredentials{
		Username: "admin",
		Password: "admin",
	},
}

func (s UserService) AddDemoUsers() {
	_ = s.repository.CreateUser(admin)
	_ = s.repository.CreateUser(guestUser)
}
