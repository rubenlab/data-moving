package main

import (
	"fmt"
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

func createSftpClient(config *Config, privateKey []byte, secret []byte) (*sftp.Client, error) {
	var err error
	if hostKeyCallback == nil {
		hostKeyCallbackImpl, err := kh.New(getKnownHostsFile(config.KnownHosts))
		if err != nil {
			return nil, err
		}
		hostKeyCallback = &hostKeyCallbackImpl
	}

	if signer == nil {
		var signerImpl ssh.Signer
		// Create the Signer for this private key.
		if len(secret) == 0 {
			signerImpl, err = ssh.ParsePrivateKey(privateKey)
		} else {
			signerImpl, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, secret)
		}
		if err != nil {
			return nil, err
		}
		signer = &signerImpl
	}

	sshClient := &ssh.ClientConfig{
		User: config.Dest.Username,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(*signer),
		},
		HostKeyCallback: *hostKeyCallback,
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", config.Dest.Host+":"+fmt.Sprint(config.Dest.Port), sshClient)
	if err != nil {
		return nil, err
	}
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	return client, nil
}

var sftpClient *sftp.Client
var connectTime time.Time

func getSftpClient(config *Config, privateKey []byte, secret []byte) (*sftp.Client, error) {
	if sftpClient == nil {
		var err error
		sftpClient, err = createSftpClient(config, privateKey, secret)
		if err != nil {
			return nil, err
		}
		connectTime = time.Now()
		return sftpClient, nil
	} else {
		now := time.Now()
		if connectTime.Add(10 * time.Minute).Before(now) {
			log.Println("reconnect sftp client")
			sftpClient.Close()
			sftpClient = nil
			var err error
			sftpClient, err = createSftpClient(config, privateKey, secret)
			if err != nil {
				return nil, err
			}
			connectTime = time.Now()
		}
		return sftpClient, nil
	}
}

func closeSftpClient() {
	if sftpClient == nil {
		return
	}
	sftpClient.Close()
}
