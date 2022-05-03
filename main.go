package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/sevlyar/go-daemon"
	"golang.org/x/term"
)

var (
	secretFile = flag.String("pwdfile", "", "file location that store the decryption secret")
	asDaemon   = flag.Bool("d", false, "run in daemon")
	encrypt    = flag.String("encrypt", "", "encrypt password")
	decrypt    = flag.String("decrypt", "", "decrypt password")
)

func main() {
	flag.Parse()

	if *encrypt != "" {
		encryptFunc()
		return
	}

	if *decrypt != "" {
		decryptFunc()
		return
	}

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

func encryptFunc() {
	secret, err := inputSecret()
	if err != nil {
		fmt.Printf("error in input password, error: %v", err)
		return
	}
	if secret == "" {
		fmt.Println("please input a none empty secret")
		return
	}

	encrypted, err := encryptToString(secret, []byte(*encrypt))
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println(encrypted)
}

func decryptFunc() {
	secret, err := inputSecret()
	if err != nil {
		fmt.Printf("error in input password, error: %v", err)
		return
	}
	if secret == "" {
		fmt.Println("please input a none empty secret")
		return
	}

	decrypted, err := decryptString(secret, *decrypt)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println(string(decrypted))
}

func inputSecret() (string, error) {
	fmt.Println("Input your secret:")
	secretArr, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	secret := string(secretArr)
	return secret, nil
}

func startMoveFile() {
	var err error
	var secretStr string
	if *secretFile != "" {
		var secret []byte
		secret, err = ioutil.ReadFile(*secretFile)
		if err != nil {
			log.Fatalf("can't read pwdfile: %v", err)
			return
		}
		os.Remove(*secretFile)
		secretStr = string(secret)
	} else if !*asDaemon {
		fmt.Println("Input password for private key:")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("error in input password, error: %v", err)
			return
		}
		secretStr = string(bytePassword)
	}
	secretStr = strings.TrimSpace(secretStr)
	config, err := LoadAppConfig("./config.yml", secretStr)
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}
	KeepFileMove(config)
}
