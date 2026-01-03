# zei

A basic CLI for storing and executing command snippets.

```bash
# Display help info
zei -h

# Execute stored snippet's command
zei <command-id>
```

> [!TIP]
> Go templating is supported within a snippet's command.
>
> For example, setting a snippet with a command `echo {{.Message}}`
> will prompt for the value of `Message` when the snippet is executed.

## Installing

Download a binary from [releases](https://github.com/Sammy-T/zei/releases) and add it to your PATH.

**or**

Install with go:

```bash
go install github.com/sammy-t/zei/cmd/zei@latest
```

## Development

#### Add Go dependencies

```bash
go get ./...
```

### Run the CLI

```bash
go run ./cmd/zei
```

### Build the CLI

```bash
go build -C ./cmd/zei
```

The binary will output to the `cmd/zei/` directory.

### Install the CLI

```bash
go install ./cmd/zei
```
