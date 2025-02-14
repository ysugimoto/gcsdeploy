package remote

import (
	"context"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/k0kubun/pp"
	"github.com/ysugimoto/gcsdeploy/local"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	c *storage.Client
}

// Create Client pointer that implements ClientInterface
func New(ctx context.Context, opts ...option.ClientOption) (ClientInterface, error) {
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, Error(err, "Failed to create GCS storage client")
	}

	return &Client{
		c: client,
	}, nil
}

// ListObjects lists all objects that exists in GCS bucket recursively
// We're guessing Objects() method will return all objects in the GCS bucket,
// but may need to call with next-page token to get all objects recursively.
func (client *Client) ListObjects(ctx context.Context, bucket *Bucket) (Objects, error) {
	iter := client.c.Bucket(bucket.Name).Objects(ctx, &storage.Query{
		Versions:   false,
		Projection: storage.ProjectionNoACL,
		Prefix:     bucket.Prefix,
	})
	o := Objects{}

	for {
		v, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			pp.Println(err)
			return nil, Error(err, "Failed to iterate object")
		}
		o[v.Name] = Object{
			Key:         v.Name,
			Bucket:      bucket,
			Checksum:    v.MD5,
			Size:        v.Size,
			ContentType: v.ContentType,
		}
	}

	return o, nil
}

// UploadObject uploads local file to GCS bucket with specified prefix
func (client *Client) UploadObject(ctx context.Context, from local.Object, to Object) error {
	fp, err := os.Open(from.FullPath)
	if err != nil {
		return Error(err, "Failed to open file "+from.FullPath)
	}
	defer fp.Close()

	w := client.c.Bucket(to.Bucket.Name).Object(to.Path()).NewWriter(ctx)
	if from.ContentType != "" {
		w.ContentType = from.ContentType
	}
	if _, err := io.Copy(w, fp); err != nil {
		return Error(err, "Failed to write remote file "+to.Path())
	}
	if err := w.Close(); err != nil {
		return Error(err, "Failed to flush writer buffer")
	}
	return nil
}

// Delete deletes file in GCS bucket.
func (client *Client) DeleteObject(ctx context.Context, from Object) error {
	if err := client.c.Bucket(from.Bucket.Name).Object(from.Path()).Delete(ctx); err != nil {
		return Error(err, "Failed to delete object "+from.Path())
	}
	return nil
}
