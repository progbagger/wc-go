# wc

`wc` is a command-line Linux utility that counts things in files.

## How to build?

Just proceed to the **src** directory and execute `make` command.

If there is an error like `package <name> is not in std` then execute from project root:

```shell
go work init; go work use src/args src/wc
```

## How to use?

There are a few flags in this utility. If one of them is specified then others can't.

- `-m` - for counting runes (**UTF-8** encoded symbols)
- `-l` - for counting lines
- `-w` - for counting words (default)
- `--help` - prints an information about flags

Rest of the arguments are treated as files to process.

This program uses *goroutines* so on each execution order of results can be different.

## Example

Launched from **src** directory after executing `make`

```shell
-> % ../build/wc -m ../README.md ../LICENSE ../.gitignore
478     ../.gitignore
734     ../README.md
1071    ../LICENSE
```
