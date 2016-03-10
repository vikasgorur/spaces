# trimspaces
A tool to remove trailing whitespaces recursively from source files in a directory.

## Usage

The tool operates in two modes: all source files in the directory, or just those files that
have been changed (only works inside git repositories).

To fix all files under the current directory:

```bash
$ trimspaces -dir
```

To fix only changed files:

```bash
$ trimspaces -changed
```

`trimspaces` is designed to be used with git repositories containing source code. It respects `.gitignore` and only
modifies source files (files with extensions it recognizes). I use it as part of an alias to
clean up whitespaces before every commit:

```bash
alias commit="trimspaces -changed; git commit $*"
```

## Why bother?

Trailing whitespaces are an annoyance. They can make pull requests and code reviews confusing
because a whitespace change causes a line to be a part of the diff even though nothing else
has changed. Most developers and open source projects prefer that they are removed.

For more, read [Why Are Trailing Whitespaces Bad?](http://www.dinduks.com/why-are-trailing-whitespaces-bad/)

## Contributing

Pull requests are welcome!

If you find the tool useful, tweet at me [@vikasgorur](https://twitter.com/vikasgorur).