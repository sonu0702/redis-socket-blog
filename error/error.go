package error

//ApplicationError is the custom Error
type ApplicationError struct {
	Code    string `json:"code,omitempty"`
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}

//New returns a new Apple error
func New(code string, name string, message string) *ApplicationError {
	return &ApplicationError{
		Code:    code,
		Name:    name,
		Message: message,
	}
}

func (err ApplicationError) Error() string {
	return err.Message
}

//ErrorCode returns the error code
func (err ApplicationError) ErrorCode() string {
	return err.Code
}

//ErrorMessage returns the extra info contained by the error
func (err ApplicationError) ErrorMessage() interface{} {
	return err.Message
}
