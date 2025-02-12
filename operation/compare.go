package operation

import (
	"bytes"
	"path/filepath"

	"github.com/ysugimoto/gcsdeploy/local"
	"github.com/ysugimoto/gcsdeploy/remote"
)

func Make(r *remote.Object, l *local.Object) (Operations, error) {
	ops := Operations{}

	// Make Add / Update operation
	for key := range l.Items {
		if rcs, ok := r.Items[key]; ok {
			// If object exists in remote, calculate checksum and compare both
			lcs, err := l.Checksum(key)
			if err != nil {
				return nil, err
			}
			// If checksum is different, create update operation
			if !bytes.Equal(rcs, lcs) {
				ops = append(ops, Operation{
					Type:   Update,
					Local:  filepath.Join(l.Root, key),
					Remote: key,
				})
			}
			continue
		}

		// If object not exists in remote, create add operation
		ops = append(ops, Operation{
			Type:   Add,
			Local:  filepath.Join(l.Root, key),
			Remote: key,
		})
	}

	// Make Delete operation
	for key := range r.Items {
		if _, ok := l.Items[key]; ok {
			continue
		}
		// If object is not in local, create delete operation
		ops = append(ops, Operation{
			Type:   Delete,
			Remote: key,
		})
	}

	return ops, nil
}
