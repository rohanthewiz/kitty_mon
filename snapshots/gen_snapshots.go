package snapshots

import (
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

func SendSnapshots(stopChan <-chan bool, doneChan chan<- bool) {
	for {
		select {
		case _, ok := <-stopChan:
			if !ok { // channel is now closed and empty
				doneChan <- true
				return
			}
		case <-time.After(6 * time.Second): // TODO config for this time
			fmt.Println("Taking a snapshot")
			takeSnapShot()
		}
	}
}

func takeSnapShot() {
	const sshUser = "sshUser"
	const privKey = "/home/me/.ssh/id_ecdsa"
	const scpServerAndPort = "myServer:<int port>"
	const serverDestPath = "/home/app/xfr/"
	const serverDestFile = "test.txt"

	clientCfg, _ := auth.PrivateKey(sshUser, privKey, ssh.InsecureIgnoreHostKey())
	client := scp.NewClient(scpServerAndPort, &clientCfg)
	err := client.Connect()
	if err != nil {
		fmt.Println("Could not establish a connection to the remote server", err)
		return
	}
	defer client.Close()

	f, err := os.Open("/home/me/test.txt")
	if err != nil {
		fmt.Println("Could not open test file")
		return
	}

	err = client.CopyFile(f, serverDestPath+serverDestFile, "0644")
	if err != nil {
		fmt.Println("Error copying file")
		return
	}
	fmt.Println("File copied successfully")
}
