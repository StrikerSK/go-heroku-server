package userDomains

const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct {
	UserCredentials
	UserID    uint    `json:"-" gorm:"primaryKey"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Role      string  `json:"-" gorm:"default:user"`
	Address   Address `json:"address" gorm:"foreignKey:UserID"`
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetRole() string {
	return u.Role
}

func (u *User) SetRole() {
	if u.Role == "" {
		u.Role = UserRole
	}
}
