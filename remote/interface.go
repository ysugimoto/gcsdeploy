package remote

import (
	"context"

	"github.com/ysugimoto/gcsdeploy/local"
)

// ClientInterface declares an interface that needs to have required mehods
type ClientInterface interface {
	ListObjects(context.Context, *Bucket) (Objects, error)
	UploadObject(context.Context, local.Object, Object) error
	DeleteObject(context.Context, Object) error
}
