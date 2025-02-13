package local

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestObjectCheckSumCalculation(t *testing.T) {
	o := Object{
		FullPath: "../examples/same/index.html",
	}

	cs, err := o.Checksum()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	v := fmt.Sprintf("%x", cs)
	if diff := cmp.Diff(v, "e2b4958d41eaa73afbdb5a4b5fad3321"); diff != "" {
		t.Errorf("Checksum calculation mismatch, diff=%s", diff)
	}
}
