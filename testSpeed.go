package main

import (
	"log"
	"fmt"
	"net"
	"sort"
	"time"
	"math/rand"
	"io/ioutil"
	"net/http"
	"strings"
)

type SpeedTestResult struct {
	IP     string
	Speed  time.Duration
}

var (
	fastestIPs []SpeedTestResult
)

func testSpeed() {
	LookupIP(filename)
	updateFixedIP()
	
	if (interval < 10) {
		interval = 10
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	for {
		<-ticker.C
		updateFixedIP()	
	}
}

func updateFixedIP() {
	// Test connection speed and sort
	fastestIPList := testAndSortIPs(ipList)	// Print sorted results

	logger.Infof("Connection speed sorting results:")
	for _, result := range fastestIPList {
		logger.Infof("IP: %s, Speed: %s\n", result.IP, result.Speed)
	}

	// Test connection speed and sort
	fastestIPListv6 := testAndSortIPs(ipListv6)	// Print sorted results

	logger.Infof("Connection speed sorting results:")
	for _, result := range fastestIPListv6 {
		logger.Infof("IP: %s, Speed: %s\n", result.IP, result.Speed)
	}

	fixedIPAddressV4 = testIPs(fastestIPList)
	if fixedIPAddressV4 != "" {
		logger.Infof("BestIP: %s", fixedIPAddressV4)
	}

	fixedIPAddressV6 = testIPs(fastestIPListv6)
	if fixedIPAddressV6 != "" {
		log.Printf("BestIPv6: %s", fixedIPAddressV6)
	}
}

func checkIPStatus(ip string) bool {
	if isIPv6(ip) {
		ip = "["+ip+"]"
	}

	url := fmt.Sprintf("http://%s/cdn-cgi/trace", ip)

	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	content := string(body)
	if strings.Contains(content, "tls=off") && strings.Contains(content, "warp=") {
		logger.Infof("Trace: %s", content)
		return true
	}

	return false
}

// func getIP(list []SpeedTestResult) string {
// 	if len(list) == 0 {
// 		return ""
// 	} else {
// 		return list[0].IP
// 	}
// }

func testIPs(ips []SpeedTestResult) string {
	for _, ip := range ips {
		if checkIPStatus(ip.IP) {
			return ip.IP
		}
	}

	return ""
}

// Test and sort IP addresses based on connection speed
func testAndSortIPs(ips []string) []SpeedTestResult {
	results := []SpeedTestResult{}

	for _, ip := range ips {
		speed, err := testConnectionSpeed(ip)
		if err != nil {
			logger.Errorf("Failed to connect to IP address %s: %s\n", ip, err)
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
func getRandomIP(ips []SpeedTestResult) string {
	if len(ips) == 0 {
		return ""
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
	var totalDuration time.Duration
	j := pingTimes

	for i := 0; i < j; i++ {
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

	averageDuration := totalDuration / time.Duration(j)
	return averageDuration, nil
}
