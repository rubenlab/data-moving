package main

import (
	"io"
	"io/fs"
	iofs "io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DirFsBase struct {
	Path string
}

func (fs *DirFsBase) abspath(path string) string {
	return filepath.Join(fs.Path, path)
}

type LocalDirFs struct {
	DirFsBase
}

func (fs *LocalDirFs) Walk(enterDir WalkDirFunc, enterFile WalkFunc, exitDir WalkDirFunc) {
	fs.walkDir(".", 1, enterDir, enterFile, exitDir)
}

func (fs *LocalDirFs) walkDir(dir string, level int, enterDir WalkDirFunc, enterFile WalkFunc, exitDir WalkDirFunc) {
	dirpath := fs.abspath(dir)
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		relpath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			enterDir(relpath, iofs.FileInfoToDirEntry(file), level, nil)
			fs.walkDir(relpath, level+1, enterDir, enterFile, exitDir)
			exitDir(relpath, iofs.FileInfoToDirEntry(file), level, nil)
		} else {
			enterFile(relpath, file, level, nil)
		}
	}
}

func (fs *LocalDirFs) MkdirAll(path string) error {
	abspath := fs.abspath(path)
	return os.MkdirAll(abspath, FileFileMode)
}

func (fs *LocalDirFs) MkdirAllAbs(rootpath string, relpath string) error {
	abspath := filepath.Join(rootpath, relpath)
	return os.MkdirAll(abspath, DirFileMode)
}

func (fs *LocalDirFs) Open(path string) (io.ReadWriteCloser, error) {
	abspath := fs.abspath(path)
	return os.Open(abspath)
}

func (fs *LocalDirFs) Create(path string) (io.ReadWriteCloser, error) {
	abspath := fs.abspath(path)
	return os.Create(abspath)
}

func (fs *LocalDirFs) OpenFile(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
	abspath := fs.abspath(name)
	return os.OpenFile(abspath, flag, perm)
}

func (fs *LocalDirFs) Chmod(path string, mode os.FileMode) error {
	abspath := fs.abspath(path)
	return os.Chmod(abspath, mode)
}

func (fs *LocalDirFs) Chown(path string, uid, gid int) error {
	return os.Chown(fs.abspath(path), uid, gid)
}

func (fs *LocalDirFs) Remove(path string) error {
	return os.Remove(fs.abspath(path))
}

func (fs *LocalDirFs) IsEmptyDir(path string) (bool, error) {
	f, err := os.Open(fs.abspath(path))
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func (fs *LocalDirFs) Move(path string, dest string) error {
	destpath := filepath.Join(dest, path)
	srcpath := fs.abspath(path)
	return os.Rename(srcpath, destpath)
}

func (fs *LocalDirFs) Lstat(p string) (os.FileInfo, error) {
	return os.Lstat(fs.abspath(p))
}

var _ DirFs = (*LocalDirFs)(nil)

type LocalDirFsCreator struct {
	fs *LocalDirFs
}

func (c *LocalDirFsCreator) create() (DirFs, error) {
	return c.fs, nil
}

func (c *LocalDirFsCreator) close() {

}

var _ DirFsCreator = (*LocalDirFsCreator)(nil)

func createLocalDirFsCreator(config DirFsConfig) (DirFsCreator, error) {
	fs := LocalDirFs{
		DirFsBase: DirFsBase{
			Path: config.Path,
		},
	}
	return &LocalDirFsCreator{
		fs: &fs,
	}, nil
}

func init() {
	registCreatorFactory("local", createLocalDirFsCreator)
}
