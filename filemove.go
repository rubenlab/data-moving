package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FsType string

const (
	Local FsType = "local"
	SMB          = "smb"
)

const DirFileMode = os.FileMode(int(0770))
const FileFileMode = os.FileMode(int(0660))

type WalkFunc func(path string, info fs.FileInfo, level int, err error) error

type WalkDirFunc func(path string, d fs.DirEntry, level int, err error) error

type DirFs interface {
	Walk(enterDir WalkDirFunc, enterFile WalkFunc, exitDir WalkDirFunc)
	MkdirAll(path string) error
	MkdirAllAbs(rootpath string, relpath string) error
	Open(path string) (io.ReadWriteCloser, error)
	Create(path string) (io.ReadWriteCloser, error)
	OpenFile(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error)
	Chmod(path string, mode os.FileMode) error
	Chown(path string, uid, gid int) error
	// remove file or (empty) dir
	Remove(path string) error
	Move(path string, destroot string) error
	Lstat(p string) (os.FileInfo, error)
}

type DirFsCreator interface {
	create() (DirFs, error)
	close()
}

var fsCreatorMap = make(map[string]func(config DirFsConfig) (DirFsCreator, error))

func registCreatorFactory(typev string, create func(config DirFsConfig) (DirFsCreator, error)) {
	fsCreatorMap[typev] = create
}

func CreateFsCreator(config DirFsConfig) (DirFsCreator, error) {
	creatorFactory, ok := fsCreatorMap[string(config.Type)]
	if !ok {
		return nil, nil
	}
	creator, err := creatorFactory(config)
	if err != nil {
		return nil, err
	}
	return creator, nil
}

func KeepFileMove(config *AppConfig) {
	sourceFsCreator, err := CreateFsCreator(config.Source)
	if err != nil {
		log.Printf("can't create source fs creator, the error is %v\n", err)
		return
	}
	destFsCreator, err := CreateFsCreator(config.Dest)
	if err != nil {
		log.Printf("can't create dest fs creator, the error is %v\n", err)
		return
	}
	for {
		log.Printf("execute file move\n")
		oneFileMove(sourceFsCreator, destFsCreator, config.Dustbin, config.Execution)
		log.Printf("move finished\n")
		time.Sleep(5 * time.Second)
	}
}

func oneFileMove(sourceFsCreator DirFsCreator, destFsCreator DirFsCreator, dustbin string, config ExecutionConfig) {
	sourceFs, err := sourceFsCreator.create()
	if err != nil {
		log.Printf("can't create source fs, the error is %v\n", err)
		return
	}
	destFs, err := destFsCreator.create()
	if err != nil {
		log.Printf("can't create dest fs, the error is %v\n", err)
		return
	}
	FileMove(sourceFs, destFs, dustbin, config)
	log.Printf("finished\n")
}

func FileMove(source DirFs, dest DirFs, dustbin string, config ExecutionConfig) {
	source.Walk(func(path string, d fs.DirEntry, level int, err error) error {
		// make dir for destination and dustbin
		if err != nil {
			return err
		}
		err = dest.MkdirAll(path)
		if err != nil {
			log.Printf("can't create parent folders on destination for folder %s,\nthe error is: %v\n", path, err)
			return err
		}
		dest.Chown(path, config.Uid, config.Gid)
		err = source.MkdirAllAbs(dustbin, path)
		if err != nil {
			log.Printf("can't create parent folders on dustbin for folder %s,\nthe error is: %v\n", path, err)
			return err
		}
		return nil
	}, func(path string, info fs.FileInfo, level int, err error) error {
		// copy file to the destination, and move source file to the dustbin
		if level < config.StartLevel {
			return nil
		}
		targetPath := path

		// rename target file if needed
		if !config.Overwrite {
			targetPath, err = avoidExistsFile2(dest, path)
			if err != nil {
				log.Printf("can't avoid exists file, the error is:\n%v", err)
				return err
			}
		}

		// copy file
		targetFile, err := dest.Create(targetPath)
		if err != nil {
			log.Printf("can't open target file, the error is:\n%v", err)
			return err
		}
		defer targetFile.Close()
		sourceFile, err := source.OpenFile(path, os.O_RDWR, 0)
		if err != nil {
			log.Printf("can't open source file, the error is:\n%v", err)
			return err
		}
		defer sourceFile.Close()
		_, err = io.Copy(targetFile, sourceFile)
		if err != nil {
			log.Printf("can't copy file to the remote server, the error is:\n%v", err)
			return err
		}
		sourceFile.Close()

		// chmod
		err = dest.Chmod(targetPath, FileFileMode)
		if err != nil {
			log.Printf("failed to change file mode, the error is:\n%v", err)
		}

		// chown
		if config.Gid != 0 {
			err = dest.Chown(targetPath, config.Uid, config.Gid)
			if err != nil {
				log.Printf("failed to change file owner, the error is:\n%v", err)
			}
		}

		// move to dustbin
		err = source.Move(path, dustbin)
		if err != nil {
			log.Printf("failed to move file to the dustbin, the error is:\n%v", err)
			filename := filepath.Base(path)
			if filename == ".DS_Store" {
				source.Remove(path)
			}
		}

		return nil
	}, func(path string, d fs.DirEntry, level int, err error) error {
		// clear empty folders
		if level >= config.StartLevel {
			source.Remove(path)
		}
		return nil
	})
}

func avoidExistsFile2(dest DirFs, path string) (string, error) {
	fileName := filepath.Base(path)
	targetFolder := filepath.Dir(path)
	_, notExistErr := dest.Lstat(path)
	for i := 1; !os.IsNotExist(notExistErr); i++ {
		if notExistErr != nil {
			return "", notExistErr
		}
		newFileName := createNewFilename(fileName, i)
		path = filepath.Join(targetFolder, newFileName)
		_, notExistErr = dest.Lstat(path)
	}
	return path, nil
}

func createNewFilename(filename string, num int) string {
	ext := filepath.Ext(filename)
	fileBaseName := strings.TrimSuffix(filename, ext)
	return fmt.Sprintf("%s(%d)%s", fileBaseName, num, ext)
}
