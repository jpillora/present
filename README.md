# present

[Go present tool](https://godoc.org/golang.org/x/tools/present/cmd/present), forked with the following changes

* Go code blocks are syntax highlighted
* Importable by other packages
* All static files are now compiled in using statik
  * This now means the binary is stand-alone, and it does not require a GOPATH

## Usage

```sh
$ go get -v github.com/jpillora/present
$ present --help
```