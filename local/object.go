package local

import (
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
)

type Checksum []byte

type Object struct {
	// Root holds the root path.
	Root string

	// Items holds the object map.
	// Key string is the path, and the value []byte is the checksum.
	Items map[string]Checksum
}

func (o *Object) Checksum(key string) (Checksum, error) {
	fp, err := os.Open(filepath.Join(o.Root, key))
	if err != nil {
		return nil, Error(err, "Failed to create path to "+key)
	}
	defer fp.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, fp); err != nil {
		return nil, Error(err, "Failed to create hash for file "+key)
	}
	sum := hash.Sum(nil)
	return sum, nil
}
