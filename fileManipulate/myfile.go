package filemanipulate

import "os"

type MyFile struct {
	name string
	buff []byte
}

func NewFile(name string) *MyFile {
	return &MyFile{name: name}
}

func (f *MyFile) Name() string {
	return f.name
}

func (f *MyFile) Buff() []byte {
	return f.buff
}

func (f *MyFile)Read() error {
	buff, err := os.ReadFile(f.name)
	f.buff = buff
	if err != nil {
		return err
	}
	return nil
}