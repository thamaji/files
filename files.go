package files

import (
	"errors"
	"io"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

// Remove は指定されたファイルを削除します。
func Remove(name string) error {
	return os.Remove(name)
}

// Exists は指定されたファイルまたはディレクトリが存在するかを確認します。
// 存在する場合は true を、存在しない場合は false を返します。
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

// MustExists は指定されたファイルまたはディレクトリが存在するかを確認します。
// 存在する場合は true を返し、存在しない場合は false を返します。
// エラーが発生した場合も false を返します。
func MustExists(name string) bool {
	ok, err := Exists(name)
	if err != nil {
		return false
	}
	return ok
}

// Rename は、指定された oldpath のファイルまたはディレクトリの名前を newpath に変更します。
// 移動先のパスが異なるデバイス上にある場合など、一部の環境では失敗することがあります。
func Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// Move はファイルまたはディレクトリを oldpath から newpath に移動します。
// まず os.Rename を試みますが、失敗した場合は Copy でコピーを行い、
// その後、元のファイルを削除します。
func Move(oldpath, newpath string) error {
	if err := os.Rename(oldpath, newpath); err == nil {
		return nil
	}

	if err := Copy(oldpath, newpath); err != nil {
		return err
	}

	return os.Remove(oldpath)
}

// Copy は、指定された srcpath から dstpath へファイルまたはディレクトリをコピーします。
// ファイルの場合は、そのまま内容をコピーし、
// ディレクトリの場合は中のファイル・ディレクトリを再帰的にコピーします。
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

// Walk は指定したディレクトリを再帰的に走査し、各ファイルのパスとファイル情報を返すイテレータを生成します。
// 例:
//
//	for path, info := range Walk("/some/path") {
//		fmt.Println(path, info.Size())
//	}
func Walk(name string) iter.Seq2[string, fs.FileInfo] {
	return func(yield func(string, fs.FileInfo) bool) {
		filepath.Walk(name, func(path string, info fs.FileInfo, err error) error {
			if !yield(path, info) {
				return filepath.SkipAll
			}
			return nil
		})
	}
}
