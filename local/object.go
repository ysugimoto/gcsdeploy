package local

import (
	"crypto/md5"
	"io"
	"os"
)

type Object struct {
	FullPath    string
	ContentType string
	Size        int64
}

func (o Object) Checksum() ([]byte, error) {
	fp, err := os.Open(o.FullPath)
	if err != nil {
		return nil, Error(err, "Failed open file "+o.FullPath)
	}
	defer fp.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, fp); err != nil {
		return nil, Error(err, "Failed to create hash for file "+o.FullPath)
	}
	sum := hash.Sum(nil)
	return sum, nil
}

type Objects map[string]Object
