package remote

import "context"

type ClientInterface interface {
	ListObjects(context.Context) (*Object, error)
	UploadObject(context.Context, string, string) error
	DeleteObject(context.Context, string) error
}
