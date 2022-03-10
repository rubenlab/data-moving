package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sevlyar/go-daemon"
)

var (
	secretFile = flag.String("pwdfile", "", "file location that store the private key password")
	asDaemon   = flag.Bool("d", false, "run in daemon")
)

func main() {
	flag.Parse()

	if *asDaemon {
		cntxt := &daemon.Context{
			PidFileName: "tohpc.pid",
			PidFilePerm: 0644,
			LogFileName: "tohpc.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
		}

		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()

		log.Print("- - - - - - - - - - - - - - -")
		log.Print("daemon started")
	}

	startMoveFile()
}

func startMoveFile() {
	config, err := LoadConfig("./config.yml")
	if err != nil {
		log.Fatalf("can't load config with error: %v", err)
		return
	}

	var secret []byte

	if *secretFile != "" {
		secret, err = ioutil.ReadFile(*secretFile)
		if err != nil {
			log.Fatalf("can't read pwdfile: %v", err)
			return
		}
		os.Remove(*secretFile)
	}

	secretStr := string(secret)
	secretStr = strings.TrimSpace(secretStr)

	privateKey, err := ioutil.ReadFile(config.Dest.IdentityFile)
	if err != nil {
		log.Fatalf("can't load identity file: %v", err)
	}

	for {
		log.Printf("execute file move\n")
		executeMoveOnce(config, privateKey, secretStr)
		log.Printf("finished\n")
		time.Sleep(5 * time.Second)
	}
}

func executeMoveOnce(config *Config, privateKey []byte, secret string) {
	client, err := createSftpClient(config, privateKey, []byte(secret))
	if err != nil {
		log.Printf("can't create sftp client: %v\n", err)
		return
	}
	defer client.Close()
	ExecuteMove(config, client)
}
