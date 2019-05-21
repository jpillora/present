//go:generate statik -dest ../ -f -p static -src files/

package static

import (
	"log"
	"net/http"
	"strings"

	"github.com/rakyll/statik/fs"
)

var cachedHFS http.FileSystem

//load on first use
func hfs() http.FileSystem {
	if cachedHFS != nil {
		return cachedHFS
	}
	f, err := fs.New()
	if err != nil {
		panic(err)
	}
	cachedHFS = f
	return f
}

func Handler() http.Handler {
	return http.FileServer(hfs())
}

func Read(file string) ([]byte, error) {
	if !strings.HasPrefix(file, "/") {
		file = "/" + file
	}
	return fs.ReadFile(hfs(), file)
}

func MustRead(file string) []byte {
	b, err := Read(file)
	if err != nil {
		log.Fatalf("static-fs read: %s: %s", file, err)
	}
	return b
}
