package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type WrongAnswerError struct {
	Message string
}

func (e *WrongAnswerError) Error() string {
	return fmt.Sprintf("Wrong Answer: %s", e.Message)
}

func SendQuery(addr string, dn string) (string, error) {
	var (
		domain  string
		rdns_ip string
	)
	if dn == "timestamp" {
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		domain = strings.Join([]string{timestamp, "-scan.echodns.xyz."}, "")
	} else {
		domain = strings.Join([]string{dn, ".echodns.xyz."}, "")
	}
	//fmt.Println(domain)
	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeA)
	m.RecursionDesired = true

	res, err := dns.Exchange(m, addr)
	if err == nil {
		if len(res.Answer) == 1 {
			if a, ok := res.Answer[0].(*dns.A); ok {
				rdns_ip = a.A.String()
			} else {
				rdns_ip = ""
				err = &WrongAnswerError{
					Message: "Wrong Record Type",
				}
			}
		} else {
			rdns_ip = ""
			err = &WrongAnswerError{
				Message: "Wrong Answer Section",
			}
		}
	}
	return rdns_ip, err
}
