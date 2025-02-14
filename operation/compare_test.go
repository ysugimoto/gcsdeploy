package operation

import (
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ysugimoto/gcsdeploy/local"
	"github.com/ysugimoto/gcsdeploy/remote"
)

func calculateChecksum(p string) []byte {
	fp, _ := os.Open(p)
	defer fp.Close()

	hash := md5.New()
	io.Copy(hash, fp)
	return hash.Sum(nil)
}

func TestMake(t *testing.T) {
	root, _ := filepath.Abs("../examples/same")
	diff, _ := filepath.Abs("../examples/diff")
	bucket := &remote.Bucket{Name: "example-bucket"}

	tests := []struct {
		name    string
		remote  remote.Objects
		local   local.Objects
		expects any
	}{
		{
			name: "No Operations when remote and local are completely same",
			remote: remote.Objects{
				"index.html": {
					Key:      "index.html",
					Checksum: calculateChecksum(filepath.Join(root, "index.html")),
					Bucket:   bucket,
				},
				"vite.svg": {
					Key:      "vite.svg",
					Checksum: calculateChecksum(filepath.Join(root, "vite.svg")),
					Bucket:   bucket,
				},
				"assets/index-n_ryQ3BS.css": {
					Key:      "assets/index-n_ryQ3BS.css",
					Checksum: calculateChecksum(filepath.Join(root, "assets/index-n_ryQ3BS.css")),
					Bucket:   bucket,
				},
				"assets/index-pGAOdsKC.js": {
					Key:      "assets/index-pGAOdsKC.js",
					Checksum: calculateChecksum(filepath.Join(root, "assets/index-pGAOdsKC.js")),
					Bucket:   bucket,
				},
				"assets/react-CHdo91hT.svg": {
					Key:      "assets/react-CHdo91hT.svg",
					Checksum: calculateChecksum(filepath.Join(root, "assets/react-CHdo91hT.svg")),
					Bucket:   bucket,
				},
			},
			local: local.Objects{
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
			expects: Operations{},
		},
		{
			name: "Add operation will present",
			remote: remote.Objects{
				"index.html": {
					Key:      "index.html",
					Checksum: calculateChecksum(filepath.Join(root, "index.html")),
					Bucket:   bucket,
				},
				"vite.svg": {
					Key:      "vite.svg",
					Checksum: calculateChecksum(filepath.Join(root, "vite.svg")),
					Bucket:   bucket,
				},
				"assets/index-n_ryQ3BS.css": {
					Key:      "assets/index-n_ryQ3BS.css",
					Checksum: calculateChecksum(filepath.Join(root, "assets/index-n_ryQ3BS.css")),
					Bucket:   bucket,
				},
				"assets/index-pGAOdsKC.js": {
					Key:      "assets/index-pGAOdsKC.js",
					Checksum: calculateChecksum(filepath.Join(root, "assets/index-pGAOdsKC.js")),
					Bucket:   bucket,
				},
				"assets/react-CHdo91hT.svg": {
					Key:      "assets/react-CHdo91hT.svg",
					Checksum: calculateChecksum(filepath.Join(root, "assets/react-CHdo91hT.svg")),
					Bucket:   bucket,
				},
			},
			local: local.Objects{
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
				"some.html": {
					FullPath:    filepath.Join(root, "some.html"),
					ContentType: "text/html",
					Size:        464,
				},
			},
			expects: Operations{
				{
					Type: Add,
					Local: local.Object{
						FullPath:    filepath.Join(root, "some.html"),
						ContentType: "text/html",
						Size:        464,
					},
					Remote: remote.Object{
						Key:    "some.html",
						Bucket: bucket,
					},
				},
			},
		},
		{
			name: "Update operation will present when checksum is different",
			remote: remote.Objects{
				"index.html": {
					Key:      "index.html",
					Checksum: calculateChecksum(filepath.Join(root, "index.html")),
					Bucket:   bucket,
				},
				"vite.svg": {
					Key:      "vite.svg",
					Checksum: calculateChecksum(filepath.Join(root, "vite.svg")),
					Bucket:   bucket,
				},
				"assets/index-n_ryQ3BS.css": {
					Key:      "assets/index-n_ryQ3BS.css",
					Checksum: calculateChecksum(filepath.Join(root, "assets/index-n_ryQ3BS.css")),
					Bucket:   bucket,
				},
				"assets/index-pGAOdsKC.js": {
					Key:      "assets/index-pGAOdsKC.js",
					Checksum: calculateChecksum(filepath.Join(root, "assets/index-pGAOdsKC.js")),
					Bucket:   bucket,
				},
				"assets/react-CHdo91hT.svg": {
					Key:      "assets/react-CHdo91hT.svg",
					Checksum: calculateChecksum(filepath.Join(root, "assets/react-CHdo91hT.svg")),
					Bucket:   bucket,
				},
			},
			local: local.Objects{
				"index.html": {
					FullPath:    filepath.Join(diff, "index.html"),
					ContentType: "text/html",
					Size:        464,
				},
				"vite.svg": {
					FullPath:    filepath.Join(diff, "vite.svg"),
					ContentType: "image/svg+xml",
					Size:        1497,
				},
				"assets/index-n_ryQ3BS.css": {
					FullPath:    filepath.Join(diff, "assets/index-n_ryQ3BS.css"),
					ContentType: "text/css",
					Size:        1391,
				},
				"assets/index-pGAOdsKC.js": {
					FullPath:    filepath.Join(diff, "assets/index-pGAOdsKC.js"),
					ContentType: "text/javascript",
					Size:        143899,
				},
				"assets/react-CHdo91hT.svg": {
					FullPath:    filepath.Join(diff, "assets/react-CHdo91hT.svg"),
					ContentType: "image/svg+xml",
					Size:        4126,
				},
			},
			expects: Operations{
				{
					Type: Update,
					Local: local.Object{
						FullPath:    filepath.Join(diff, "index.html"),
						ContentType: "text/html",
						Size:        464,
					},
					Remote: remote.Object{
						Key:      "index.html",
						Checksum: calculateChecksum(filepath.Join(root, "index.html")),
						Bucket:   bucket,
					},
				},
			},
		},
		{
			name: "Delete operation will present when not exists in local",
			remote: remote.Objects{
				"index.html": {
					Key:      "index.html",
					Checksum: calculateChecksum(filepath.Join(root, "index.html")),
					Bucket:   bucket,
				},
				"vite.svg": {
					Key:      "vite.svg",
					Checksum: calculateChecksum(filepath.Join(root, "vite.svg")),
					Bucket:   bucket,
				},
				"assets/index-n_ryQ3BS.css": {
					Key:      "assets/index-n_ryQ3BS.css",
					Checksum: calculateChecksum(filepath.Join(root, "assets/index-n_ryQ3BS.css")),
					Bucket:   bucket,
				},
				"assets/index-pGAOdsKC.js": {
					Key:      "assets/index-pGAOdsKC.js",
					Checksum: calculateChecksum(filepath.Join(root, "assets/index-pGAOdsKC.js")),
					Bucket:   bucket,
				},
				"assets/react-CHdo91hT.svg": {
					Key:      "assets/react-CHdo91hT.svg",
					Checksum: calculateChecksum(filepath.Join(root, "assets/react-CHdo91hT.svg")),
					Bucket:   bucket,
				},
			},
			local: local.Objects{
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
			expects: Operations{
				{
					Type: Delete,
					Remote: remote.Object{
						Key:      "index.html",
						Checksum: calculateChecksum(filepath.Join(root, "index.html")),
						Bucket:   bucket,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ops, err := Make(bucket, tt.remote, tt.local)
			if err != nil {
				t.Errorf("Unexpected error returns: %s", err)
				return
			}
			if diff := cmp.Diff(ops, tt.expects); diff != "" {
				t.Errorf("Operation result mismatch, diff=%s", diff)
			}
		})
	}
}
