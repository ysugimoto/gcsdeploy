package local

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	root string
}

func New(root string) (ClientInterface, error) {
	if _, err := os.Stat(root); err != nil {
		return nil, Error(err, "Failed to stat root path of "+root)
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, Error(err, "Failed to convert to absolute path "+root)
	}
	return &Client{
		root: abs + "/",
	}, nil
}

func (client *Client) ListObjects() (*Object, error) {
	o := &Object{
		Root:  client.root,
		Items: make(map[string]Checksum),
	}

	err := filepath.Walk(client.root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip if walked path is the directory
		if info.IsDir() {
			return nil
		}
		// Add items, but we won't calculate checksum at this time.
		// We should calculate checksum when we need to compare between remote and local for performance.
		o.Items[strings.TrimPrefix(path, client.root)] = nil
		return nil
	})
	if err != nil {
		return nil, Error(err, "Failed to find target local files")
	}

	return o, nil
}
