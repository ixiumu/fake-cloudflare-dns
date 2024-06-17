package main

// https://github.com/katakonst/go-dns-proxy/blob/master/dns_proxy.go

import (
	"fmt"
	"net"
	"regexp"

	"github.com/miekg/dns"
)

type DNSProxy struct {
	domains       map[string]interface{}
	servers       map[string]interface{}
	defaultServer string
}

func (proxy *DNSProxy) getResponse(requestMsg *dns.Msg) (*dns.Msg, error) {
	responseMsg := new(dns.Msg)
	if len(requestMsg.Question) > 0 {
		question := requestMsg.Question[0]

		dnsServer := upstreamServer

		// valid := false
		queryMsg := new(dns.Msg)
		requestMsg.CopyTo(queryMsg)
		queryMsg.Question = []dns.Question{question}

		msg, err := lookup(upstreamServer, queryMsg)
		if err != nil {
			return nil, err
		}

		if len(msg.Answer) == 0 {
			return responseMsg, fmt.Errorf("not found")
		}

		switch question.Qtype {
		case dns.TypeA:
			answer, err := proxy.processTypeA(dnsServer, &question, requestMsg)
			if err != nil {
				return responseMsg, err
			}
			responseMsg.Answer = append(responseMsg.Answer, *answer)

		case dns.TypeAAAA:
			answer, err := proxy.processTypeAAAA(dnsServer, &question, requestMsg)
			if err != nil {
				return responseMsg, err
			}
			responseMsg.Answer = append(responseMsg.Answer, *answer)

		default:
			answer, err := proxy.processOtherTypes(dnsServer, &question, requestMsg)
			if err != nil {
				return responseMsg, err
			}
			responseMsg.Answer = append(responseMsg.Answer, *answer)
		}
	}

	return responseMsg, nil
}

func (proxy *DNSProxy) processOtherTypes(dnsServer string, q *dns.Question, requestMsg *dns.Msg) (*dns.RR, error) {
	queryMsg := new(dns.Msg)
	requestMsg.CopyTo(queryMsg)
	queryMsg.Question = []dns.Question{*q}

	msg, err := lookup(dnsServer, queryMsg)
	if err != nil {
		return nil, err
	}

	if len(msg.Answer) > 0 {
		return &msg.Answer[0], nil
	}
	return nil, fmt.Errorf("not found")
}

func (proxy *DNSProxy) processTypeA(dnsServer string, q *dns.Question, requestMsg *dns.Msg) (*dns.RR, error) {
	ip := fixedIPAddressV4

	if ip == "" {
		return nil, fmt.Errorf("not found")
	}

	answer, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

func (proxy *DNSProxy) processTypeAAAA(dnsServer string, q *dns.Question, requestMsg *dns.Msg) (*dns.RR, error) {
	ip := fixedIPAddressV6
	// ip := "::1"

	if ip == "" {
		return nil, fmt.Errorf("not found")
	}

	answer, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, ip))
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

func (dnsProxy *DNSProxy) getIPFromConfigs(domain string, configs map[string]interface{}) string {

	for k, v := range configs {
		match, _ := regexp.MatchString(k+"\\.", domain)
		if match {
			return v.(string)
		}
	}
	return ""
}

func GetOutboundIP() (net.IP, error) {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func lookup(server string, m *dns.Msg) (*dns.Msg, error) {
	dnsClient := new(dns.Client)
	dnsClient.Net = "udp"
	response, _, err := dnsClient.Exchange(m, server)
	if err != nil {
		return nil, err
	}

	return response, nil
}