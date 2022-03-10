package main

import (
	"fmt"
	"os"
	"path/filepath"

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

func createSftpClient(config *Config, privateKey []byte, secret []byte) (*sftp.Client, error) {
	hostKeyCallback, err := kh.New(getKnownHostsFile(config.KnownHosts))
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	// Create the Signer for this private key.
	if secret == nil {
		signer, err = ssh.ParsePrivateKey(privateKey)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, secret)
	}
	if err != nil {
		return nil, err
	}

	sshClient := &ssh.ClientConfig{
		User: config.Dest.Username,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
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
