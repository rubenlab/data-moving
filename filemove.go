package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
)

const DirFileMode = os.FileMode(int(0770))
const FileFileMode = os.FileMode(int(0660))

func ExecuteMove(config *Config, client *sftp.Client) {
	rootDir := config.Source.RootDir
	targetDir := config.Dest.Path
	dustbin := config.Dustbin
	startLevel := config.Source.StartLevel
	overwrite := config.Source.Overwrite
	executeMoveInternal(client, overwrite, rootDir, startLevel, targetDir, dustbin, rootDir, 1)
}

func executeMoveInternal(client *sftp.Client, overwrite bool, rootDir string, startLevel int, targetDir string, dustbin string, currentDir string, currentLevel int) {
	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		log.Printf("can't open folder %s,\nthe error is: %v\n", currentDir, err)
		return
	}
	for _, file := range files {
		filePath := filepath.Join(currentDir, file.Name())
		if file.IsDir() {
			err = mkdirAll(client, rootDir, targetDir, dustbin, filePath)
			if err != nil {
				log.Printf("can't create parent folders for folder %s,\nthe error is: %v\n", currentDir, err)
				continue
			}
			executeMoveInternal(client, overwrite, rootDir, startLevel, targetDir, dustbin, filePath, currentLevel+1)
			if currentLevel >= startLevel {
				removeEmptyFolder(filePath)
			}
		} else {
			if currentLevel < startLevel {
				continue
			}
			moveFile(client, overwrite, rootDir, targetDir, dustbin, filePath)
		}
	}
}

func mkdirAll(client *sftp.Client, rootDir string, targetDir string, dustbin string, currentDir string) error {
	targetFolder, err := replacePath(currentDir, rootDir, targetDir)
	if err != nil {
		return errors.Wrap(err, "can't replace path for remote dir")
	}
	err = client.MkdirAll(targetFolder)
	if err != nil {
		return errors.Wrap(err, "can't mkdir for remote folder")
	}
	err = client.Chmod(targetFolder, DirFileMode)
	if err != nil {
		log.Printf("failed to change file mode of folder, the error is:\n%v", err)
	}
	dustbinDir, err := replacePath(currentDir, rootDir, dustbin)
	if err != nil {
		return errors.Wrap(err, "can't replace path for dustbin")
	}
	err = os.MkdirAll(dustbinDir, DirFileMode)
	if err != nil {
		return errors.Wrap(err, "can't mkdir for dustbin folder")
	}
	return nil
}

func moveFile(client *sftp.Client, overwrite bool, rootDir string, targetDir string, dustbin string, filePath string) {
	targetPath, err := replacePath(filePath, rootDir, targetDir)
	if err != nil {
		log.Printf("error in generating targetPath, the error is:\n%v", err)
		return
	}
	targetPathDir := filepath.Dir(targetPath)
	dustbinPath, err := replacePath(filePath, rootDir, dustbin)
	if err != nil {
		log.Printf("error in generating dustbinPath, the error is:\n%v", err)
		return
	}

	if !overwrite {
		targetPath, err = avoidExistsFile(client, targetPath, targetPathDir)
		if err != nil {
			log.Printf("can't avoid exists file, the error is:\n%v", err)
			return
		}
	}

	targetFile, err := client.Create(targetPath)
	if err != nil {
		log.Printf("can't open target file, the error is:\n%v", err)
		return
	}
	defer targetFile.Close()
	sourceFile, err := os.Open(filePath)
	if err != nil {
		log.Printf("can't open source file, the error is:\n%v", err)
		return
	}
	defer sourceFile.Close()
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		log.Printf("can't copy file to the remote server, the error is:\n%v", err)
		return
	}
	sourceFile.Close()
	err = client.Chmod(targetPath, FileFileMode)
	if err != nil {
		log.Printf("failed to change file mode, the error is:\n%v", err)
	}
	err = os.Rename(filePath, dustbinPath)
	if err != nil {
		log.Printf("failed to move file to the dustbin, the error is:\n%v", err)
	}
}

func avoidExistsFile(client *sftp.Client, filePath string, targetFolder string) (string, error) {
	fileName := filepath.Base(filePath)
	_, notExistErr := client.Lstat(filePath)
	for i := 1; !os.IsNotExist(notExistErr); i++ {
		if notExistErr != nil {
			return "", notExistErr
		}
		newFileName := createNewFilename(fileName, i)
		filePath = filepath.Join(targetFolder, newFileName)
		_, notExistErr = client.Lstat(filePath)
	}
	return filePath, nil
}

func createNewFilename(filename string, num int) string {
	ext := filepath.Ext(filename)
	fileBaseName := strings.TrimSuffix(filename, ext)
	return fmt.Sprintf("%s(%d)%s", fileBaseName, num, ext)
}

func removeEmptyFolder(folderPath string) {
	f, err := os.Open(folderPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		err = os.Remove(folderPath)
		if err != nil {
			log.Println(err)
		}
	}
}

func replacePath(path string, sourceDir string, destDir string) (string, error) {
	relpath, err := filepath.Rel(sourceDir, path)
	if err != nil {
		return "", err
	}
	newpath := filepath.Join(destDir, relpath)
	return newpath, nil
}
