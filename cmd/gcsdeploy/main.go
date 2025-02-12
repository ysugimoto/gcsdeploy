package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"
	"github.com/ysugimoto/gcsdeploy/local"
	"github.com/ysugimoto/gcsdeploy/operation"
	"github.com/ysugimoto/gcsdeploy/remote"
	"golang.org/x/sync/errgroup"
)

const (
	flagNameDryRun      = "dry-run"
	flagNameBucket      = "bucket"
	flagNameLocalPath   = "local"
	flagNameConcurrency = "concurrency"
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
			},
			&cli.StringFlag{
				Name:     flagNameLocalPath,
				Aliases:  []string{"l"},
				Usage:    "Specify local root directory to deploy",
				Required: true,
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

func action(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bucket := c.String(flagNameBucket)
	localPath := c.String(flagNameLocalPath)
	concurrency := c.Int(flagNameConcurrency)

	r, err := remote.New(ctx, bucket)
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

	// If --dry-run flas is provided, print operation detail
	if c.Bool(flagNameDryRun) {
		// printOperations(ops)
		return nil
	}

	// Execute operations for each concurrency
	for _, task := range divideOperationsByConcurrency(ops, concurrency) {
		var eg errgroup.Group
		for i := range task {
			eg.Go(func() error {
				var err error
				switch task[i].Type {
				case operation.Add:
					err = r.UploadObject(ctx, task[i].Local, task[i].Remote)
				case operation.Update:
					err = r.UploadObject(ctx, task[i].Local, task[i].Remote)
				case operation.Delete:
					err = r.DeleteObject(ctx, task[i].Remote)
				}
				return err
			})
			if err := eg.Wait(); err != nil {
				return err
			}
		}
	}

	return nil
}

func divideOperationsByConcurrency(ops operation.Operations, concurrency int) []operation.Operations {
	tasks := []operation.Operations{}
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
	return tasks
}
