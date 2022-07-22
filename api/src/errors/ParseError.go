package errors

type ParseError struct {
	Message string
}

func NewParseError(message string) ParseError {
	return ParseError{
		Message: message,
	}
}

func (e ParseError) Error() string {
	return e.Message
}
