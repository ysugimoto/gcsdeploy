package main

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ysugimoto/gcsdeploy/local"
	"github.com/ysugimoto/gcsdeploy/operation"
	"github.com/ysugimoto/gcsdeploy/remote"
)

func TestDivideOperationByConcurrency(t *testing.T) {
	root, _ := filepath.Abs("../../examples/same")
	tests := []struct {
		name        string
		concurrency int
		operations  operation.Operations
		expects     any
	}{
		{
			name:        "Concurrency 1",
			concurrency: 1,
			operations: operation.Operations{
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
			},
			expects: []operation.Operations{
				{
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
				},
				{
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
				},
			},
		},
		{
			name:        "Concurrency 3",
			concurrency: 3,
			operations: operation.Operations{
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
				{
					Type: operation.Add,
					Local: local.Object{
						FullPath: filepath.Join(root, "index.html"),
					},
					Remote: remote.Object{
						Key: "index.html",
					},
				},
			},
			expects: []operation.Operations{
				{
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
				},
				{
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
				},
				{
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
					{
						Type: operation.Add,
						Local: local.Object{
							FullPath: filepath.Join(root, "index.html"),
						},
						Remote: remote.Object{
							Key: "index.html",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks := divideOperationsByConcurrency(tt.operations, tt.concurrency)
			if diff := cmp.Diff(tasks, tt.expects); diff != "" {
				t.Errorf("Divided tasks result mismatch, diff=%s", diff)
			}
		})
	}
}
