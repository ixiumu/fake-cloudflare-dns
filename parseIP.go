package main

import (
	"log"
	"net"
	"os"
	"strings"
)

var (
	ipList   []string
	ipListv6 []string
)

func UpdateIPList(ip string) {
	if isIPv4(ip) {
		ipList = append(ipList, ip)
		logger.Infof("IP: %s", ip)
	} else if isIPv6(ip) {
		ipListv6 = append(ipListv6, ip)
		logger.Infof("IPv6: %s", ip)
	}
}

func LookupDefaultIP(domain string) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		logger.Infof("Unable to resolve domain %s: %s\n", domain, err)
	} else {
		for _, ip := range ips {
			if fixedIPAddressV4 == "" && isIPv4(ip.String()) {
				fixedIPAddressV4 = ip.String()
				logger.Infof("DefaultIP: %s", ip.String())
			} else if fixedIPAddressV6 == "" && isIPv6(ip.String()) {
				fixedIPAddressV6 = ip.String()
				logger.Infof("DefaultIPv6: %s", ip.String())
			}
		}
	}
}

func LookupIP(ipFile string) {
	// Open the file
	file, err := os.Open(ipFile)
	if err != nil {
		logger.Errorf("Unable to open the file: %s", err)
		return
		// return []string{}
	}
	defer file.Close()

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
						UpdateIPList(ip.String())
					}
				}
			} else if isIPv4(entry) || isIPv6(entry) {
				UpdateIPList(entry)
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
	// return ipAddresses
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
	if len(entry) < 10 {
		return false
	}
	return net.ParseIP(entry) != nil && strings.Contains(entry, ":")
}