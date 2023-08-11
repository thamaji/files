package files

import (
	"os"
	"path/filepath"
)

func OpenFileWriter(name string, opt ...Option) (*FileWriter, error) {
	f, err := CreateTemp(filepath.Dir(name), filepath.Base(name)+"_*", opt...)
	if err != nil {
		return nil, err
	}
	w := &FileWriter{name: name, f: f, err: nil}
	return w, nil
}

type FileWriter struct {
	name string
	f    *os.File
	err  error
}

func (w *FileWriter) Write(p []byte) (int, error) {
	n, err := w.f.Write(p)
	if err != nil {
		w.err = err
	}
	return n, err
}

func (w *FileWriter) WriteAt(b []byte, off int64) (int, error) {
	n, err := w.f.WriteAt(b, off)
	if err != nil {
		w.err = err
	}
	return n, err
}

func (w *FileWriter) Close() error {
	err := w.err
	if err1 := w.f.Sync(); err1 != nil && err == nil {
		err = err1
	}
	if err1 := w.f.Close(); err1 != nil && err == nil {
		err = err1
	}
	if err != nil {
		_ = os.Remove(w.f.Name())
		return err
	}
	if err := os.Rename(w.f.Name(), w.name); err != nil {
		_ = os.Remove(w.f.Name())
		return err
	}
	return nil
}
