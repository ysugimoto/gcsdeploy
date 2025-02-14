package operation

import (
	"bytes"

	"github.com/ysugimoto/gcsdeploy/local"
	"github.com/ysugimoto/gcsdeploy/remote"
)

// Make creates operations that needs to run modification in GCS bucket (add, update, or delete).
// Compare remote files and local files and judge what operation should do
func Make(bucket *remote.Bucket, ro remote.Objects, lo local.Objects) (Operations, error) {
	ops := Operations{}

	// Make Add / Update operation
	for key, obj := range lo {
		if rcs, ok := ro[key]; ok {
			// If object exists in remote, calculate checksum and compare both
			lcs, err := obj.Checksum()
			if err != nil {
				return nil, err
			}
			// If checksum is different, create update operation
			if !bytes.Equal(rcs.Checksum, lcs) {
				ops = append(ops, Operation{
					Type:   Update,
					Local:  obj,
					Remote: rcs,
				})
			}
			continue
		}

		// If object not exists in remote, create add operation
		ops = append(ops, Operation{
			Type:  Add,
			Local: obj,
			Remote: remote.Object{
				Key:    key,
				Bucket: bucket,
			},
		})
	}

	// Make Delete operation
	for key, obj := range ro {
		if _, ok := lo[key]; ok {
			continue
		}
		// If object is not in local, create delete operation
		ops = append(ops, Operation{
			Type:   Delete,
			Remote: obj,
		})
	}

	return ops, nil
}
