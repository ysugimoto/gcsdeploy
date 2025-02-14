package remote

import (
	"net/url"
	"strings"
)

// Bucket struct holds the bucket name and prefix.
// CLI flag of bucket string may includes a sub-directory like gs://[bucket-name]/[sub-directory].
type Bucket struct {
	Name   string
	Prefix string
}

// Get GCS bucket URL.
func (b *Bucket) String() string {
	str := "gs://" + b.Name
	if b.Prefix != "" {
		str += "/" + b.Prefix
	}
	return str
}

// Parse from provided bucket string to Bucket pointer.
func ParseBucket(bucket string) (*Bucket, error) {
	parsed, err := url.Parse(bucket)
	if err != nil {
		return nil, Error(err, "Failed to parse bucket string")
	}

	// Provided with gs:// protocol like gs://[bucket-name]
	if parsed.Scheme == "gs" {
		return &Bucket{
			Name:   parsed.Host,
			Prefix: strings.Trim(parsed.Path, "/"),
		}, nil
	}

	// Otherwise, provide bucket namne like "bucket-name"
	spl := strings.SplitN(parsed.Path, "/", 2)
	var prefix string
	if len(spl) > 1 {
		prefix = spl[1]
	}
	return &Bucket{
		Name:   spl[0],
		Prefix: strings.Trim(prefix, "/"),
	}, nil
}
