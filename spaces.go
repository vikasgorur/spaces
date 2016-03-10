package spaces

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/svent/sift/gitignore"
	"gopkg.in/fatih/set.v0"
)

// list stolen from https://github.com/ggreer/the_silver_searcher/blob/master/src/lang.c
var srcExtensions set.Interface = set.New(
	"as", "mxml", // actionscript
	"ada", "adb", "ads", // ada
	"asm", "s", // asm
	"bat", "cmd", // batch
	"bb", "bbappend", "bbclass", "inc", // bitbake
	"bro", "bif", // bro
	"c", "h", "xs", // cc
	"cfc", "cfm", "cfml", // cfmx
	"clj", "cljs", "cljc", "cljx", // clojure
	"coffee", "cjsx", // coffee
	"cpp", "cc", "C", "cxx", "m", "hpp", "hh", "h", "H", "hxx", // cpp
	"cr", "ecr", // crystal
	"cs",                                                                                  // csharp
	"css",                                                                                 // css
	"pas", "int", "dfm", "nfm", "dof", "dpk", "dproj", "groupproj", "bdsgroup", "bdsproj", // delphi
	"ebuild", "eclass", // ebuild
	"el",        // elisp
	"ex", "exs", // elixir
	"erl", "hrl", // erlang
	"f", "f77", "f90", "f95", "f03", "for", "ftn", "fpp", // fortran
	"fs", "fsi", "fsx", // fsharp
	"po", "pot", "mo", // gettext
	"go",                               // go
	"groovy", "gtmpl", "gpp", "grunit", // groovy
	"haml",      // haml
	"hs", "lhs", // haskell
	"h",                             // hh
	"htm", "html", "shtml", "xhtml", // html
	"ini",                // ini
	"jade",               // jade
	"java", "properties", // java
	"js", "jsx", // js
	"json",                         // json
	"jsp", "jspx", "jhtm", "jhtml", // jsp
	"jl",          // julia
	"less",        // less
	"liquid",      // liquid
	"lisp", "lsp", // lisp
	"lua",                    // lua
	"m4",                     // m4
	"Makefiles", "mk", "mak", // make
	"mako",                                           // mako
	"markdown", "mdown", "mdwn", "mkdn", "mkd", "md", // markdown
	"mas", "mhtml", "mpl", "mtxt", // mason
	"m",       // matlab
	"m", "wl", // mathematica
	"m", "moo", // mercury
	"nim",    // nim
	"m", "h", // objc
	"mm", "h", // objcpp
	"ml", "mli", "mll", "mly", // ocaml
	"m",                                            // octave
	"pir", "pasm", "pmc", "ops", "pod", "pg", "tg", // parrot
	"pl", "pm", "pm6", "pod", "t", // perl
	"php", "phpt", "php3", "php4", "php5", "phtml", // php
	"pike", "pmod", // pike
	"pt", "cpt", "metadata", "cpy", "py", "xml", "zcml", // plone
	"pp",               // puppet
	"py",               // python
	"qml",              // qml
	"rkt", "ss", "scm", // racket
	"Rakefiles",                       // rake
	"rst",                             // restructuredtext
	"rs",                              // rs
	"R", "Rmd", "Rnw", "Rtex", "Rrst", // r
	"rdoc",                                              // rdoc
	"rb", "rhtml", "rjs", "rxml", "erb", "rake", "spec", // ruby
	"rs",           // rust
	"sls",          // salt
	"sass", "scss", // sass
	"scala",     // scala
	"scm", "ss", // scheme
	"sh", "bash", "csh", "tcsh", "ksh", "zsh", // shell
	"st",                       // smalltalk
	"sml", "fun", "mlb", "sig", // sml
	"sql", "ctl", // sql
	"styl",               // stylus
	"swift",              // swift
	"tcl", "itcl", "itk", // tcl
	"tex", "cls", "sty", // tex
	"tt", "tt2", "ttml", // tt
	"toml",      // toml
	"ts", "tsx", // ts
	"vala", "vapi", // vala
	"bas", "cls", "frm", "ctl", "vb", "resx", // vb
	"vm", "vtl", "vsl", // velocity
	"v", "vh", "sv", // verilog
	"vhd", "vhdl", // vhdl
	"vim",        // vim
	"wxi", "wxs", // wix
	"wsdl",                             // wsdl
	"wadl",                             // wadl
	"xml", "dtd", "xsl", "xslt", "ent", // xml
	"yaml", "yml", // yaml
)

var ignoreChecker *gitignore.Checker

// IsIgnored returns true if the path matches .gitignore.
func IsIgnored(path string, info os.FileInfo) bool {
	if ignoreChecker == nil {
		ignoreChecker = gitignore.NewChecker()
		ignoreChecker.LoadBasePath(".")
	}

	return ignoreChecker.Check(path, info)
}

// IsSourceFile returns true if the path is a source file.
func IsSourceFile(path string, info os.FileInfo) bool {
	ext := filepath.Ext(path)
	if info.Mode().IsRegular() && strings.HasPrefix(ext, ".") && srcExtensions.Has(filepath.Ext(path)[1:]) {
		return true
	}

	return false
}

// WalkFunc is a function that is called on every input file.
// boolean return value indicates whether the file was affected.
type WalkFunc func(path string, info os.FileInfo, err error) (bool, error)

// WalkDir walks every file under the cwd.
// returns a list of affected files.
func WalkDir(fn WalkFunc) []string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "trimspaces: could not get current directory: %v\n", err)
		os.Exit(2)
	}

	var affectedFiles []string
	walker := func(path string, info os.FileInfo, err error) error {
		affected, e := fn(path, info, err)
		if affected {
			affectedFiles = append(affectedFiles, path)
		}
		return e
	}
	filepath.Walk(cwd, walker)
	return affectedFiles
}

// extractPath returns a path contained in a line of 'git status' output
// if it's a path we're interested in (added, modified, ...)
// returns nil otherwise
func extractPath(line string) string {
	pieces := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(pieces) < 2 {
		return ""
	}

	code := pieces[0]
	if code == "M" || code == "A" || code == "??" {
		return strings.TrimSpace(pieces[1])
	} else if []rune(code)[0] == 'R' {
		return strings.TrimSpace(strings.SplitN(pieces[1], "->", 2)[1])
	}
	//TODO: this will still fail if the filenames contain leading or trailing spaces
	// use 'git status -z' to handle it properly

	return ""
}

// IsGitRoot returns true if the cwd is the root of a git repo.
func IsGitRoot() bool {
	info, err := os.Stat(".git")
	if err == nil && info.IsDir() {
		return true
	}

	return false
}

// changedFiles returns a slice of file names that have been modified/added
// to the git repository.
func changedFiles() ([]string, error) {
	if !IsGitRoot() {
		return nil, errors.New("cwd is not the root of a git repository")
	}

	status := exec.Command("git", "status", "--porcelain")
	output, err := status.Output()
	if err != nil {
		return nil, err
	}

	var paths []string
	lines := bufio.NewScanner(bytes.NewReader(output))
	for lines.Scan() {
		path := extractPath(lines.Text())
		if path != "" {
			paths = append(paths, path)
		}
	}

	return paths, nil
}

// WalkChanged walks only changed/added files in the git repository.
// returns a list of affected files.
func WalkChanged(fn WalkFunc) []string {
	var affectedFiles []string
	paths, err := changedFiles()
	if err != nil {
		fmt.Println(err)
		return affectedFiles
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			return affectedFiles
		}
		affected, _ := fn(path, info, err)
		if affected {
			affectedFiles = append(affectedFiles, path)
		}
	}

	return affectedFiles
}

// WalkList walks a given list (slice) of files.
// returns a list of affected files.
func WalkList(files []string, fn WalkFunc) []string {
	var affectedFiles []string

	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			return affectedFiles
		}

		affected, _ := fn(path, info, err)
		if affected {
			affectedFiles = append(affectedFiles, path)
		}
	}

	return affectedFiles
}
