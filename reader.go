package files

import (
	"os"
)

func OpenFileReader(name string) (*FileReader, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	r := &FileReader{f: f}
	return r, nil
}

type FileReader struct {
	f *os.File
}

func (r *FileReader) Read(p []byte) (int, error) {
	return r.f.Read(p)
}

func (r *FileReader) ReadAt(p []byte, off int64) (int, error) {
	return r.f.ReadAt(p, off)
}

func (r *FileReader) Close() error {
	return r.f.Close()
}
