package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func nslookup(host string) ([]string, error) {

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf(`error while looking up host "%s" -> %v`, host, err.Error())
	}

	ipStrings := make([]string, len(ips))
	for i, ip := range ips {
		ipStrings[i] = ip.String()
	}

	return ipStrings, nil
}

func getPortsList(portVar string) ([]int, error) {

	// if port argument is like : 22,80,23
	if strings.Contains(portVar, ",") {
		portsList := strings.Split(portVar, ",")
		portsIntsList := make([]int, len(portsList))
		for p := range portsList {
			port, err := strconv.Atoi(portsList[p])
			if err != nil {
				return nil, fmt.Errorf("invalid Port : %s", portsList[p])
			}
			portsIntsList[p] = port
		}
		return portsIntsList, nil

	} else if strings.Contains(portVar, "-") {
		// if port argument is like : 1-1024

		portMinAndMax := strings.Split(portVar, "-")

		portMin, err := strconv.Atoi(portMinAndMax[0])
		if err != nil {
			return nil, fmt.Errorf("invalid Port Min Range : %s", portMinAndMax[0])
		}

		portMax, err := strconv.Atoi(portMinAndMax[1])
		if err != nil {
			return nil, fmt.Errorf("invalid Port Max Range : %s", portMinAndMax[1])
		}

		var portsTempList []int

		for pMin := portMin; pMin <= portMax; pMin++ {
			portsTempList = append(portsTempList, pMin)
		}

		return portsTempList, nil

	}

	// if port is single number like : 80
	port, err := strconv.Atoi(portVar) // check if port is correct (int)
	if err != nil {
		return nil, fmt.Errorf("invalid Port : %s", portVar)
	}

	return []int{port}, nil
}
