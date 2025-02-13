package local

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestListNewClient(t *testing.T) {
	tests := []struct {
		name    string
		root    string
		isError bool
	}{
		{
			name:    "raise error when root is not found",
			root:    "/path/to/notfound",
			isError: true,
		},
		{
			name: "create client successfully",
			root: "../examples/same",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.root)
			if err != nil {
				if !tt.isError {
					t.Errorf("Unexpected error returned: %s", err)
				}
				return
			}
		})
	}
}

func TestListObjects(t *testing.T) {
	root, _ := filepath.Abs("../examples/same")
	tests := []struct {
		name    string
		root    string
		expects any
	}{
		{
			name: "Return expected files",
			root: root,
			expects: Objects{
				"index.html": {
					FullPath:    filepath.Join(root, "index.html"),
					ContentType: "text/html",
					Size:        464,
				},
				"vite.svg": {
					FullPath:    filepath.Join(root, "vite.svg"),
					ContentType: "image/svg+xml",
					Size:        1497,
				},
				"assets/index-n_ryQ3BS.css": {
					FullPath:    filepath.Join(root, "assets/index-n_ryQ3BS.css"),
					ContentType: "text/css",
					Size:        1391,
				},
				"assets/index-pGAOdsKC.js": {
					FullPath:    filepath.Join(root, "assets/index-pGAOdsKC.js"),
					ContentType: "text/javascript",
					Size:        143899,
				},
				"assets/react-CHdo91hT.svg": {
					FullPath:    filepath.Join(root, "assets/react-CHdo91hT.svg"),
					ContentType: "image/svg+xml",
					Size:        4126,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.root)
			if err != nil {
				t.Errorf("Unexpected error returned: %s", err)
			}
			ret, err := c.ListObjects()
			if err != nil {
				t.Errorf("Unexpected error returned: %s", err)
			}
			if diff := cmp.Diff(ret, tt.expects); diff != "" {
				t.Errorf("ListObjects() result mismatch, diff=%s", diff)
			}
		})
	}
}
