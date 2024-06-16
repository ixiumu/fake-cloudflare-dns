package main

import (
	"log"
	"net"
	"os"
	"strings"
)

func LookupIP(ipFile string) []string {
	// Open the file
	file, err := os.Open(ipFile)
	if err != nil {
		log.Printf("Unable to open the file: %s", err)
		return []string{}
	}
	defer file.Close()

	// Create a string slice to store the resolved IP addresses
	ipAddresses := []string{}

	// Read the file line by line
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if n == 0 || err != nil {
			break
		}

		// Parse domains and IP addresses
		lines := string(buffer[:n])
		entries := parseEntries(lines)
		for _, entry := range entries {
			if isDomain(entry) {
				ips, err := net.LookupIP(entry)
				if err != nil {
					log.Printf("Unable to resolve domain %s: %s\n", entry, err)
				} else {
					for _, ip := range ips {
						ipAddresses = append(ipAddresses, ip.String())
					}
				}
			} else if isIPv4(entry) || isIPv6(entry) {
				ipAddresses = append(ipAddresses, entry)
			} else {
				log.Printf("Unrecognized entry: %s\n", entry)
			}
		}
	}

	// Print the resolved IP addresses
	// fmt.Println("Resolution results:")
	// for _, ip := range ipAddresses {
	// 	log.Println(ip)
	// }
	return ipAddresses
}

// Parse domains and IP addresses, extract non-empty entries
func parseEntries(content string) []string {
	lines := []string{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

// Check if an entry is a domain name
func isDomain(entry string) bool {
	return !isIPv4(entry) && !isIPv6(entry)
}

// Check if an entry is an IPv4 address
func isIPv4(entry string) bool {
	return net.ParseIP(entry) != nil && strings.Contains(entry, ".")
}

// Check if an entry is an IPv6 address
func isIPv6(entry string) bool {
	return net.ParseIP(entry) != nil && strings.Contains(entry, ":")
}