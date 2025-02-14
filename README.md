# gcsdeploy

Simple CLI tool that manages deployment to the Google Cloud Storage bucket.
This tool uses MD5 checksum in GCS object, and compare with local file, and update if file has changed.

This tool is heavily inspired in [s3deploy](https://github.com/bep/s3deploy).

## Installation

Download a latest binary from [Releases](https://github.com/ysugimoto/gcsdeploy/releases) and put it to the `$PATH`.
Or install via command:

```shell
go install github.com/ysugimoto/gcsdeploy
```

## Usage / Options

This is a single command tool with kind of options:

```shell
NAME:
   gcsdeploy - GCS deploy management with rsync-like operation

USAGE:
   gcsdeploy [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dry-run                      Dry run (default: false)
   --delete                       Delete GCS object if not exists in local (default: false)
   --bucket value, -b value       Specify deploy destination bucket
   --source value, -s value       Specify local root directory to deploy (default: ".")
   --credential value             Specify credential file path
   --concurrency value, -c value  Specify operation concurrency (default: 1)
   --help, -h                     show help
```

### --dry-run

Dry-Run, print operation plan what the GCS bucket will be modified.
We don't modify GCS object if this flag is provided.

### --delete

Enable delete operation that deletes GCS object that is not exist in local.
Default is `false` so if you need to completely synchronize between local and GCS bucket, provide this option.

### --bucket, -b

Specify the destination GCS bucket to modify.

### --source, -s

Specify root of local files. Default is the current working directory.

### --credential

Provide custom Google Cloud credential filepath.
This tool uses application default credential as default, use this credential file if provided.

### --concurrency, -c

Specify the concurrency for each operations.
Default is `1`, and it can get better performance of deployment if this value is increased.
The max concurrency is `10`.


## Contribution

- Fork this repository
- Customize / Fix problem
- Send PR :-)
- Or feel free to create issues for us. We'll look into it

## License

MIT License

## Contributors

- [@ysugimoto](https://github.com/ysugimoto)
