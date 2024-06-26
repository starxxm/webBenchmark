package main

import (
	"fmt"
	"github.com/miekg/dns"
)

func nslookup(targetAddress, server string) (res []string) {
	if server == "" {
		server = "8.8.8.8"
	}
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(targetAddress+".", dns.TypeA)

	ns := server + ":53"
	r, t, err := c.Exchange(&m, ns)
	if err != nil {
		fmt.Printf("nameserver %s error: %v\n", ns, err)
		return res
	}
	fmt.Printf("nameserver %s took %v", ns, t)
	if len(r.Answer) == 0 {
		return res
	}
	for _, ans := range r.Answer {
		if ans.Header().Rrtype == dns.TypeA {
			Arecord := ans.(*dns.A)
			res = append(res, fmt.Sprintf("%s", Arecord))
		}
	}
	return
}
