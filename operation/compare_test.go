package operation

import (
	"crypto/md5"
	"io"
	"os"
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
	tests := []struct {
		name    string
		remote  *remote.Object
		local   *local.Object
		expects any
	}{
		{
			name: "No Operations when remote and local are completely same",
			remote: &remote.Object{
				Items: map[string]remote.Checksum{
					"index.html":                calculateChecksum("../examples/same/index.html"),
					"vite.svg":                  calculateChecksum("../examples/same/vite.svg"),
					"assets/index-n_ryQ3BS.css": calculateChecksum("../examples/same/assets/index-n_ryQ3BS.css"),
					"assets/index-pGAOdsKC.js":  calculateChecksum("../examples/same/assets/index-pGAOdsKC.js"),
					"assets/react-CHdo91hT.svg": calculateChecksum("../examples/same/assets/react-CHdo91hT.svg"),
				},
			},
			local: &local.Object{
				Root: "../examples/same/",
				Items: map[string]local.Checksum{
					"index.html":                nil,
					"vite.svg":                  nil,
					"assets/index-n_ryQ3BS.css": nil,
					"assets/index-pGAOdsKC.js":  nil,
					"assets/react-CHdo91hT.svg": nil,
				},
			},
			expects: Operations{},
		},
		{
			name: "Add operation will present",
			remote: &remote.Object{
				Items: map[string]remote.Checksum{
					"index.html":                calculateChecksum("../examples/same/index.html"),
					"vite.svg":                  calculateChecksum("../examples/same/vite.svg"),
					"assets/index-n_ryQ3BS.css": calculateChecksum("../examples/same/assets/index-n_ryQ3BS.css"),
					"assets/index-pGAOdsKC.js":  calculateChecksum("../examples/same/assets/index-pGAOdsKC.js"),
					"assets/react-CHdo91hT.svg": calculateChecksum("../examples/same/assets/react-CHdo91hT.svg"),
				},
			},
			local: &local.Object{
				Root: "../examples/same/",
				Items: map[string]local.Checksum{
					"index.html":                nil,
					"vite.svg":                  nil,
					"assets/index-n_ryQ3BS.css": nil,
					"assets/index-pGAOdsKC.js":  nil,
					"assets/react-CHdo91hT.svg": nil,
					"some.html":                 nil,
				},
			},
			expects: Operations{
				{
					Type:   Add,
					Local:  "../examples/same/some.html",
					Remote: "some.html",
				},
			},
		},
		{
			name: "Update operation will present when checksum is different",
			remote: &remote.Object{
				Items: map[string]remote.Checksum{
					"index.html":                calculateChecksum("../examples/same/index.html"),
					"vite.svg":                  calculateChecksum("../examples/same/vite.svg"),
					"assets/index-n_ryQ3BS.css": calculateChecksum("../examples/same/assets/index-n_ryQ3BS.css"),
					"assets/index-pGAOdsKC.js":  calculateChecksum("../examples/same/assets/index-pGAOdsKC.js"),
					"assets/react-CHdo91hT.svg": calculateChecksum("../examples/same/assets/react-CHdo91hT.svg"),
				},
			},
			local: &local.Object{
				Root: "../examples/diff/",
				Items: map[string]local.Checksum{
					"index.html":                nil,
					"vite.svg":                  nil,
					"assets/index-n_ryQ3BS.css": nil,
					"assets/index-pGAOdsKC.js":  nil,
					"assets/react-CHdo91hT.svg": nil,
				},
			},
			expects: Operations{
				{
					Type:   Update,
					Local:  "../examples/diff/index.html",
					Remote: "index.html",
				},
			},
		},
		{
			name: "Delete operation will present when not exists in local",
			remote: &remote.Object{
				Items: map[string]remote.Checksum{
					"index.html":                calculateChecksum("../examples/same/index.html"),
					"vite.svg":                  calculateChecksum("../examples/same/vite.svg"),
					"assets/index-n_ryQ3BS.css": calculateChecksum("../examples/same/assets/index-n_ryQ3BS.css"),
					"assets/index-pGAOdsKC.js":  calculateChecksum("../examples/same/assets/index-pGAOdsKC.js"),
					"assets/react-CHdo91hT.svg": calculateChecksum("../examples/same/assets/react-CHdo91hT.svg"),
				},
			},
			local: &local.Object{
				Root: "../examples/diff/",
				Items: map[string]local.Checksum{
					"vite.svg":                  nil,
					"assets/index-n_ryQ3BS.css": nil,
					"assets/index-pGAOdsKC.js":  nil,
					"assets/react-CHdo91hT.svg": nil,
				},
			},
			expects: Operations{
				{
					Type:   Delete,
					Remote: "index.html",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ops, err := Make(tt.remote, tt.local)
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
