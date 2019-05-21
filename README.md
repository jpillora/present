# present

[Go present tool](https://godoc.org/golang.org/x/tools/present), forked with the following changes

* Go code blocks are syntax highlighted
* Importable by other packages
* All static files are now compiled in using statik
  * This now means the binary is stand-alone, and it does not require a GOPATH