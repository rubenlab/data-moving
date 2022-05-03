package main

import (
	"fmt"
	"io"
	"io/fs"
	iofs "io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hirochachacha/go-smb2"
)

type SmbDirFs struct {
	DirFsBase
	conn    *net.Conn
	session *smb2.Session
	share   *smb2.Share
}

func (fs *SmbDirFs) Walk(enterDir WalkDirFunc, enterFile WalkFunc, exitDir WalkDirFunc) {
	fs.walkDir(".", 1, enterDir, enterFile, exitDir)
}

func (fs *SmbDirFs) walkDir(dir string, level int, enterDir WalkDirFunc, enterFile WalkFunc, exitDir WalkDirFunc) {
	dirpath := fs.abspath(dir)
	files, err := fs.share.ReadDir(dirpath)
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

func (fs *SmbDirFs) MkdirAll(path string) error {
	abspath := fs.abspath(path)
	err := fs.share.MkdirAll(abspath, DirFileMode)
	if err != nil {
		return err
	}
	fs.share.Chmod(abspath, DirFileMode)
	return nil
}

func (fs *SmbDirFs) MkdirAllAbs(rootpath string, relpath string) error {
	abspath := filepath.Join(rootpath, relpath)
	return fs.share.MkdirAll(abspath, DirFileMode)
}

func (fs *SmbDirFs) Open(path string) (io.ReadWriteCloser, error) {
	abspath := fs.abspath(path)
	return fs.share.Open(abspath)
}

func (fs *SmbDirFs) Create(path string) (io.ReadWriteCloser, error) {
	abspath := fs.abspath(path)
	return fs.share.Create(abspath)
}

func (fs *SmbDirFs) OpenFile(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
	abspath := fs.abspath(name)
	return fs.share.OpenFile(abspath, flag, perm)
}

func (fs *SmbDirFs) Chmod(path string, mode os.FileMode) error {
	abspath := fs.abspath(path)
	return fs.share.Chmod(abspath, mode)
}

func (fs *SmbDirFs) Chown(path string, uid, gid int) error {
	// can't chown on smb fs
	return nil
}

func (fs *SmbDirFs) Remove(path string) error {
	abspath := fs.abspath(path)
	return fs.share.Remove(abspath)
}

func (fs *SmbDirFs) Move(path string, destroot string) error {
	abspath := fs.abspath(path)
	destpath := filepath.Join(destroot, path)
	return fs.share.Rename(abspath, destpath)
}

func (fs *SmbDirFs) Lstat(p string) (os.FileInfo, error) {
	abspath := fs.abspath(p)
	return fs.share.Lstat(abspath)
}

var _ DirFs = (*SmbDirFs)(nil)

type SmbDirFsCreator struct {
	fs         *SmbDirFs
	createTime time.Time
	config     *DirFsConfig
}

func (c *SmbDirFsCreator) create() (DirFs, error) {
	now := time.Now()
	if c.fs.session == nil || c.createTime.Add(10*time.Minute).Before(now) {
		if c.fs.session != nil {
			log.Println("reconnect smb client")
			c.close()
		}
		var err error
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.config.Host, c.config.Port))
		if err != nil {
			panic(err)
		}

		d := &smb2.Dialer{
			Initiator: &smb2.NTLMInitiator{
				User:     c.config.Username,
				Password: c.config.Password,
				Domain:   c.config.Domain,
			},
		}

		session, err := d.Dial(conn)
		if err != nil {
			panic(err)
		}

		share, err := session.Mount(c.config.ShareName)
		if err != nil {
			panic(err)
		}

		c.createTime = time.Now()
		c.fs.conn = &conn
		c.fs.session = session
		c.fs.share = share
	}
	return c.fs, nil
}

func (c *SmbDirFsCreator) close() {
	if c.fs != nil {
		if c.fs.share != nil {
			c.fs.session.Logoff()
			(*c.fs.conn).Close()
			c.fs.share = nil
			c.fs.session = nil
			c.fs.conn = nil
		}
	}
}

var _ DirFsCreator = (*SftpDirFsCreator)(nil)

func createSmbDirFsCreator(config DirFsConfig) (DirFsCreator, error) {
	fs := SmbDirFs{
		DirFsBase: DirFsBase{
			Path: config.Path,
		},
	}
	return &SmbDirFsCreator{
		fs:         &fs,
		createTime: time.Now(),
		config:     &config,
	}, nil
}

func init() {
	registCreatorFactory("smb", createSmbDirFsCreator)
}
