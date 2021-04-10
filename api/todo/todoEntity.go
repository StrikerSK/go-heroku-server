package todo

type Todo struct {
	Id          uint   `json:"id"`
	UserID      uint   `json:"-"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type UserNumber struct {
	Id     uint   `json:"id"`
	Number uint   `json:"userNumber"`
	Todos  []Todo `json:"todos" gorm:"foreignkey:StudentId"`
}
