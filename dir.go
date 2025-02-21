package files

import (
	"io"
	"io/fs"
	"os"
)

// Mkdir は指定した名前のディレクトリを作成します。
// オプションでディレクトリのパーミッションを指定できます。
func Mkdir(name string, opt ...Option) error {
	o := getOptions(opt...)
	return os.Mkdir(name, o.DirPerm)
}

// MkdirAll は指定したパスに必要なディレクトリを再帰的に作成します。
// オプションでディレクトリのパーミッションを指定できます。
func MkdirAll(path string, opt ...Option) error {
	o := getOptions(opt...)
	return os.MkdirAll(path, o.DirPerm)
}

// Readdir は指定されたディレクトリ内のファイル情報を一覧で取得します。
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

// MustReaddir は指定されたディレクトリ内のファイル情報を一覧で取得します。
// エラーが発生した場合、空のスライスを返します。
func MustReaddir(name string) []fs.FileInfo {
	list, err := Readdir(name)
	if err != nil {
		return nil
	}
	return list
}

// ReadDirnames は指定されたディレクトリ内のファイル名のみを一覧で取得します。
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

// MustReadDirnames は指定されたディレクトリ内のファイル名のみを一覧で取得します。
// エラーが発生した場合、空のスライスを返します。
func MustReadDirnames(name string) []string {
	list, err := ReadDirnames(name)
	if err != nil {
		return nil
	}
	return list
}

// ReadDir は指定されたディレクトリ内のエントリ（ディレクトリやファイル）の情報を一覧で取得します。
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

// MustReadDir は指定されたディレクトリ内のエントリ（ディレクトリやファイル）の情報を一覧で取得します。
// エラーが発生した場合、空のスライスを返します。
func MustReadDir(name string) []fs.DirEntry {
	list, err := ReadDir(name)
	if err != nil {
		return nil
	}
	return list
}

// RemoveAll は指定されたパスを削除します。
// パスがディレクトリの場合、再帰的に中身も削除されます。
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// IsDir は指定された名前がディレクトリかどうかを確認します。
func IsDir(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return true, err
	}
	return fi.IsDir(), nil
}

// MustIsDir は指定された名前がディレクトリかどうかを確認します。
// エラーが発生した場合、false を返します。
func MustIsDir(name string) bool {
	ok, err := IsDir(name)
	if err != nil {
		return false
	}
	return ok
}

// IsEmptyDir は指定されたディレクトリが空かどうかを確認します。
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

// MustIsEmptyDir は指定されたディレクトリが空かどうかを確認します。
// エラーが発生した場合、false を返します。
func MustIsEmptyDir(name string) bool {
	ok, err := IsEmptyDir(name)
	if err != nil {
		return false
	}
	return ok
}
