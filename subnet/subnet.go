package subnet

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

// SubnetInfo holds all calculated subnet info
type SubnetInfo struct {
	IP              net.IP   `json:"-"`
	CIDR            int      `json:"cidr"`
	IsIPv6          bool     `json:"is_ipv6"`
	AddressRange    string   `json:"address_range"`
	NumHosts        *big.Int `json:"-"`
	UsableRange     string   `json:"usable_range"`
	NetworkAddr     net.IP   `json:"-"`
	LastAddr        net.IP   `json:"-"`
	UsableHosts     *big.Int `json:"-"`
	SubnetMask      string   `json:"subnet_mask"`
	WildcardMask    string   `json:"wildcard_mask"`
	SubnetMaskBin   string   `json:"subnet_mask_binary"`
	WildcardMaskBin string   `json:"wildcard_mask_binary"`
	IPClass         string   `json:"ip_class,omitempty"`
	IsPrivate       bool     `json:"is_private"`
	ReverseDNS      string   `json:"reverse_dns"`
	HexIP           string   `json:"hex_ip"`
	HexNetwork      string   `json:"hex_network"`
	ParentCIDR      string   `json:"parent_cidr"`
	ParentRange     string   `json:"parent_range"`
	Subnets         []string `json:"subnets,omitempty"`
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
		mask = mask[:4]
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

	// Convert masks to binary and decimal
	subnetMaskBin := maskToBinary(mask)
	wildcardMaskBin := maskToBinary(invertMask(mask))
	subnetMask := maskToDecimal(mask, isIPv6)
	wildcardMask := maskToDecimal(invertMask(mask), isIPv6)

	// Determine IP class (IPv4 only)
	ipClass := ""
	if !isIPv6 {
		ipClass = getIPClass(ip)
	}

	// Reverse DNS
	reverseDNS := computeReverseDNS(network, cidr, isIPv6)

	// Hex representations
	hexIP := ipToHex(ip)
	hexNetwork := ipToHex(network)

	// Parent network
	parentCIDR, parentRange := computeParent(network, cidr, totalBits, isIPv6)

	// Subnet split
	subnets := computeSubnets(network, cidr, totalBits, isIPv6)

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
		SubnetMask:      subnetMask,
		WildcardMask:    wildcardMask,
		SubnetMaskBin:   subnetMaskBin,
		WildcardMaskBin: wildcardMaskBin,
		IPClass:         ipClass,
		IsPrivate:       ip.IsPrivate(),
		ReverseDNS:      reverseDNS,
		HexIP:           hexIP,
		HexNetwork:      hexNetwork,
		ParentCIDR:      parentCIDR,
		ParentRange:     parentRange,
		Subnets:         subnets,
	}

	return info, nil
}

// ParseMaskNotation converts "IP mask" notation to CIDR notation
func ParseMaskNotation(ipStr, maskStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address: %s", ipStr)
	}

	maskIP := net.ParseIP(maskStr)
	if maskIP == nil {
		return "", fmt.Errorf("invalid subnet mask: %s", maskStr)
	}

	maskIPv4 := maskIP.To4()
	if maskIPv4 == nil {
		return "", fmt.Errorf("subnet mask must be IPv4: %s", maskStr)
	}

	mask := net.IPMask(maskIPv4)
	ones, bits := mask.Size()
	if bits == 0 {
		return "", fmt.Errorf("non-contiguous subnet mask: %s", maskStr)
	}

	return fmt.Sprintf("%s/%d", ipStr, ones), nil
}

// Converts mask to decimal string format
func maskToDecimal(mask net.IPMask, isIPv6 bool) string {
	if !isIPv6 {
		parts := make([]string, len(mask))
		for i, b := range mask {
			parts[i] = fmt.Sprintf("%d", b)
		}
		return strings.Join(parts, ".")
	}

	// IPv6: colon-separated hex groups
	parts := make([]string, 0, 8)
	for i := 0; i < len(mask); i += 2 {
		parts = append(parts, fmt.Sprintf("%02x%02x", mask[i], mask[i+1]))
	}
	return strings.Join(parts, ":")
}

// Converts mask to binary string format
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

// Determines the IP class based on the first octet
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

// Converts an IP address to hex representation
func ipToHex(ip net.IP) string {
	hex := "0x"
	for _, b := range ip {
		hex += fmt.Sprintf("%02X", b)
	}
	return hex
}

// Computes the reverse DNS zone string
func computeReverseDNS(network net.IP, cidr int, isIPv6 bool) string {
	if isIPv6 {
		return computeReverseDNSv6(network, cidr)
	}
	return computeReverseDNSv4(network, cidr)
}

func computeReverseDNSv4(network net.IP, cidr int) string {
	octets := make([]string, 4)
	for i, b := range network {
		octets[i] = fmt.Sprintf("%d", b)
	}

	// Number of full octets in the network portion
	fullOctets := cidr / 8
	if fullOctets == 0 {
		return "in-addr.arpa"
	}

	// Reverse the network octets
	reversed := make([]string, fullOctets)
	for i := 0; i < fullOctets; i++ {
		reversed[i] = octets[fullOctets-1-i]
	}

	zone := strings.Join(reversed, ".") + ".in-addr.arpa"

	// For non-octet-boundary CIDRs, prepend the partial octet with CIDR
	if cidr%8 != 0 {
		zone = fmt.Sprintf("%s/%d.%s", octets[fullOctets], cidr, zone)
	}

	return zone
}

func computeReverseDNSv6(network net.IP, cidr int) string {
	// Expand to full 32 hex nibbles
	full := fmt.Sprintf("%032x", []byte(network))

	// Use nibble boundary (round up to nearest multiple of 4)
	nibbles := cidr / 4
	if cidr%4 != 0 {
		nibbles++
	}

	if nibbles == 0 {
		return "ip6.arpa"
	}

	// Take the network nibbles and reverse them
	chars := make([]string, nibbles)
	for i := 0; i < nibbles; i++ {
		chars[i] = string(full[nibbles-1-i])
	}

	return strings.Join(chars, ".") + ".ip6.arpa"
}

// Computes the parent (supernet) network
func computeParent(network net.IP, cidr, totalBits int, isIPv6 bool) (string, string) {
	if cidr == 0 {
		return "N/A", "N/A"
	}

	parentCIDR := cidr - 1
	parentMask := net.CIDRMask(parentCIDR, totalBits)

	// Normalize mask length to match IP length
	if !isIPv6 {
		parentMask = parentMask[len(parentMask)-4:]
	}

	parentNet := network.Mask(parentMask)
	parentLast := make(net.IP, len(parentNet))
	for i := range parentNet {
		parentLast[i] = parentNet[i] | ^parentMask[i]
	}

	cidrStr := fmt.Sprintf("%s/%d", parentNet, parentCIDR)
	rangeStr := fmt.Sprintf("%s - %s", parentNet, parentLast)
	return cidrStr, rangeStr
}

// Computes the two subnets from splitting the current prefix
func computeSubnets(network net.IP, cidr, totalBits int, isIPv6 bool) []string {
	maxCIDR := 32
	if isIPv6 {
		maxCIDR = 128
	}
	if cidr >= maxCIDR {
		return nil
	}

	newCIDR := cidr + 1

	// First subnet: same network address
	first := fmt.Sprintf("%s/%d", network, newCIDR)

	// Second subnet: flip the bit at position `cidr` (0-indexed from MSB)
	secondNet := make(net.IP, len(network))
	copy(secondNet, network)
	byteIndex := cidr / 8
	bitIndex := uint(7 - cidr%8)
	secondNet[byteIndex] |= 1 << bitIndex

	second := fmt.Sprintf("%s/%d", secondNet, newCIDR)

	return []string{first, second}
}
