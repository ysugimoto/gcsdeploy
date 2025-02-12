package remote

type RemoteError struct {
	err     error
	message string
}

func Error(err error, message string) *RemoteError {
	return &RemoteError{
		err:     err,
		message: message,
	}
}

func (e *RemoteError) Error() string {
	return e.message
}

func (e *RemoteError) Unwrap() error {
	return e.err
}
