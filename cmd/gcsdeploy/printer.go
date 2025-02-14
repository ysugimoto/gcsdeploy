package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/ysugimoto/gcsdeploy/operation"
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

// Format int to thousand comma string
func numberFormat(v int64) string {
	return message.NewPrinter(language.English).Sprintf("%d", v)
}

// Print operation plan (add, update, delete)
func printDryRunOperations(ops operation.Operations, enableDelete bool) {
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
				"- %s -> %s (%s %s bytes)",
				a.Local.FullPath,
				a.Remote.URL(),
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
				"- %s -> %s (%s %s bytes)",
				u.Local.FullPath,
				u.Remote.URL(),
				u.Local.ContentType,
				numberFormat(u.Local.Size),
			)
		}
	}
	if len(del) > 0 {
		str := "Delete Files"
		if !enableDelete {
			str += " (will not execute, provide --delete flag if you'd like to be enable)"
		}
		str += ":"
		writeln(red, "%s", str)
		for _, d := range del {
			writeln(
				white,
				"- %s (%s %s bytes)",
				d.Remote.URL(),
				d.Remote.ContentType,
				numberFormat(d.Remote.Size),
			)
		}
	}
}

// Print add operation information
func printAddOperation(op operation.Operation) {
	write(green, "%8s: ", "Adding")
	writeln(
		white,
		"%s -> %s (%s %s bytes)",
		op.Local.FullPath,
		op.Remote.URL(),
		op.Local.ContentType,
		numberFormat(op.Local.Size),
	)
}

// Print update operation information
func printUpdateOperation(op operation.Operation) {
	write(yellow, "%8s: ", "Updating")
	writeln(
		white,
		"%s -> %s (%s %s bytes)",
		op.Local.FullPath,
		op.Remote.URL(),
		op.Local.ContentType,
		numberFormat(op.Local.Size),
	)
}

// Print delete operation information
func printDeleteOperation(op operation.Operation) {
	write(red, "%8s: ", "Deleting")
	writeln(
		white,
		"%s (%s %s bytes)",
		op.Remote.URL(),
		op.Remote.ContentType,
		numberFormat(op.Remote.Size),
	)
}
