package files

import (
	"io"
	"io/fs"
	"os"
)

func Mkdir(name string, opt ...Option) error {
	o := getOptions(opt...)
	return os.Mkdir(name, o.DirPerm)
}

func MkdirAll(path string, opt ...Option) error {
	o := getOptions(opt...)
	return os.MkdirAll(path, o.DirPerm)
}

func Readdir(name string) ([]fs.FileInfo, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	_ = f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func MustReaddir(name string) []fs.FileInfo {
	list, err := Readdir(name)
	if err != nil {
		return nil
	}
	return list
}

func ReadDirnames(name string) ([]string, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdirnames(-1)
	_ = f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func MustReadDirnames(name string) []string {
	list, err := ReadDirnames(name)
	if err != nil {
		return nil
	}
	return list
}

func ReadDir(name string) ([]fs.DirEntry, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	list, err := f.ReadDir(-1)
	_ = f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func MustReadDir(name string) []fs.DirEntry {
	list, err := ReadDir(name)
	if err != nil {
		return nil
	}
	return list
}

func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func IsDir(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return true, err
	}
	return fi.IsDir(), nil
}

func MustIsDir(name string) bool {
	ok, err := IsDir(name)
	if err != nil {
		return false
	}
	return ok
}

func IsEmptyDir(name string) (bool, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return false, err
	}
	_, err = f.Readdirnames(1)
	_ = f.Close()
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func MustIsEmptyDir(name string) bool {
	ok, err := IsEmptyDir(name)
	if err != nil {
		return false
	}
	return ok
}
