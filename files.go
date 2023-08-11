package files

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

func Remove(name string) error {
	return os.Remove(name)
}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func MustExists(name string) bool {
	ok, err := Exists(name)
	if err != nil {
		return false
	}
	return ok
}

func Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func Move(oldpath, newpath string) error {
	if err := os.Rename(oldpath, newpath); err == nil {
		return nil
	}

	if err := Copy(oldpath, newpath); err != nil {
		return err
	}

	return os.Remove(oldpath)
}

func Copy(srcpath, dstpath string) error {
	sf, err := os.OpenFile(srcpath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	sfi, err := sf.Stat()
	if err != nil {
		_ = sf.Close()
		return err
	}

	if sfi.IsDir() {
		temp, err := MkdirTemp(filepath.Dir(dstpath), filepath.Base(dstpath)+"_*", WithDirPerm(sfi.Mode().Perm()))
		if err != nil {
			return err
		}

		names, err := sf.Readdirnames(-1)
		_ = sf.Close()
		if err != nil {
			return err
		}

		for _, name := range names {
			if err := copy(filepath.Join(srcpath, name), filepath.Join(temp, name)); err != nil {
				_ = os.RemoveAll(temp)
				return err
			}
		}

		if err := os.Rename(temp, dstpath); err != nil {
			_ = os.RemoveAll(temp)
			return err
		}

		return nil
	}

	df, err := OpenFileWriter(dstpath, WithFilePerm(sfi.Mode().Perm()))
	if err != nil {
		return err
	}
	_, err = io.Copy(df, sf)
	if err1 := df.Close(); err1 != nil && err == nil {
		err = err1
	}
	if err != nil {
		return err
	}

	return nil
}

func copy(srcpath, dstpath string) error {
	sf, err := os.OpenFile(srcpath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	sfi, err := sf.Stat()
	if err != nil {
		_ = sf.Close()
		return err
	}

	if sfi.IsDir() {
		if err := os.Mkdir(dstpath, sfi.Mode().Perm()); err != nil {
			_ = sf.Close()
			return err
		}

		names, err := sf.Readdirnames(-1)
		_ = sf.Close()
		if err != nil {
			return err
		}
		for _, name := range names {
			if err := copy(filepath.Join(srcpath, name), filepath.Join(dstpath, name)); err != nil {
				return err
			}
		}
		return nil
	}

	df, err := os.OpenFile(dstpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sfi.Mode().Perm())
	if err != nil {
		_ = sf.Close()
		return err
	}
	_, err = io.Copy(df, sf)
	_ = sf.Close()
	if err1 := df.Sync(); err1 != nil && err == nil {
		err = err1
	}
	if err1 := df.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
