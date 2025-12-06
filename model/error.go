package model

// ErrNotFound is returned when a requested resource is not found.
type ErrNotFound struct{}

// Error implements the error interface for ErrNotFound.
func (ErrNotFound) Error() string {
    return "not found"
}
