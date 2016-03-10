# spaces

## Install
Assuming you have the standard go tools installed and `$GOPATH/bin` is in your `PATH`:

```bash
$ go get github.com/vikasgorur/spaces/trimspaces github.com/vikasgorur/spaces/trimspaces
```

## trimspaces
A tool to remove trailing whitespaces recursively from source files in a directory.

## Usage

```
Usage: ./trimspaces [-check] [-dir|-changed] [file1 ...]

Fix trailing spaces in input files (or stdin).

  -changed
    	operate only on files that have been changed (only works in git repos).
  -dir
    	operate recursively on all source files in the current directory.
  -verbose
    	run in verbose mode
```

Input must be provided either by the mode (`dir/changed`) or by a list of files on the command-line,
or by standard input.

To fix all files under the current directory:

```bash
$ trimspaces -dir
```

To fix only changed files:

```bash
$ trimspaces -changed
```

To fix a given list of files:
```bash
$ trimspaces src/*.js
```

`trimspaces` is designed to be used with git repositories containing source code. It respects `.gitignore` and only
modifies source files (files with extensions it recognizes). I use it as part of an alias to
clean up whitespaces before every commit:

```bash
alias commit="trimspaces -changed; git commit $*"
```

## showspaces

`showspaces` is a tool to highlight trailing spaces in files. It accepts the same arguments as `trimspaces` with one addition,
a 'check' mode. In this mode no output is printed but the exit status is non-zero if any of the input files contain trailing
spaces.

To highlight spaces in all changed files:

```bash
$ showspaces -changed
```

To check a given list of files:
```bash
$ showspaces -check src/*.js
```


## Why bother?

Trailing whitespaces are an annoyance. They can make pull requests and code reviews confusing
because a whitespace change causes a line to be a part of the diff even though nothing else
has changed. Most developers and open source projects prefer that they are removed.

For more, read [Why Are Trailing Whitespaces Bad?](http://www.dinduks.com/why-are-trailing-whitespaces-bad/)

## Contributing

Pull requests are welcome!

If you find the tool useful, tweet at me [@vikasgorur](https://twitter.com/vikasgorur).
