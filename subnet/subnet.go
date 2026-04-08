package subnet

import (
	"fmt"
	"math/big"
	"net"
)

// SubnetInfo holds all calculated subnet info
type SubnetInfo struct {
	IP              net.IP
	CIDR            int
	IsIPv6          bool
	AddressRange    string
	NumHosts        *big.Int
	UsableRange     string
	NetworkAddr     net.IP
	LastAddr        net.IP
	UsableHosts     *big.Int
	SubnetMaskBin   string
	WildcardMaskBin string
	IPClass         string
	IsPrivate       bool
}

// Computes subnet details from a CIDR notation
func CalculateSubnet(input string) (*SubnetInfo, error) {
	_, ipNet, err := net.ParseCIDR(input)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation: %v", err)
	}

	ip := ipNet.IP
	mask := ipNet.Mask
	cidr, totalBits := mask.Size()

	// Detect IPv4 vs IPv6
	isIPv6 := ip.To4() == nil

	// Normalize IP to canonical form
	if !isIPv6 {
		ip = ip.To4()
		mask = mask[:4] // ensure 4-byte mask for IPv4
	}

	// Calculate number of hosts
	hostBits := totalBits - cidr
	numHosts := new(big.Int).Lsh(big.NewInt(1), uint(hostBits))

	// Calculate network and last addresses
	network := ip.Mask(mask)
	lastAddr := make(net.IP, len(ip))
	for i := range ip {
		lastAddr[i] = network[i] | ^mask[i]
	}

	// Calculate usable hosts and range
	usableHosts := new(big.Int)
	var usableRange string

	if isIPv6 {
		switch {
		case cidr == 128:
			usableHosts.SetInt64(1)
			usableRange = network.String()
		case cidr == 127:
			usableHosts.SetInt64(2)
			usableRange = fmt.Sprintf("%s - %s", network, lastAddr)
		default:
			usableHosts.Set(numHosts)
			usableRange = fmt.Sprintf("%s - %s", network, lastAddr)
		}
	} else {
		switch {
		case cidr == 32:
			usableHosts.SetInt64(1)
			usableRange = network.String()
		case cidr == 31:
			usableHosts.SetInt64(2)
			usableRange = fmt.Sprintf("%s - %s", network, lastAddr)
		default:
			usableHosts.Sub(numHosts, big.NewInt(2))
			usableFirst := incrementIP(network)
			usableLast := decrementIP(lastAddr)
			usableRange = fmt.Sprintf("%s - %s", usableFirst, usableLast)
		}
	}

	// Convert masks to binary
	subnetMaskBin := maskToBinary(mask)
	wildcardMaskBin := maskToBinary(invertMask(mask))

	// Determine IP class (IPv4 only)
	ipClass := "N/A"
	if !isIPv6 {
		ipClass = getIPClass(ip)
	}

	info := &SubnetInfo{
		IP:              ip,
		CIDR:            cidr,
		IsIPv6:          isIPv6,
		AddressRange:    fmt.Sprintf("%s - %s", network, lastAddr),
		NumHosts:        numHosts,
		UsableRange:     usableRange,
		NetworkAddr:     network,
		LastAddr:        lastAddr,
		UsableHosts:     usableHosts,
		SubnetMaskBin:   subnetMaskBin,
		WildcardMaskBin: wildcardMaskBin,
		IPClass:         ipClass,
		IsPrivate:       ip.IsPrivate(),
	}

	return info, nil
}

// Converts mask to binary string format
// IPv4: dot-separated 8-bit groups (11111111.11111111.11000000.00000000)
// IPv6: colon-separated 16-bit groups
func maskToBinary(mask net.IPMask) string {
	if len(mask) == 4 {
		binary := ""
		for i, b := range mask {
			binary += fmt.Sprintf("%08b", b)
			if i < len(mask)-1 {
				binary += "."
			}
		}
		return binary
	}

	// IPv6: colon-separated 16-bit groups
	binary := ""
	for i := 0; i < len(mask); i += 2 {
		binary += fmt.Sprintf("%08b%08b", mask[i], mask[i+1])
		if i < len(mask)-2 {
			binary += ":"
		}
	}
	return binary
}

// Creates a wildcard mask by inverting the subnet mask
func invertMask(mask net.IPMask) net.IPMask {
	inverted := make(net.IPMask, len(mask))
	for i, b := range mask {
		inverted[i] = ^b
	}
	return inverted
}

// Increments an IP address by 1
func incrementIP(ip net.IP) net.IP {
	newIP := make(net.IP, len(ip))
	copy(newIP, ip)
	for i := len(newIP) - 1; i >= 0; i-- {
		if newIP[i] < 255 {
			newIP[i]++
			break
		}
		newIP[i] = 0
	}
	return newIP
}

// Decrements an IP address by 1
func decrementIP(ip net.IP) net.IP {
	newIP := make(net.IP, len(ip))
	copy(newIP, ip)
	for i := len(newIP) - 1; i >= 0; i-- {
		if newIP[i] > 0 {
			newIP[i]--
			break
		}
		newIP[i] = 255
	}
	return newIP
}

// Determines the IP class (A, B, C, D, E) based on the first octet
func getIPClass(ip net.IP) string {
	firstOctet := ip[0]
	switch {
	case firstOctet <= 127:
		return "A"
	case firstOctet <= 191:
		return "B"
	case firstOctet <= 223:
		return "C"
	case firstOctet <= 239:
		return "D"
	default:
		return "E"
	}
}
