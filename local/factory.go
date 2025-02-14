package local

import (
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	root string
}

// Create *Client pointer that implements ClientInterface
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

// Find local objects recursively
func (client *Client) ListObjects() (Objects, error) {
	o := Objects{}

	err := filepath.Walk(client.root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip if walked path is directory
		if info.IsDir() {
			return nil
		}
		m := mime.TypeByExtension(filepath.Ext(path))
		if idx := strings.Index(m, ";"); idx != -1 {
			// Trim charset section
			m = m[:idx]
		}
		// Add items, but we won't calculate checksum at this time.
		// We should calculate checksum when we need to compare between remote and local.
		o[strings.TrimPrefix(path, client.root)] = Object{
			ContentType: m,
			Size:        info.Size(),
			FullPath:    path,
		}

		return nil
	})
	if err != nil {
		return nil, Error(err, "Failed to find target local files")
	}

	return o, nil
}
