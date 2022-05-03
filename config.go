package main

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type DirFsConfig struct {
	Type         FsType
	Host         string
	Port         int
	Path         string // root dir path
	ShareName    string `yaml:"share-name"` // share name of smb server
	IdentityFile string `yaml:"identity-file"`
	Username     string
	Password     string
	Domain       string
	KnownHosts   string `yaml:"known-hosts"` // hosts file location
}

type ExecutionConfig struct {
	StartLevel int  `yaml:"start-level"` // Files parallel to the start level are ignored, all files should be placed in the start level directory or deeper.
	Overwrite  bool // Overwrite existing file on the remote server
	Gid        int  // if not zero, will be used to set file group on destination
	Uid        int  // if gid is not zero, a correct uid value should be set on destination
}

type AppConfig struct {
	Source     DirFsConfig
	Dest       DirFsConfig
	Dustbin    string
	Execution  ExecutionConfig
	KnownHosts string `yaml:"known-hosts"` // hosts file location
}

func LoadAppConfig(path string, secret string) (*AppConfig, error) {
	config := AppConfig{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "can not open config file")
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Wrap(err, "can not unmarshal config data")
	}
	if secret != "" {
		err = decryptConfig(&config.Source, secret)
		if err != nil {
			return nil, errors.Wrap(err, "cannot decrypt source password")
		}
		err = decryptConfig(&config.Dest, secret)
		if err != nil {
			return nil, errors.Wrap(err, "cannot decrypt dest password")
		}
	}
	if config.KnownHosts != "" {
		config.Source.KnownHosts = config.KnownHosts
		config.Dest.KnownHosts = config.KnownHosts
	}
	return &config, nil
}

func decryptConfig(config *DirFsConfig, secret string) error {
	if secret == "" {
		return nil
	}
	if config.Password != "" {
		arr, err := decryptString(secret, config.Password)
		if err != nil {
			return errors.Wrap(err, "cannot decrypt dest password")
		}
		config.Password = string(arr)
	}
	return nil
}
