package files

import (
	"os"
)

// OpenFileReader は、指定されたファイルを読み込み専用で開き、
// FileReader インスタンスを返します。
func OpenFileReader(name string) (*FileReader, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	r := &FileReader{f: f}
	return r, nil
}

// FileReader はファイルの読み取り操作をラップする構造体です。
// Read、ReadAt、Close メソッドを提供します。
type FileReader struct {
	f *os.File
}

// Read は、ファイルからデータを読み取ります。
// 引数として渡されたバイトスライスに読み取ったデータを格納し、
// 実際に読み取ったバイト数とエラーを返します。
func (r *FileReader) Read(p []byte) (int, error) {
	return r.f.Read(p)
}

// ReadAt は、指定されたオフセット位置からデータを読み取ります。
// 引数として渡されたオフセットからデータを読み取り、
// 実際に読み取ったバイト数とエラーを返します。
func (r *FileReader) ReadAt(p []byte, off int64) (int, error) {
	return r.f.ReadAt(p, off)
}

// Close は、ファイルを閉じます。
// ファイルを正常に閉じた場合は nil を返し、エラーが発生した場合はそのエラーを返します。
func (r *FileReader) Close() error {
	return r.f.Close()
}
