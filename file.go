package files

import (
	"errors"
	"os"
	"path/filepath"
)

func Open(name string, opt ...Option) (*os.File, error) {
	return OpenFile(name, os.O_RDONLY, opt...)
}

func Create(name string, opt ...Option) (*os.File, error) {
	return OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, opt...)
}

func OpenFile(name string, flag int, opt ...Option) (*os.File, error) {
	o := getOptions(opt...)
	f, err := os.OpenFile(name, flag, o.FilePerm)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) || flag&os.O_CREATE == 0 {
			return nil, err
		}
		if err = os.MkdirAll(filepath.Dir(name), o.DirPerm); err != nil {
			return nil, err
		}
		f, err = os.OpenFile(name, flag, o.FilePerm)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

func ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func WriteFile(name string, data []byte, opt ...Option) error {
	f, err := OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, opt...)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Sync(); err1 != nil && err == nil {
		err = err1
	}
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
