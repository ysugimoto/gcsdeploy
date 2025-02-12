package remote

import (
	"context"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Client struct {
	c      *storage.Client
	bucket string
}

func New(ctx context.Context, bucket string) (ClientInterface, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, Error(err, "Failed to create GCS storage client")
	}

	// If bucket name starts with "gs://" protocol, trim it
	bucket = strings.TrimPrefix(bucket, "gs://")

	return &Client{
		c:      client,
		bucket: bucket,
	}, nil
}

func (client *Client) ListObjects(ctx context.Context) (*Object, error) {
	iter := client.c.Bucket(client.bucket).Objects(ctx, &storage.Query{
		Versions:   false,
		Projection: storage.ProjectionNoACL,
	})
	o := &Object{
		Bucket: client.bucket,
		Items:  make(map[string]Checksum),
	}

	for {
		v, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, Error(err, "Failed to iterate object")
		}
		o.Items[v.Name] = v.MD5
	}

	return o, nil
}

func (client *Client) UploadObject(ctx context.Context, from, to string) error {
	fp, err := os.Open(from)
	if err != nil {
		return Error(err, "Failed to open file "+from)
	}
	defer fp.Close()

	w := client.c.Bucket(client.bucket).Object(to).NewWriter(ctx)
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
	if err := client.c.Bucket(client.bucket).Object(from).Delete(ctx); err != nil {
		return Error(err, "Failed to delete object "+from)
	}
	return nil
}
