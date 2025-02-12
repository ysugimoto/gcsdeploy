package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ysugimoto/gcsdeploy/operation"
)

func TestDivideOperationByConcurrency(t *testing.T) {
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
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
			},
			expects: []operation.Operations{
				{
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				},
				{
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				},
			},
		},
		{
			name:        "Concurrency 3",
			concurrency: 3,
			operations: operation.Operations{
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
			},
			expects: []operation.Operations{
				{
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				},
				{
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
				},
				{
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
					{Type: operation.Add, Local: "../../examples/same/index.html", Remote: "index.html"},
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
