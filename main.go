package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/txn2/txeh"
)

var blockList = []string{
	"facebook.com", "www.facebook.com",
	"x.com", "www.x.com",
	"github.com", "www.github.com",
	"instagram.com", "www.instagram.com",
	"linkedin.com", "www.linkedin.com",
}

const (
	redirectIP = "0.0.0.0"
	name       = "study"
)

func printHelp() {
	fmt.Printf("Usage: sudo %s <command>\n", name)
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  on      Block distractions")
	fmt.Println("  off     Restore access")
	fmt.Println("  status  Check if study mode is active")
}

func main() {

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	// loads the current /etc/hosts file
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		log.Fatalf("Error loading /etc/hosts: %v", err)
	}

	switch command {
	case "on":
		enable(hosts)
	case "off":
		disable(hosts)
	case "status":
		showStatus(hosts)
	default:
		fmt.Printf("Unknown command %s\n", command)
		printHelp()
	}
}

func enable(hosts *txeh.Hosts) {
	for _, site := range blockList {
		hosts.AddHost(redirectIP, site)
	}

	saveAndFlush(hosts)
	fmt.Println("Just study")
}

func disable(hosts *txeh.Hosts) {
	for _, site := range blockList {
		hosts.RemoveHost(site)
	}

	saveAndFlush(hosts)
	fmt.Println("Please remember if you use this while studying you'll become a failure and disgrace to yourself, family and humanity.")
}

func saveAndFlush(hosts *txeh.Hosts) {
	err := hosts.Save()
	if err != nil {
		log.Fatalf("Failed to save hosts file (are you root?): %v", err)
	}

	// OS caches the real IP either way even after redirection, so flush the dns cache forcing it to always read /etc/hosts
	cmd := exec.Command("resolvectl", "flush-caches")
	err = cmd.Run()
	if err != nil {
		_ = exec.Command("systemd-resolve", "--flush-caches").Run()
	}
}

func showStatus(hosts *txeh.Hosts) {
	bait := blockList[0]
	found, ip, _ := hosts.HostAddressLookup(bait, txeh.IPFamilyV4)

	if found && ip == redirectIP {
		fmt.Println("You should be studying rn")
		fmt.Printf("Distractions are currently redirected to %s.\n", redirectIP)
	} else {
		fmt.Println("You're not studying")
	}
}
