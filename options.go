package files

import "io/fs"

var DefaultFilePerm fs.FileMode = 0644 // ファイルパーミッションのデフォルト値（0644）
var DefaultDirPerm fs.FileMode = 0755  // ディレクトリパーミッションのデフォルト値（0755）

// SetDefaultFileMode は、デフォルトのファイルパーミッションを設定します。
// この関数を呼び出すことで、ファイルの作成時に使用されるデフォルトのパーミッションが変更されます。
func SetDefaultFileMode(perm fs.FileMode) {
	DefaultDirPerm = perm
}

// SetDefaultDirFileMode は、デフォルトのディレクトリパーミッションを設定します。
// この関数を呼び出すことで、ディレクトリの作成時に使用されるデフォルトのパーミッションが変更されます。
func SetDefaultDirFileMode(perm fs.FileMode) {
	DefaultDirPerm = perm
}

type options struct {
	FilePerm fs.FileMode
	DirPerm  fs.FileMode
}

func getOptions(opts ...Option) options {
	o := options{
		FilePerm: DefaultFilePerm,
		DirPerm:  DefaultDirPerm,
	}
	for _, of := range opts {
		o = of(o)
	}
	return o
}

// Option は、設定を変更する関数の型です。
type Option func(options) options

// WithFilePerm は、ファイルのパーミッションを設定するオプションを返します。
// このオプションは、ファイルの作成時に使用されるパーミッションを指定します。
func WithFilePerm(perm fs.FileMode) Option {
	return func(o options) options {
		o.FilePerm = perm
		return o
	}
}

// WithDirPerm は、ディレクトリのパーミッションを設定するオプションを返します。
// このオプションは、ディレクトリの作成時に使用されるパーミッションを指定します。
func WithDirPerm(perm fs.FileMode) Option {
	return func(o options) options {
		o.FilePerm = perm
		return o
	}
}
