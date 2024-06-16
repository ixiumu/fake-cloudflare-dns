package main

import (
	"log"
	"net"
	"time"
	"flag"

	"github.com/miekg/dns"
)

type SpeedTestResult struct {
	IP     string
	Speed  time.Duration
}

var (
	filename   string
	port       string
	IPs        []string
	fastestIPs []SpeedTestResult
)

func main() {
	flag.StringVar(&filename, "f", "ip.txt", "File name")
	flag.StringVar(&port, "p", "53", "Port number")
	flag.Parse()

	log.Printf("Listen: %s", port)

	server := dns.Server{
		Addr: ":" + string(port),
		Net:  "udp",
	}

	server.Handler = dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {

		resp := new(dns.Msg)
		resp.SetReply(r)
		ip := getIP(fastestIPs)

		for _, q := range r.Question {
			if (isIPv6(ip)) {
				rr := new(dns.AAAA)
				rr.Hdr = dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    3600,
				}
				rr.AAAA = net.ParseIP(ip)
				resp.Answer = append(resp.Answer, rr)
			} else {
				rr := new(dns.A)
				rr.Hdr = dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    3600,
				}
				rr.A = net.ParseIP(ip)
				resp.Answer = append(resp.Answer, rr)
			}
		}

		log.Printf("IP: %s", ip)

		err := w.WriteMsg(resp)
		if err != nil {
			log.Printf("Failed to send DNS response: %s", err)
		}
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("DNS server start sucess")
		}
	}()

	testSpeed()
}

func testSpeed() {
	IPs = LookupIP(filename)

	// Test connection speed and sort
	fastestIPs = testAndSortIPs(IPs)	// Print sorted results

	log.Printf("Connection speed sorting results:")
	for _, result := range fastestIPs {
		log.Printf("IP: %s, Speed: %s\n", result.IP, result.Speed)
	}

	ip := getIP(fastestIPs)
	log.Printf(ip)

	ticker := time.NewTicker(6 * time.Hour)
	for {
		<-ticker.C
		fastestIPs = testAndSortIPs(IPs)	// Print sorted results
	}
}
