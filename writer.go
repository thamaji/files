package files

import (
	"os"
	"path/filepath"
)

// OpenFileWriterは、指定されたファイル名に対して安全に書き込むためのWriterを作成します。
// 書き込み中は一時ファイルにデータが書き込まれ、
// 書き込みが完了した後に指定された名前のファイルにデータが反映されます。
// オプションでファイルに関する設定を指定できます。
func OpenFileWriter(name string, opt ...Option) (*FileWriter, error) {
	f, err := CreateTemp(filepath.Dir(name), filepath.Base(name)+"_*", opt...)
	if err != nil {
		return nil, err
	}
	w := &FileWriter{name: name, f: f, err: nil}
	return w, nil
}

// FileWriter はファイルへの書き込みを管理する構造体です。
type FileWriter struct {
	name string
	f    *os.File
	err  error
}

// Write は指定されたバイトスライスをファイルに書き込みます。
// 書き込みに成功した場合は書き込まれたバイト数を返します。
func (w *FileWriter) Write(p []byte) (int, error) {
	n, err := w.f.Write(p)
	if err != nil {
		w.err = err
	}
	return n, err
}

// WriteAt は指定されたオフセット位置にデータをファイルに書き込みます。
// 書き込みに成功した場合は書き込まれたバイト数を返します。
func (w *FileWriter) WriteAt(b []byte, off int64) (int, error) {
	n, err := w.f.WriteAt(b, off)
	if err != nil {
		w.err = err
	}
	return n, err
}

// Close はファイルをクローズし、必要に応じてファイルの同期、名前の変更を行います。
// 書き込みエラーが発生した場合はそのエラーを返し、
// 最後にエラーが発生した場合は一時ファイルを削除します。
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
