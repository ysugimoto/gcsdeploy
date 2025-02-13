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
)

const (
	flagNameDryRun          = "dry-run"
	flagNameBucket          = "bucket"
	flagNameLocalPath       = "local"
	flagNameConcurrency     = "concurrency"
	flagNameCrendentialPath = "credential"
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
				Name:     flagNameLocalPath,
				Aliases:  []string{"l"},
				Usage:    "Specify local root directory to deploy",
				Required: true,
			},
			&cli.StringFlag{
				Name:    flagNameCrendentialPath,
				Aliases: []string{"c"},
				Usage:   "Specify credential file path",
			},
			&cli.IntFlag{
				Name:    flagNameConcurrency,
				Aliases: []string{"p"},
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

func action(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We don't need error check because validation has already done in CLI flag parsing
	bucket, _ := remote.ParseBucket(c.String(flagNameBucket)) // nolint:errcheck
	localPath := c.String(flagNameLocalPath)
	credential := c.String(flagNameCrendentialPath)
	concurrency := c.Int(flagNameConcurrency)

	var r remote.ClientInterface
	var err error
	if credential == "" {
		r, err = remote.New(ctx, bucket)
	} else {
		r, err = remote.NewWithCredential(ctx, bucket, credential)
	}
	if err != nil {
		return err
	}
	dest, err := r.ListObjects(ctx)
	if err != nil {
		return err
	}
	l, err := local.New(localPath)
	if err != nil {
		return err
	}
	src, err := l.ListObjects()
	if err != nil {
		return err
	}
	ops, err := operation.Make(dest, src)
	if err != nil {
		return err
	}

	// If --dry-run flag is provided, print operation detail
	if c.Bool(flagNameDryRun) {
		printDryRunOperations(ops, bucket)
		return nil
	}

	// Execute operations for each concurrency
	var messages []string
	for _, task := range divideOperationsByConcurrency(ops, concurrency) {
		if err := runTask(task, bucket); err != nil {
			messages = append(messages, err.Error())
		}
	}

	if len(messages) > 0 {
		return fmt.Errorf("%s", strings.Join(messages, "\n"))
	}

	return nil
}

func runTask(tasks operation.Operations, bucket *remote.Bucket) error {
	var eg errgroup.Group
	for i := range tasks {
		task := tasks[i] // trap variable in this scope
		eg.Go(func() error {
			var err error
			switch task.Type {
			case operation.Add:
				printAddOperation(task, bucket)
				// err = r.UploadObject(ctx, task[i].Local, task[i].Remote)
			case operation.Update:
				printUpdateOperation(task, bucket)
				// err = r.UploadObject(ctx, task[i].Local, task[i].Remote)
			case operation.Delete:
				printDeleteOperation(task, bucket)
				// err = r.DeleteObject(ctx, task[i].Remote)
			}
			return err
		})
	}
	return eg.Wait()
}

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
