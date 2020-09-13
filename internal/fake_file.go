package internal

import (
	"os"
	"time"
)

type fakeFile struct {
	name    string
	size    int64
	modTime time.Time
	isDir   bool
}

func (f fakeFile) Name() string {
	return f.name
}

func (f fakeFile) Size() int64 {
	return f.size
}

func (f fakeFile) Mode() os.FileMode {
	panic("implement me")
}

func (f fakeFile) ModTime() time.Time {
	return f.modTime
}

func (f fakeFile) IsDir() bool {
	return f.isDir
}

func (f fakeFile) Sys() interface{} {
	panic("implement me")
}
