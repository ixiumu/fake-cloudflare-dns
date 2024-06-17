package main

import (
	"github.com/miekg/dns"
	"fmt"
	"flag"
	"regexp"
)

var (
	filename string
	port int
	defaultDomain string
	upstreamServer string
	pingTimes int
	interval int
	fixedIPAddressV4 = ""
	fixedIPAddressV6 = ""
	dnsProxy = DNSProxy{}
	logger *Log
	logLevel string
)

func main() {
	flag.StringVar(&filename, "f", "ip.txt", "File name")
	flag.IntVar(&port, "p", 53, "Port number")
	flag.StringVar(&defaultDomain, "domain", "creativecommons.org", "Default domain")
	flag.IntVar(&pingTimes, "t", 3, "Ping times")
	flag.StringVar(&logLevel, "log", "info", "Log level: none | err | info")
	flag.StringVar(&upstreamServer, "dns", "8.8.8.8:53", "Upstream DNS Server")
	flag.IntVar(&interval, "i", 360, "Each speed measurement interval")
	flag.Parse()

	logger = NewLogger(logLevel)

	logger.Infof("Upstream DNS Server: %s", upstreamServer)

	go func() {
		LookupDefaultIP(defaultDomain)
		testSpeed()
	}()

	// Create a new DNS server
	dns.HandleFunc(".", handleDNSRequest)

	server := &dns.Server{Addr: fmt.Sprintf(":%d", port), Net: "udp"}
	logger.Infof("Starting DNS server on :%d", port)
	err := server.ListenAndServe()
	if err != nil {
		logger.Errorf("Failed to start server: %s\n", err.Error())
	}
	defer server.Shutdown()
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	switch r.Opcode {
	case dns.OpcodeQuery:
		m, err := dnsProxy.getResponse(r)
		if err != nil {
			logger.Errorf("Failed lookup for %s with error: %s\n", r, err.Error())
			m.SetReply(r)
			w.WriteMsg(m)
			return
		}
		if len(m.Answer) > 0 {
			pattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
			ipAddress := pattern.FindAllString(m.Answer[0].String(), -1)

			if len(ipAddress) > 0 {
				logger.Infof("Lookup for %s with ip %s\n", m.Answer[0].Header().Name, ipAddress[0])
			} else {
				logger.Infof("Lookup for %s with response %s\n", m.Answer[0].Header().Name, m.Answer[0])
			}
		}
		m.SetReply(r)
		w.WriteMsg(m)
	}
}
