package errors

type DatabaseError struct {
	Message string
}

func NewDatabaseError(message string) DatabaseError {
	return DatabaseError{
		Message: message,
	}
}

func (e DatabaseError) Error() string {
	return e.Message
}
