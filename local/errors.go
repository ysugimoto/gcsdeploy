package local

type LocalError struct {
	err     error
	message string
}

func Error(err error, message string) *LocalError {
	return &LocalError{
		err:     err,
		message: message,
	}
}

func (e *LocalError) Error() string {
	return e.message
}

func (e *LocalError) Unwrap() error {
	return e.err
}
