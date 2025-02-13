package remote

type Object struct {
	Key         string
	Bucket      *Bucket
	Checksum    []byte
	Size        int64
	ContentType string
}

func (o Object) Path() string {
	return o.Bucket.String() + "/" + o.Key
}

type Objects map[string]Object
