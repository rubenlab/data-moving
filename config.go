package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

type Config struct {
	Source struct {
		RootDir    string `yaml:"root-dir"`
		StartLevel int    `yaml:"start-level"` // Files parallel to the start level are ignored, all files should be placed in the start level directory or deeper.
		Overwrite  bool   // Overwrite existing file on the remote server
	}
	Dest struct {
		Host         string
		Port         int
		IdentityFile string `yaml:"identity-file"`
		Username     string
		Path         string
		Gid          int // if not zero, will be used to set file group
		Uid          int // if gid is not zero, a correct uid value should be set
	}
	Dustbin    string
	KnownHosts string `yaml:"known-hosts"`
}

func LoadConfig(path string) (*Config, error) {
	config := Config{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "can not open config file")
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Wrap(err, "can not unmarshal config data")
	}
	return &config, nil
}
