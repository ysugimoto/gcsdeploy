package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	cli "github.com/urfave/cli/v2"
	"github.com/ysugimoto/gcsdeploy/local"
	"github.com/ysugimoto/gcsdeploy/operation"
	"github.com/ysugimoto/gcsdeploy/remote"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/option"
)

// Declare CLI command flag names
const (
	flagNameDryRun          = "dry-run"
	flagNameBucket          = "bucket"
	flagNameSource          = "source"
	flagNameConcurrency     = "concurrency"
	flagNameCrendentialPath = "credential"
	flagNameDelete          = "delete"
)

func main() {
	app := &cli.App{
		Name:  "gcsdeploy",
		Usage: "GCS deploy management with rsync-like operation",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  flagNameDryRun,
				Usage: "Dry run",
			},
			&cli.BoolFlag{
				Name:  flagNameDelete,
				Usage: "Delete GCS object if not exists in local",
			},
			&cli.StringFlag{
				Name:     flagNameBucket,
				Aliases:  []string{"b"},
				Usage:    "Specify deploy destination bucket",
				Required: true,
				Action: func(ctx *cli.Context, v string) error {
					if _, err := remote.ParseBucket(v); err != nil {
						return fmt.Errorf("Invalid bucket provided: %s", v)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    flagNameSource,
				Aliases: []string{"s"},
				Usage:   "Specify local root directory to deploy",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:  flagNameCrendentialPath,
				Usage: "Specify credential file path",
			},
			&cli.IntFlag{
				Name:    flagNameConcurrency,
				Aliases: []string{"c"},
				Usage:   "Specify operation concurrency",
				Value:   1,
				Action: func(ctx *cli.Context, v int) error {
					if v < 1 || v > 10 {
						return fmt.Errorf("Concurrency flag %d is out of range [1-10]", v)
					}
					return nil
				},
			},
		},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// Execute CLI action
func action(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We don't need error check because validation has already done in CLI flag parsing
	bucket, _ := remote.ParseBucket(c.String(flagNameBucket)) // nolint:errcheck
	source := c.String(flagNameSource)
	credential := c.String(flagNameCrendentialPath)
	concurrency := c.Int(flagNameConcurrency)
	enableDelete := c.Bool(flagNameDelete)

	var options []option.ClientOption
	if credential != "" {
		options = append(options, option.WithCredentialsFile(credential))
	}
	r, err := remote.New(ctx, options...)
	if err != nil {
		return err
	}
	dest, err := r.ListObjects(ctx, bucket)
	if err != nil {
		return err
	}
	l, err := local.New(source)
	if err != nil {
		return err
	}
	src, err := l.ListObjects()
	if err != nil {
		return err
	}
	ops, err := operation.Make(bucket, dest, src)
	if err != nil {
		return err
	}

	// If --dry-run flag is provided, only prints operation plans
	if c.Bool(flagNameDryRun) {
		printDryRunOperations(ops, enableDelete)
		return nil
	}

	// Execute operations for each concurrency
	var messages []string
	for _, task := range divideOperationsByConcurrency(ops, concurrency) {
		if err := runTask(ctx, r, task, enableDelete); err != nil {
			messages = append(messages, err.Error())
		}
	}

	if len(messages) > 0 {
		return fmt.Errorf("%s", strings.Join(messages, "\n"))
	}

	return nil
}

// runTask runs for each task unit that is divided by parallel count
func runTask(ctx context.Context, r remote.ClientInterface, tasks operation.Operations, enableDelete bool) error {
	var eg errgroup.Group
	for i := range tasks {
		task := tasks[i] // trap variable in this scope
		eg.Go(func() error {
			var err error
			switch task.Type {
			case operation.Add:
				printAddOperation(task)
				err = r.UploadObject(ctx, task.Local, task.Remote)
			case operation.Update:
				printUpdateOperation(task)
				err = r.UploadObject(ctx, task.Local, task.Remote)
			case operation.Delete:
				if !enableDelete {
					printDeleteOperation(task)
					err = r.DeleteObject(ctx, task.Remote)
				}
			}
			return err
		})
	}
	return eg.Wait()
}

// divideOperationsByConcurrency divides each task unit per the concurrency
func divideOperationsByConcurrency(ops operation.Operations, concurrency int) (tasks []operation.Operations) {
	var task operation.Operations
	for i := range ops {
		task = append(task, ops[i])
		if len(task) == concurrency {
			tasks = append(tasks, task)
			task = operation.Operations{}
		}
	}
	if len(task) > 0 {
		tasks = append(tasks, task)
	}
	return
}
