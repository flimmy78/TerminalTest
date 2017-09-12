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

	var code string
	switch key {
	case "0":
		code = "11"
		break
	case "1":
		code = "2"
		break
	case "2":
		code = "3"
		break
	case "3":
		code = "4"
		break
	case "4":
		code = "5"
		break
	case "5":
		code = "6"
		break
	case "6":
		code = "7"
		break
	case "7":
		code = "8"
		break
	case "8":
		code = "9"
		break
	case "9":
		code = "10"
		break
	case "up":
		code = "104"
		break
	case "down":
		code = "109"
		break
	case "enter":
		code = "28"
		break
	case "esc":
		code = "1"
		break
	}

	clt.Run("simulate_key /dev/input/event1 " + code)
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
