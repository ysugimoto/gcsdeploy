package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/ysugimoto/gcsdeploy/operation"
)

var (
	output = colorable.NewColorableStderr()
	yellow = color.New(color.FgYellow)
	white  = color.New(color.FgWhite)
	red    = color.New(color.FgRed)
	green  = color.New(color.FgGreen)
)

func write(c *color.Color, format string, args ...interface{}) {
	c.Fprint(output, fmt.Sprintf(format, args...))
}

func writeln(c *color.Color, format string, args ...interface{}) {
	write(c, format+"\n", args...)
}

func printDryRunOperations(ops operation.Operations, bucket string) {
	var add, update, del operation.Operations
	for i := range ops {
		switch ops[i].Type {
		case operation.Add:
			add = append(add, ops[i])
		case operation.Update:
			update = append(update, ops[i])
		case operation.Delete:
			del = append(del, ops[i])
		}
	}
	writeln(white, "================================ %s ================================", "Dry-Run Result")
	if len(add) > 0 {
		writeln(green, "%s", "Add Files:")
		for i := range add {
			writeln(white, "%2s- %s -> gs://%s/%s", " ", add[i].Local, bucket, add[i].Remote)
		}
	}
	if len(update) > 0 {
		writeln(yellow, "%s", "Update Files:")
		for i := range update {
			writeln(white, "%2s- %s -> gs://%s/%s", " ", update[i].Local, bucket, update[i].Remote)
		}
	}
	if len(del) > 0 {
		writeln(white, "%s", "Delete Files:")
		for i := range del {
			writeln(red, "%2s- %s -> gs://%s/%s", " ", del[i].Local, bucket, del[i].Remote)
		}
	}
	writeln(white, "================================================================================")
}

func printAddOperation(op operation.Operation, bucket string) {
	writeln(white, "%8s: %s -> gs://%s/%s ...", "Adding", op.Local, bucket, op.Remote)
}

func printUpdateOperation(op operation.Operation, bucket string) {
	writeln(white, "%8s: %s -> gs://%s/%s ...", "Updating", op.Local, bucket, op.Remote)
}

func printDeleteOperation(op operation.Operation, bucket string) {
	writeln(white, "%8s gs://%s/%s ...", "Deleting", bucket, op.Remote)
}
