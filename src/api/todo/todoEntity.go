package todo

import (
	"errors"
)

type Todo struct {
	Id          uint   `json:"id"`
	UserID      uint   `json:"-"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (t *Todo) validateAccess(userID uint) error {
	if t.UserID != userID {
		return errors.New("access denied")
	} else {
		return nil
	}
}
