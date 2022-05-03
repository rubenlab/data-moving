package main

import (
	"io"
	"io/fs"
	iofs "io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
)

type SftpDirFs struct {
	DirFsBase
	client *sftp.Client
}

func (fs *SftpDirFs) Walk(enterDir WalkDirFunc, enterFile WalkFunc, exitDir WalkDirFunc) {
	fs.walkDir(".", 1, enterDir, enterFile, exitDir)
}

func (fs *SftpDirFs) walkDir(dir string, level int, enterDir WalkDirFunc, enterFile WalkFunc, exitDir WalkDirFunc) {
	dirpath := fs.abspath(dir)
	files, err := fs.client.ReadDir(dirpath)
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

func (fs *SftpDirFs) MkdirAll(path string) error {
	abspath := fs.abspath(path)
	err := fs.client.MkdirAll(abspath)
	if err != nil {
		return err
	}
	fs.client.Chmod(abspath, DirFileMode)
	return nil
}

func (fs *SftpDirFs) MkdirAllAbs(rootpath string, relpath string) error {
	abspath := filepath.Join(rootpath, relpath)
	err := fs.client.MkdirAll(abspath)
	if err != nil {
		return err
	}
	fs.client.Chmod(abspath, DirFileMode)
	return nil
}

func (fs *SftpDirFs) Open(path string) (io.ReadWriteCloser, error) {
	abspath := fs.abspath(path)
	return fs.client.Open(abspath)
}

func (fs *SftpDirFs) Create(path string) (io.ReadWriteCloser, error) {
	abspath := fs.abspath(path)
	return fs.client.Create(abspath)
}

func (fs *SftpDirFs) OpenFile(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
	abspath := fs.abspath(name)
	return fs.client.OpenFile(abspath, flag)
}

func (fs *SftpDirFs) Chmod(path string, mode os.FileMode) error {
	abspath := fs.abspath(path)
	return fs.client.Chmod(abspath, mode)
}

func (fs *SftpDirFs) Chown(path string, uid, gid int) error {
	abspath := fs.abspath(path)
	return fs.client.Chown(abspath, uid, gid)
}

func (fs *SftpDirFs) Remove(path string) error {
	abspath := fs.abspath(path)
	return fs.client.Remove(abspath)
}

func (fs *SftpDirFs) Move(path string, destroot string) error {
	abspath := fs.abspath(path)
	destpath := filepath.Join(destroot, path)
	return fs.client.Rename(abspath, destpath)
}

func (fs *SftpDirFs) Lstat(p string) (os.FileInfo, error) {
	abspath := fs.abspath(p)
	return fs.client.Lstat(abspath)
}

var _ DirFs = (*SftpDirFs)(nil)

type SftpDirFsCreator struct {
	fs         *SftpDirFs
	createTime time.Time
	config     *DirFsConfig
}

func (c *SftpDirFsCreator) create() (DirFs, error) {
	now := time.Now()
	if c.fs.client == nil || c.createTime.Add(10*time.Minute).Before(now) {
		if c.fs.client != nil {
			log.Println("reconnect sftp client")
			c.fs.client.Close()
			c.fs.client = nil
		}
		var err error
		sftpClient, err := createSftpClient(c.config)
		if err != nil {
			return nil, err
		}
		c.createTime = time.Now()
		c.fs.client = sftpClient
	}
	return c.fs, nil
}

func (c *SftpDirFsCreator) close() {
	if c.fs != nil {
		c.fs.client.Close()
	}
}

var _ DirFsCreator = (*SftpDirFsCreator)(nil)

func createSftpDirFsCreator(config DirFsConfig) (DirFsCreator, error) {
	fs := SftpDirFs{
		DirFsBase: DirFsBase{
			Path: config.Path,
		},
	}
	return &SftpDirFsCreator{
		fs:         &fs,
		createTime: time.Now(),
		config:     &config,
	}, nil
}

func init() {
	registCreatorFactory("sftp", createSftpDirFsCreator)
}
