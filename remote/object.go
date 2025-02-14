package remote

type Objects map[string]Object

// Remote object declaration.
// struct should have key, GCS bucket info, checksum mimetype, and its filesize.
type Object struct {
	Key         string
	Bucket      *Bucket
	Checksum    []byte
	ContentType string
	Size        int64
}

// Returns URL-formed string of GCS bucket.
func (o Object) URL() string {
	return o.Bucket.String() + "/" + o.Key
}

// Returns path in the GCS bucket from root.
func (o Object) Path() string {
	if o.Bucket.Prefix != "" {
		return o.Bucket.Prefix + "/" + o.Key
	}
	return o.Key
}
