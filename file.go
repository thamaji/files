package files

import (
	"bufio"
	"errors"
	"iter"
	"os"
	"path/filepath"
)

// Open は指定されたファイルを読み取り専用で開きます。
// オプションでファイルに関する設定を指定できます。
func Open(name string, opt ...Option) (*os.File, error) {
	return OpenFile(name, os.O_RDONLY, opt...)
}

// Create は指定されたファイルを読み書きモードで開き、ファイルが存在しない場合は作成します。
// オプションでファイルに関する設定を指定できます。
func Create(name string, opt ...Option) (*os.File, error) {
	return OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, opt...)
}

// OpenFile は指定されたファイルを指定されたフラグで開きます。
// 必要に応じてディレクトリを作成し、ファイルを開き直します。
// オプションでファイルに関する設定を指定できます。
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

// ReadFile は指定されたファイルを読み込み、その内容をバイトスライスとして返します。
func ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// WriteFile は指定されたファイルにデータを書き込みます。
// ファイルが存在しない場合は新しく作成し、既存の内容は上書きされます。
// オプションでファイルに関する設定を指定できます。
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

// ReadLine は指定されたファイルを1行ずつ読み込むイテレータを返します。
// 例:
//
//	for line := range ReadLine("example.txt") {
//		fmt.Println(line)
//	}
func ReadLine(name string) iter.Seq[string] {
	return func(yield func(string) bool) {
		f, err := os.OpenFile(name, os.O_RDONLY, 0)
		if err != nil {
			return
		}
		s := bufio.NewScanner(f)
		defer f.Close()
		for s.Scan() {
			if !yield(s.Text()) {
				return
			}
		}
	}
}
