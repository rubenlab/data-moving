package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("./config-test.yml")
	if err != nil {
		t.Error(err)
	}
	if config.Source.StartLevel != 2 {
		t.Errorf("start level is not correct, the loaded config level is: %d", config.Source.StartLevel)
		t.Error(config)
	}
	if config.Dest.Path != "/usr/yi1" {
		t.Errorf("path is not correct, the loaded path is: %v", config.Dest.Path)
		t.Error(config)
	}
}

func TestGetKnownHost(t *testing.T) {
	hostFile := getKnownHostsFile("/home/ytm/.ssh/known_hosts")
	if hostFile != "/home/ytm/.ssh/known_hosts" {
		t.Errorf("unexpected host file location: %v", hostFile)
	}
}

func TestCreateNewFileName(t *testing.T) {
	fileName := "test.txt"
	newFileName := createNewFilename(fileName, 2)
	if newFileName != "test(2).txt" {
		t.Errorf("new file name is not correct: %v", newFileName)
	}
}

func TestRemoveEmptyFolder(t *testing.T) {
	path := "./emptyTestFolder"
	os.Mkdir(path, DirFileMode)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("create empty folder failed")
		return
	}
	removeEmptyFolder(path)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("empty folder not deleted")
	}
}

func TestEncryption(t *testing.T) {
	data := "random text, hahaha"
	secret := "lalala_this@is#password"
	encryptedStr, err := encryptToString(secret, []byte(data))
	if err != nil {
		t.Error(err)
	}
	dataBytes, err := decryptString(secret, encryptedStr)
	if err != nil {
		t.Error(err)
	}
	if data != string(dataBytes) {
		t.Errorf("encrypt and decrypt failed, the result is: %s", string(dataBytes))
	}
}
