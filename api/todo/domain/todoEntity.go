package todoDomains

type Todo struct {
	Id          uint   `json:"id"`
	Username    string `json:"-"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
