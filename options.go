package files

import "io/fs"

var DefaultFilePerm fs.FileMode = 0644
var DefaultDirPerm fs.FileMode = 0755

func SetDefaultFileMode(perm fs.FileMode) {
	DefaultDirPerm = perm
}

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

type Option func(options) options

func WithFilePerm(perm fs.FileMode) Option {
	return func(o options) options {
		o.FilePerm = perm
		return o
	}
}

func WithDirPerm(perm fs.FileMode) Option {
	return func(o options) options {
		o.FilePerm = perm
		return o
	}
}
