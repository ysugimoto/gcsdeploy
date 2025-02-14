package local

// ClientInterface declares an interface that needs to have required mehods
type ClientInterface interface {
	ListObjects() (Objects, error)
}
