package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	host := os.Args[2:2]
	port := "22"

	sshConfig := &ssh.ClientConfig{
		User:            "simple",
		Auth:            []ssh.AuthMethod{ssh.Password("simple")},
		HostKeyCallback: func(host string, remote net.Addr, key ssh.PublicKey) error { return nil },
	}

	connection, err := ssh.Dial("tcp", host+":"+port, sshConfig)

	if err != nil {
		fmt.Printf("Failed to dial: %s", err)
		os.Exit(1)
		//return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()

	if err != nil {
		fmt.Printf("Failed to dial: %s", err)
		os.Exit(1)
		//return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		fmt.Printf("request for pseudo terminal failed: %s", err)
		os.Exit(1)
		//return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Printf("Unable to setup stdin for session: %v", err)
		os.Exit(1)
		//return fmt.Errorf("Unable to setup stdin for session: %v", err)
	}

	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		fmt.Printf("Unable to setup stdout for session: %v", err)
		os.Exit(1)
		//return fmt.Errorf("Unable to setup stdout for session: %v", err)
	}

	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		fmt.Printf("Unable to setup stderr for session: %v", err)
		os.Exit(1)
		//return fmt.Errorf("Unable to setup stderr for session: %v", err)
	}

	go io.Copy(os.Stderr, stderr)

	if err := session.Shell(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	session.Wait()

}
