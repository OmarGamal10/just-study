package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/txn2/txeh"
)

var blockList = []string{
	"facebook.com", "www.facebook.com",
	"x.com", "www.x.com",
	"instagram.com", "www.instagram.com",
	"tiktok.com", "www.tiktok.com",
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

	ipChan := make(chan string, 100)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for ip := range ipChan {
			exec.Command("ss", "-K", "dst", ip).Run()
		}
	}()

	var lookupWg sync.WaitGroup
	for _, domain := range blockList {
		hosts.AddHost(redirectIP, domain)
		lookupWg.Add(1)
		go func(domain string) {
			defer lookupWg.Done()
			ips, _ := net.LookupIP(domain)
			for _, ip := range ips {
				ipChan <- ip.String()
				// fmt.Printf("Blocking %s (%s)\n", domain, ip.String())
			}
		}(domain)
	}

	lookupWg.Wait()
	close(ipChan)
	wg.Wait()
	saveAndFlush(hosts)
	fmt.Println("Just study")
}

func disable(hosts *txeh.Hosts) {
	for _, domain := range blockList {
		hosts.RemoveHost(domain)
	}

	saveAndFlush(hosts)
	fmt.Println("Please remember if you use this while studying you'll become a failure and disgrace to yourself, family and humanity.")
}

func saveAndFlush(hosts *txeh.Hosts) {
	err := hosts.Save()
	if err != nil {
		log.Fatalf("Failed to save hosts file (are you root?): %v", err)
	}

	// OS caches the real IP either way even after redirection, flusing the dns cache forces it to read /etc/hosts
	cmd := exec.Command("resolvectl", "flush-caches")
	err = cmd.Run()
	if err != nil {
		_ = exec.Command("systemd-resolve", "--flush-caches").Run()
	}
}

func showStatus(hosts *txeh.Hosts) {
	canary := blockList[0] // coal miners back then carried a canary to the mines to detect toxic gases(if the canary died) :D
	found, ip, _ := hosts.HostAddressLookup(canary, txeh.IPFamilyV4)

	const (
		ColorGreen = "\033[32m"
		ColorRed   = "\033[31m"
		ColorReset = "\033[0m"
	)

	if found && ip == redirectIP {
		fmt.Printf("%sYou should be studying rn\n", ColorRed)
		fmt.Printf("Distractions are currently redirected to %s%s.\n", redirectIP, ColorReset)
	} else {
		fmt.Printf("%sYou're not studying\n%s", ColorGreen, ColorReset)
	}
}
