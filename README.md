# spaces

## Install
Assuming you have the standard go tools installed and `$GOPATH/bin` is in your `PATH`:

```bash
$ go get github.com/vikasgorur/spaces/showspaces github.com/vikasgorur/spaces/trimspaces
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
function commit {
    git add $(trimspaces -list-fixed -changed)
    git commit $*
}
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

## FAQ

### Why bother?

Trailing whitespaces are an annoyance. They can make pull requests and code reviews confusing
because a whitespace change causes a line to be a part of the diff even though nothing else
has changed. Most developers and open source projects prefer that they are removed.

For more, read [Why Are Trailing Whitespaces Bad?](http://www.dinduks.com/why-are-trailing-whitespaces-bad/)

### Why not a shell one-liner?

I could do this with a one-liner, but I'd rather have a tool that understands git and will work for
repositories of any language.

### Why not let the editor do this on save?

Configuring your editor to do this is a good practice. Sometimes, though, I'll use a different editor for a quick fix (vim)
or I need to fix files that someone else edited.

## Contributing

Please open an issue if you find a bug or contact me via email or twitter. Pull requests are welcome!

If you find the tool useful, tweet at me [@vikasgorur](https://twitter.com/vikasgorur).
