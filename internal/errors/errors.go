package errors

// New is a convenience type for creating a const error using familiar syntax
type New string

// Error returns the error string
func (e New) Error() string {
	return string(e)
}
