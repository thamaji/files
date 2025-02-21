package files

import (
	"errors"
	"os"
)

// CreateTemp は指定されたディレクトリに一時ファイルを作成します。
// ディレクトリが存在しない場合は作成し、指定されたパターンに基づいて一時ファイルを作成します。
// オプションでファイルに関する設定を指定できます。
func CreateTemp(dir string, pattern string, opt ...Option) (*os.File, error) {
	o := getOptions(opt...)
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		if err = os.MkdirAll(dir, o.DirPerm); err != nil {
			return nil, err
		}
		f, err = os.CreateTemp(dir, pattern)
		if err != nil {
			return nil, err
		}
	}
	if err := f.Chmod(o.FilePerm); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return nil, err
	}
	return f, nil
}

// MkdirTemp は指定されたディレクトリに一時ディレクトリを作成します。
// 作成したディレクトリに指定されたパーミッションを設定します。
// オプションでファイルに関する設定を指定できます。
func MkdirTemp(dir string, pattern string, opt ...Option) (string, error) {
	o := getOptions(opt...)
	name, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return "", err
	}
	if err := os.Chmod(name, o.DirPerm); err != nil {
		_ = os.Remove(name)
		return "", err
	}
	return name, nil
}
