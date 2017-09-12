package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"path"
	"time"
)

var (
	keymap = map[string]string{"0": "11", "1": "2", "2": "3", "3": "4", "4": "5", "5": "6", "6": "7", "7": "8", "8": "9", "9": "10",
		"A": "30", "B": "48", "C": "46", "D": "32", "E": "18", "F": "33", "G": "34", "H": "35", "I": "23", "J": "36",
		"K": "37", "L": "38", "M": "50", "N": "49", "O": "24", "P": "25", "Q": "16", "R": "19", "S": "31", "T": "20",
		"U": "22", "V": "47", "W": "17", "X": "45", "Y": "21", "Z": "44",
		"up": "104", "down": "109", "esc": "1", "enter": "28", "shift": "42",
	}
)

// ClientSSH clientssh struct
type ClientSSH struct {
	client *ssh.Client
}

var (
	client *ClientSSH
)

// NewClientSSH create a ssh client
func NewClientSSH() *ClientSSH {
	if client == nil {
		client = &ClientSSH{}
	}
	return client
}

func connect(user, password, host string, port int) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	return client, err

}

//Open open
func (clt *ClientSSH) Open(usr, password string) error {

	s, err := connect(usr, password, "192.168.1.6", 22)
	clt.client = s
	return err
}

//Run run command
func (clt *ClientSSH) Run(cmd string) (string, error) {

	var b bytes.Buffer

	session, err := clt.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	session.Stdout = &b
	err = session.Run(cmd)

	return b.String(), err

}

//KeyPress simulate keyboard press
func (clt *ClientSSH) KeyPress(key string) {

	if code, ok := keymap[key]; ok {
		clt.Run("simulate_key /dev/input/event1 " + code)
	}

}

//FileTransfer transfer file
func (clt *ClientSSH) FileTransfer(remoteFile string) {

	sftpClient, err := sftp.NewClient(clt.client)
	if err != nil {
		return
	}
	defer sftpClient.Close()

	srcFile, err1 := sftpClient.Open(remoteFile)
	if err1 != nil {
		return
	}

	defer srcFile.Close()

	var localFileName = path.Base(remoteFile)
	dstFile, err := os.Create(path.Join("/tmp", localFileName))
	if err != nil {
		return
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		return
	}

	fmt.Println("copy file from remote server finished!")

}
