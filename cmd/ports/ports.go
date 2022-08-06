/*
Copyright © 2022 Jace Walker <jc@jcwlkr.io>
*/
package ports

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
)

func ScanPorts(ipAddress net.IP, portsToScan []int) {

	var openPortsSlice = []int{0}

	for _, port := range portsToScan {
		conn, err := net.Dial("tcp", ipAddress.String()+":"+strconv.Itoa(port))
		if err == nil {
			conn.Close()
			openPortsSlice = append(openPortsSlice, port)
			fmt.Println(ipAddress.String() + ":" + strconv.Itoa(port))
		}
	}
}

func ConvertPortsStringToSlice(portsToScanString string) []int {

	portRange := []int{}

	if strings.Contains(portsToScanString, "-") {
		for _, port := range strings.Split(portsToScanString, "-") {
			portInt, _ := strconv.Atoi(port)
			portRange = append(portRange, portInt)
		}
	} else {
		portInt, _ := strconv.Atoi(portsToScanString)
		portRange = append(portRange, portInt)
	}

	if strings.Contains(portsToScanString, "-") {
		for i := portRange[0]; i <= portRange[1]; i++ {
			portRange = append(portRange, i)
		}

		return sort.IntSlice(portRange[2:])
	} else {
		return portRange
	}

}
