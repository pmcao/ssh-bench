package main

import (
	"fmt"
	"github.com/paulbellamy/ratecounter"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

//getEnvWithDefault returns the environment value for key
//returning fallback instead if it is missing or blank
func getEnvWithDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <user> <host:port> <command>", os.Args[0])
	}
	// We're recording marks-per-1second
	counter := ratecounter.NewRateCounter(1 * time.Second)

	for i := 0; i < 99999999999; i++ {
		// Record an event happening
		counter.Incr(1)
		connectToHost(os.Args[1], os.Args[2])
		fmt.Print("%d auth tries \r", counter.Rate())
	}
}

func connectToHost(user, host string) {
	pass := "ssh-bench"
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return
	}

	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer client.Close()
	defer session.Close()
}
