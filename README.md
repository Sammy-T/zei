# zei

A basic CLI for storing and executing command snippets.

<img width="1001" height="216" alt="screen-01" src="https://github.com/user-attachments/assets/603b3efb-0cad-4985-bc31-65f5754b0b82" />

```bash
# Display help info
zei -h

# Execute stored snippet's command
zei <command-id>
```

> [!TIP]
> Go templating is supported within a snippet's command.
>
> For example, setting a snippet with a command like
>
> ```bash
> echo {{.Message}}
> ```
> 
> will prompt for the value of `Message` when the snippet is executed.

> [!TIP]
> Piped commands are supported within a snippet's command.
>
> Setting a snippet with a command like
> 
> ```bash
> <command-1> | <command-2>
> ```
> 
> will pipe command-1's output to command-2 when the snippet is executed.

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
