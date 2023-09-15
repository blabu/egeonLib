package golang

import (
	"io"
	"os"
)

// ReadFile - read all file from path
func ReadFile(Path string) ([]byte, error) {
	f, e := os.Open(Path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return io.ReadAll(f)
}
