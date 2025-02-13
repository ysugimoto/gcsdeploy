package operation

import (
	"bytes"

	"github.com/ysugimoto/gcsdeploy/local"
	"github.com/ysugimoto/gcsdeploy/remote"
)

func Make(ro remote.Objects, lo local.Objects) (Operations, error) {
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
				Key: key,
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
