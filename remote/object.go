package remote

type Checksum []byte

type Object struct {
	// Bucket holds the bucket name to operate.
	Bucket string

	// Items holds the object map.
	// Key string is the path, and the value []byte is the checksum.
	Items map[string]Checksum
}
