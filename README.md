# action-docs

A CLI utility that generates documentation for github actions and workflows.

## Install

### Go install

Run `go install` to download the binary to the go's binary folder:

```bash
go install github.com/nu12/action-docs@latest
```

Note: go's binary folder (tipically `~/go/bin`) should be added to your PATH.

### From release (x86_64 only)

Download a tagged release binary for your OS (ubuntu, macos, windows) placing it in a folder in your PATH and make it executable (may require elevated permissions).

### From source

Clone this repo and compile the source code:

```bash
git clone github.com/nu12/action-docs
cd action-docs
go build -o action-docs main.go
```

Move binary to a bin folder in your PATH (may require elevated permissions):
```bash
mv action-docs /usr/local/bin/
```

## Basic usage

```
Create documentation for github actions and workflows

Usage:
  action-docs [command]

Available Commands:
  actions     Generate documentation for github actions
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Show current version
  workflows   Generate documentation for github workflows

Flags:
      --config string   config file (default is $HOME/.action-docs.yaml)
  -h, --help            help for action-docs

Use "action-docs [command] --help" for more information about a command.
```