package app

// App sentinel errors.
const (
	ErrBookNotFound = Error("book not found")
)

// Error is a type for app sentinel errors.
type Error string

func (e Error) Error() string {
	return string(e)
}
