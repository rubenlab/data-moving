package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

func getKnownHostsFile(providedPath string) string {
	if providedPath != "" {
		return providedPath
	}
	dirname, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(dirname, ".ssh/known_hosts")
	} else {
		return ""
	}
}

var hostKeyCallback *ssh.HostKeyCallback
var signer *ssh.Signer

func createSftpClient(config *DirFsConfig) (*sftp.Client, error) {
	var err error
	if hostKeyCallback == nil {
		hostKeyCallbackImpl, err := kh.New(getKnownHostsFile(config.KnownHosts))
		if err != nil {
			return nil, err
		}
		hostKeyCallback = &hostKeyCallbackImpl
	}

	privateKey, err := ioutil.ReadFile(config.IdentityFile)
	if err != nil {
		log.Fatalf("can't load identity file: %v", err)
	}

	if signer == nil {
		var signerImpl ssh.Signer
		secret := config.Password
		// Create the Signer for this private key.
		if secret == "" {
			signerImpl, err = ssh.ParsePrivateKey(privateKey)
		} else {
			signerImpl, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(secret))
		}
		if err != nil {
			return nil, err
		}
		signer = &signerImpl
	}

	sshClient := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(*signer),
		},
		HostKeyCallback: *hostKeyCallback,
		Timeout:         10 * time.Second,
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", config.Host+":"+fmt.Sprint(config.Port), sshClient)
	if err != nil {
		return nil, err
	}
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	return client, nil
}
