package todo

type Todo struct {
	Id          uint   `json:"id"`
	UserID      uint   `json:"-"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
