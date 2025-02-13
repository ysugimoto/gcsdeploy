package local

type ClientInterface interface {
	ListObjects() (Objects, error)
}
