margo
=====

`margo` (the Marvelous Go tutorial) is a simple repository holding sources for an introduction to `Go`.

The slides are here:

- [part 1 (basic intro)](https://talks.godoc.org/github.com/sbinet/margo/part-1/talk.slide)
- [part 2 (concurrency)](https://talks.godoc.org/github.com/sbinet/margo/part-2/talk.slide)
- [part 3 (cgo)](https://talks.godoc.org/github.com/sbinet/margo/part-3/talk.slide)
- [part 4 (web servers)](https://talks.godoc.org/github.com/sbinet/margo/part-4/talk.slide)
- [part 5 (a DAQ application)](https://talks.godoc.org/github.com/sbinet/margo/part-5/talk.slide)
- [part 6 (a chat application)](https://talks.godoc.org/github.com/sbinet/margo/part-6/talk.slide)
- [part 7 (biology)](https://talks.godoc.org/github.com/sbinet/margo/part-7/talk.slide)
- [part 8 (modules)](https://talks.godoc.org/github.com/sbinet/margo/part-8/talk.slide)

## Bootstrapping the work environment

### Installing the `Go` toolchain

The `Go` hands-on session obviously needs a working `Go` toolchain.

There are 3 ways to achieve this:
- install `Go` via your favorite package manager (`yum`, `apt-get`, `fink`, ...)
- install `Go` via `docker`
- install `Go` manually.

Installing and configuring `Go` can be achieved by following the instructions of the official `Go` installation page: https://golang.org/doc/install

### Setting up the work environment

Like `python` and its `$PYTHONPATH` environment variable, `Go` uses
`$GOPATH` to locate packages' source trees.
You can choose whatever you like (obviously a directory under which
you have read/write access, though.)

In the following, we'll assume you chose `$HOME/go`:

```sh
$ mkdir -p $HOME/go
$ export GOPATH=$HOME/go
$ export PATH=$GOPATH/bin:$PATH
```

Make sure the `go` tool is correctly setup:

```sh
$ go env
GOARCH="amd64"
GOBIN=""
GOCHAR="6"
GOEXE=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOOS="linux"
GOPATH="$HOME/go"
GORACE=""
GOROOT="/usr/lib/golang"
GOTOOLDIR="$HOME/go/pkg/tool/linux_amd64"
CC="gcc"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0"
CXX="g++"
CGO_ENABLED="1"
```

(on other platforms/architectures, the output might differ
slightly. The important env.vars. are `GOPATH` and `GOROOT`.)

### Testing `go`

Now that the `go` tool is correctly setup, let's try to fetch some
code.
For this part, you'll need the following tools installed to actually retrieve the code from the repositories:

- `git`

Without further ado:

```sh
$ go install -v github.com/sbinet/margo/cmd/margo-hello@latest
go: downloading github.com/sbinet/margo v0.0.0-20200115082936-ccffedad4555
github.com/sbinet/margo/cmd/margo-hello
```

`go install` downloaded (cloned, in `git` speak) the whole
`github.com/sbinet/margo` repository (under `$GOPATH/pkg/mod`) and
compiled the `margo-hello` command.
As the compilation was successful, it also installed the `margo-hello`
command under `$GOPATH/bin`.

The `margo-hello` command is now available from your shell:

```sh
$ margo-hello
Hello MarGo!

$ margo-hello you
Hello you!
```

## Setting up your favorite editor

Extensive documentation on how to setup your editor (for code
highlighting, code completion, ...) is available here:

 https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins
 
At the very least, you should try to install and setup `goimports` as
explained here:

 https://godoc.org/golang.org/x/tools/cmd/goimports

`goimports` provides automatic code formating as well as automated
insertion/deletion of used/unused packages (in your `import` package
statements.)

## Documentation

The `Go` programming language is quite new (released in 2009) but
ships already with quite a fair amount of documentation.
Here are a few pointers:

- https://golang.org/doc/code.html
- https://tour.golang.org
- https://golang.org/doc/effective_go.html
- https://dave.cheney.net/resources-for-new-go-programmers
- https://gobyexample.com
- https://greatercommons.com/cwg
- https://learn-golang.com/

For more advanced topics:

- https://talks.golang.org
- https://blog.golang.org
- https://groups.google.com/d/forum/golang-nuts (`Go` community forum)
- the internets. typical queries are `go something` or `golang something`
