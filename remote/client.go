package remote

import (
	"context"
	"io"
	"mime"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	c      *storage.Client
	bucket *Bucket
}

func New(ctx context.Context, bucket *Bucket) (ClientInterface, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, Error(err, "Failed to create GCS storage client")
	}

	return &Client{
		c:      client,
		bucket: bucket,
	}, nil
}

func NewWithCredential(ctx context.Context, bucket *Bucket, creds string) (ClientInterface, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(creds))
	if err != nil {
		return nil, Error(err, "Failed to create GCS storage client")
	}

	return &Client{
		c:      client,
		bucket: bucket,
	}, nil
}

func (client *Client) ListObjects(ctx context.Context) (Objects, error) {
	iter := client.c.Bucket(client.bucket.Name).Objects(ctx, &storage.Query{
		Versions:   false,
		Projection: storage.ProjectionNoACL,
		Prefix:     client.bucket.Prefix,
	})
	o := Objects{}

	for {
		v, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, Error(err, "Failed to iterate object")
		}
		o[v.Name] = Object{
			Key:         v.Name,
			Bucket:      client.bucket,
			Checksum:    v.MD5,
			Size:        v.Size,
			ContentType: v.ContentType,
		}
	}

	return o, nil
}

func (client *Client) UploadObject(ctx context.Context, from, to string) error {
	fp, err := os.Open(from)
	if err != nil {
		return Error(err, "Failed to open file "+from)
	}
	defer fp.Close()

	w := client.c.Bucket(client.bucket.Name).Object(to).NewWriter(ctx)
	m, err := mime.ExtensionsByType(filepath.Ext(from))
	if err == nil && m != nil {
		w.ContentType = m[0]
	}
	if _, err := io.Copy(w, fp); err != nil {
		return Error(err, "Failed to write remote file "+from)
	}
	return nil
}

func (client *Client) DeleteObject(ctx context.Context, from string) error {
	if err := client.c.Bucket(client.bucket.Name).Object(from).Delete(ctx); err != nil {
		return Error(err, "Failed to delete object "+from)
	}
	return nil
}
