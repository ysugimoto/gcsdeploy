package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/ysugimoto/gcsdeploy/operation"
	"github.com/ysugimoto/gcsdeploy/remote"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

func numberFormat(v int64) string {
	return message.NewPrinter(language.English).Sprintf("%d", v)
}

func printDryRunOperations(ops operation.Operations, bucket *remote.Bucket) {
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
	writeln(white, "[Dry-Run]")
	if len(add) > 0 {
		writeln(green, "%s", "Add Files:")
		for _, a := range add {
			writeln(
				white,
				"- %s -> %s/%s (%s %s bytes)",
				a.Local.FullPath,
				bucket.String(),
				a.Remote.Key,
				a.Local.ContentType,
				numberFormat(a.Local.Size),
			)
		}
	}
	if len(update) > 0 {
		writeln(yellow, "%s", "Update Files:")
		for _, u := range update {
			writeln(
				white,
				"- %s -> %s/%s (%s %s bytes)",
				u.Local.FullPath,
				bucket.String(),
				u.Remote.Key,
				u.Local.ContentType,
				numberFormat(u.Local.Size),
			)
		}
	}
	if len(del) > 0 {
		writeln(red, "%s", "Delete Files:")
		for _, d := range del {
			writeln(
				white,
				"- %s/%s (%s %s bytes)",
				bucket.String(),
				d.Remote.Key,
				d.Remote.ContentType,
				numberFormat(d.Remote.Size),
			)
		}
	}
}

func printAddOperation(op operation.Operation, bucket *remote.Bucket) {
	writeln(
		white,
		"%8s: %s -> %s/%s (%s %s bytes)",
		"Adding",
		op.Local.FullPath,
		bucket.String(),
		op.Remote.Key,
		op.Local.ContentType,
		numberFormat(op.Local.Size),
	)
}

func printUpdateOperation(op operation.Operation, bucket *remote.Bucket) {
	writeln(
		white,
		"%8s: %s -> %s/%s (%s %s bytes)",
		"Updating",
		op.Local.FullPath,
		bucket.String(),
		op.Remote.Key,
		op.Local.ContentType,
		numberFormat(op.Local.Size),
	)
}

func printDeleteOperation(op operation.Operation, bucket *remote.Bucket) {
	writeln(
		white,
		"%8s: %s/%s (%s %s bytes)",
		"Deleting",
		bucket.String(),
		op.Remote.Key,
		op.Remote.ContentType,
		numberFormat(op.Remote.Size),
	)
}
