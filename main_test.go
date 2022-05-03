package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadAppConfig("./config-test.yml", "")
	if err != nil {
		t.Error(err)
	}
	if config.Execution.StartLevel != 2 {
		t.Errorf("start level is not correct, the loaded start level is: %d", config.Execution.StartLevel)
	}
	if config.Source.Type != "local" {
		t.Errorf("Source type is not correct, the loaded Source type is: %s", config.Source.Type)
	}
	if config.Dest.ShareName != "ukln-all$" {
		t.Errorf("Dest ShareName is not correct, the loaded ShareName is: %s", config.Dest.ShareName)
	}
	if config.Dest.KnownHosts != "/home/ytm/.ssh/known_hosts" {
		t.Errorf("Dest KnownHosts is not correct, the loaded KnownHosts is: %s", config.Dest.KnownHosts)
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
