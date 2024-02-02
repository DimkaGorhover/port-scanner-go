package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type scanner struct {
	ip      string
	port    int
	timeout time.Duration
	debug   bool
}

func (s *scanner) scan(ctx context.Context) error {
	d := net.Dialer{Timeout: s.timeout}
	address := fmt.Sprintf(`%s:%d`, s.ip, s.port)
	if s.debug {
		fmt.Printf("scanning %s ...\n", address)
	}
	_, err := d.Dial("tcp", address)
	if err != nil {
		if addrError, ok := err.(*net.AddrError); ok {
			if addrError.Timeout() {
				return nil
			}
		} else if opError, ok := err.(*net.OpError); ok {

			// handle lacked sufficient buffer space error

			if strings.TrimSpace(opError.Err.Error()) == "bind: An operation on a socket could not be performed because "+
				"the system lacked sufficient buffer space or because a queue was full." {

				time.Sleep(s.timeout + (3 * time.Second))

				_, errAe := d.Dial("tcp", address)

				if errAe != nil {
					if addErr, ok := err.(*net.AddrError); ok {
						if addErr.Timeout() {
							return nil
						}
					}
				}
			}

		} else {
			return err
		}

		return nil
	}

	fmt.Printf("[+] Port %15s %5d/TCP is open\n", s.ip, s.port)

	return nil
}
