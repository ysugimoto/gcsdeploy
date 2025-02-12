package local

type ClientInterface interface {
	ListObjects() (*Object, error)
}
