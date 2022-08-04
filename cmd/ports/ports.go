package ports

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
)

// var commonPorts = []int{22, 80, 443}

func ScanPorts(ipAddress net.IP, portsString string) []int {
	var openPortsSlice = []int{}
	portsToScan := convertPortsStringToSlice(portsString)

	for _, port := range portsToScan {
		conn, err := net.Dial("tcp", ipAddress.String()+":"+strconv.Itoa(port))
		if err == nil {
			conn.Close()
			openPortsSlice = append(openPortsSlice, port)
			fmt.Println(ipAddress.String(), "- Port", port, "is open")
		}
	}

	return openPortsSlice
}

func convertPortsStringToSlice(portsToScanString string) []int {
	// if portsToScan is a single port, return a slice with one port
	// if portsToScan is a range, return a slice with all ports in the range
	// portsToScan = "0-1000"
	// split on -
	// convert to int

	portRange := []int{}

	for _, port := range strings.Split(portsToScanString, "-") {
		portInt, _ := strconv.Atoi(port)
		portRange = append(portRange, portInt)
	}

	// add all integers between the two ports to the slice
	for i := portRange[0]; i <= portRange[1]; i++ {
		portRange = append(portRange, i)
	}

	// sort portRange in ascending order
	return sort.IntSlice(portRange[2:])
}
