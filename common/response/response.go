package response

type Success struct {
	Message string `json:"message"`
}

type Error struct {
	Error string `json:"error"`
}

func NewSuccess(message string) *Success {
	return &Success{
		Message: message,
	}
}

func NewError(err error) *Error {
	return &Error{
		Error: err.Error(),
	}
}
