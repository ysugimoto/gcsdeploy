package remote

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseBucket(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect any
	}{
		{
			name:  "provide with gs:// protocol format",
			input: "gs://example-bucket",
			expect: &Bucket{
				Name:   "example-bucket",
				Prefix: "",
			},
		},
		{
			name:  "provide with gs:// protocol format with prefix",
			input: "gs://example-bucket/prefix",
			expect: &Bucket{
				Name:   "example-bucket",
				Prefix: "prefix",
			},
		},
		{
			name:  "provide with raw bucket name",
			input: "example-bucket",
			expect: &Bucket{
				Name:   "example-bucket",
				Prefix: "",
			},
		},
		{
			name:  "provide with raw bucket name with prefix",
			input: "example-bucket/prefix",
			expect: &Bucket{
				Name:   "example-bucket",
				Prefix: "prefix",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := ParseBucket(tt.input)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
				return
			}
			if diff := cmp.Diff(b, tt.expect); diff != "" {
				t.Errorf("Return bucket mismatch, diff=%s", diff)
			}
		})
	}
}
