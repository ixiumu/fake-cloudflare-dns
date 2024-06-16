package main

import (
	"log"
	"net"
	"sort"
	"time"
	"math/rand"
)

// Test and sort IP addresses based on connection speed
func testAndSortIPs(ips []string) []SpeedTestResult {
	results := []SpeedTestResult{}

	for _, ip := range ips {
		speed, err := testConnectionSpeed(ip)
		if err != nil {
			log.Printf("Failed to connect to IP address %s: %s\n", ip, err)
			continue
		}

		results = append(results, SpeedTestResult{
			IP:    ip,
			Speed: speed,
		})
	}

	// Sort by speed in ascending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Speed < results[j].Speed
	})

	return results
}

// Randomly select one from the first five elements
func getIP(ips []SpeedTestResult) string {
	if len(ips) == 0 {
		return "1.0.0.1"
	}
	
	var firstFive []SpeedTestResult
	if len(ips) >= 5 {
		firstFive = ips[:5]
	} else {
		firstFive = ips
	}

	// Randomly select one from the first five elements
	randomIndex := rand.Intn(len(firstFive))
	randomElement := firstFive[randomIndex]

	return randomElement.IP
}

func testConnectionSpeed(ip string) (time.Duration, error) {
	pingTimes := 2
	var totalDuration time.Duration

	for i := 0; i < pingTimes; i++ {
		start := time.Now()
		var host string

		if isIPv6(ip) {
			host = "[" + ip + "]:80"
		} else {
			host = ip + ":80"
		}

		conn, err := net.DialTimeout("tcp", host, time.Millisecond*1000)
		if err != nil {
			return 0, err
		}

		conn.Close()

		duration := time.Since(start)
		totalDuration += duration
	}

	averageDuration := totalDuration / time.Duration(pingTimes)
	return averageDuration, nil
}
