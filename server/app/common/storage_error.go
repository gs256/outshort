package common

const (
	ErrorAny                = 0
	ErrorUniqueViolation    = 1
	ErrorNotFound           = 2
	ErrorInvalidCredentials = 3
)

type StorageError struct {
	Code    int
	message string
}

func NewStorageError(code int, message string) *StorageError {
	return &StorageError{code, message}
}

func (this *StorageError) Error() string {
	return this.message
}
