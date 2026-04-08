package main

import (
	"bitsplitter/subnet"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"text/tabwriter"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	jsonOutput  = flag.Bool("json", false, "Output as JSON")
	showTable   = flag.Bool("table", false, "Show CIDR reference table")
	showVersion = flag.Bool("version", false, "Show version")
	noColor     = flag.Bool("no-color", false, "Disable color output")
)

// ANSI color codes
const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	cyan   = "\033[36m"
	yellow = "\033[33m"
)

var useColor bool

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: bitsplitter [flags] <IP/CIDR> [<IP/CIDR>...]\n")
		fmt.Fprintf(os.Stderr, "       bitsplitter [flags] <IP> <subnet mask>\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  bitsplitter 192.168.0.0/24\n")
		fmt.Fprintf(os.Stderr, "  bitsplitter 2001:db8::/32\n")
		fmt.Fprintf(os.Stderr, "  bitsplitter 192.168.1.0 255.255.255.0\n")
		fmt.Fprintf(os.Stderr, "  bitsplitter 10.0.0.0/8 10.1.0.0/16\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Detect TTY for color support
	useColor = !*noColor && !*jsonOutput && isTTY()

	if *showVersion {
		fmt.Printf("bitsplitter %s (commit: %s, built: %s)\n", version, commit, date)
		return
	}

	if *showTable {
		printCIDRTable()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Single CIDR or mask notation
	if len(args) == 1 || (len(args) == 2 && !containsSlash(args[1])) {
		var input string
		if len(args) == 2 {
			cidr, err := subnet.ParseMaskNotation(args[0], args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			input = cidr
		} else {
			input = args[0]
		}

		info, err := subnet.CalculateSubnet(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if *jsonOutput {
			printJSON(info)
		} else {
			printSubnetInfo(info)
		}
		return
	}

	// Multiple CIDRs — overlap detection
	var infos []*subnet.SubnetInfo
	for _, arg := range args {
		info, err := subnet.CalculateSubnet(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %q: %v\n", arg, err)
			os.Exit(1)
		}
		infos = append(infos, info)
	}

	if *jsonOutput {
		printMultiJSON(infos, args)
	} else {
		printMultiSubnetInfo(infos, args)
	}
}

func containsSlash(s string) bool {
	return strings.Contains(s, "/")
}

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func colorHeader(s string) string {
	if useColor {
		return bold + cyan + s + reset
	}
	return s
}

func colorLabel(s string) string {
	if useColor {
		return yellow + s + reset
	}
	return s
}

// --- Single subnet display ---

func printSubnetInfo(info *subnet.SubnetInfo) {
	fmt.Println()
	fmt.Println(colorHeader("# Overview (IP/CIDR)"))
	fmt.Printf("%s %s/%d\n", colorLabel("IP address:"), info.IP, info.CIDR)

	fmt.Println()
	fmt.Println(colorHeader("# Address Details"))
	fmt.Printf("%s %s\n", colorLabel("Address range:"), info.AddressRange)
	fmt.Printf("%s %s\n", colorLabel("Number of hosts:"), info.NumHosts.String())
	fmt.Println()
	fmt.Printf("%s %s\n", colorLabel("Network address:"), info.NetworkAddr)
	if info.IsIPv6 {
		fmt.Printf("%s %s\n", colorLabel("Last address:"), info.LastAddr)
	} else {
		fmt.Printf("%s %s\n", colorLabel("Broadcast address:"), info.LastAddr)
	}
	fmt.Println()
	fmt.Printf("%s %s\n", colorLabel("Usable range:"), info.UsableRange)
	fmt.Printf("%s %s\n", colorLabel("Usable hosts:"), info.UsableHosts.String())
	fmt.Println()
	fmt.Printf("%s %s\n", colorLabel("Reverse DNS:"), info.ReverseDNS)
	fmt.Printf("%s %s\n", colorLabel("Hex (IP):"), info.HexIP)
	fmt.Printf("%s %s\n", colorLabel("Hex (network):"), info.HexNetwork)

	fmt.Println()
	fmt.Println(colorHeader("# Mask Information"))
	fmt.Printf("%s %s\n", colorLabel("Subnet mask:"), info.SubnetMask)
	fmt.Printf("%s %s\n", colorLabel("Wildcard mask:"), info.WildcardMask)
	fmt.Printf("%s %s\n", colorLabel("Subnet mask (binary):"), info.SubnetMaskBin)
	fmt.Printf("%s %s\n", colorLabel("Wildcard mask (binary):"), info.WildcardMaskBin)

	fmt.Println()
	fmt.Println(colorHeader("# Classification"))
	fmt.Printf("%s %s\n", colorLabel("IP type:"), map[bool]string{true: "Private", false: "Public"}[info.IsPrivate])
	if !info.IsIPv6 {
		fmt.Printf("%s %s\n", colorLabel("IP class:"), info.IPClass)
	}

	fmt.Println()
	fmt.Println(colorHeader("# Supernet"))
	fmt.Printf("%s %s\n", colorLabel("Parent CIDR:"), info.ParentCIDR)
	fmt.Printf("%s %s\n", colorLabel("Parent range:"), info.ParentRange)

	if len(info.Subnets) > 0 {
		fmt.Println()
		fmt.Println(colorHeader("# Subnet Split"))
		for _, s := range info.Subnets {
			fmt.Printf("  %s\n", s)
		}
	}
}

// --- Multi-subnet display with overlap ---

func printMultiSubnetInfo(infos []*subnet.SubnetInfo, cidrs []string) {
	for i, info := range infos {
		fmt.Printf("\n%s\n", colorHeader(fmt.Sprintf("# Network %d: %s/%d", i+1, info.NetworkAddr, info.CIDR)))
		fmt.Printf("  %s %s\n", colorLabel("Range:"), info.AddressRange)
		fmt.Printf("  %s %s\n", colorLabel("Hosts:"), info.NumHosts.String())
		fmt.Printf("  %s %s\n", colorLabel("Usable:"), info.UsableHosts.String())
	}

	results, err := subnet.CheckAllOverlaps(cidrs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking overlaps: %v\n", err)
		return
	}

	fmt.Printf("\n%s\n", colorHeader("# Overlap Analysis"))
	for _, r := range results {
		status := "no overlap"
		if r.Overlaps {
			switch r.Relation {
			case "identical":
				status = "identical networks"
			case "a contains b":
				status = fmt.Sprintf("%s contains %s", r.A, r.B)
			case "b contains a":
				status = fmt.Sprintf("%s contains %s", r.B, r.A)
			}
		}
		fmt.Printf("  %s vs %s: %s\n", r.A, r.B, status)
	}
}

// --- JSON output ---

type jsonSubnetInfo struct {
	IP              string   `json:"ip"`
	CIDR            int      `json:"cidr"`
	IsIPv6          bool     `json:"is_ipv6"`
	AddressRange    string   `json:"address_range"`
	NumHosts        string   `json:"num_hosts"`
	UsableRange     string   `json:"usable_range"`
	NetworkAddr     string   `json:"network_address"`
	LastAddr        string   `json:"last_address"`
	UsableHosts     string   `json:"usable_hosts"`
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

func toJSONInfo(info *subnet.SubnetInfo) jsonSubnetInfo {
	return jsonSubnetInfo{
		IP:              info.IP.String(),
		CIDR:            info.CIDR,
		IsIPv6:          info.IsIPv6,
		AddressRange:    info.AddressRange,
		NumHosts:        info.NumHosts.String(),
		UsableRange:     info.UsableRange,
		NetworkAddr:     info.NetworkAddr.String(),
		LastAddr:        info.LastAddr.String(),
		UsableHosts:     info.UsableHosts.String(),
		SubnetMask:      info.SubnetMask,
		WildcardMask:    info.WildcardMask,
		SubnetMaskBin:   info.SubnetMaskBin,
		WildcardMaskBin: info.WildcardMaskBin,
		IPClass:         info.IPClass,
		IsPrivate:       info.IsPrivate,
		ReverseDNS:      info.ReverseDNS,
		HexIP:           info.HexIP,
		HexNetwork:      info.HexNetwork,
		ParentCIDR:      info.ParentCIDR,
		ParentRange:     info.ParentRange,
		Subnets:         info.Subnets,
	}
}

func printJSON(info *subnet.SubnetInfo) {
	data, _ := json.MarshalIndent(toJSONInfo(info), "", "  ")
	fmt.Println(string(data))
}

func printMultiJSON(infos []*subnet.SubnetInfo, cidrs []string) {
	type multiOutput struct {
		Networks []jsonSubnetInfo      `json:"networks"`
		Overlaps []subnet.OverlapResult `json:"overlaps"`
	}

	out := multiOutput{}
	for _, info := range infos {
		out.Networks = append(out.Networks, toJSONInfo(info))
	}

	results, err := subnet.CheckAllOverlaps(cidrs)
	if err == nil {
		out.Overlaps = results
	}

	data, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(data))
}

// --- CIDR reference table ---

func printCIDRTable() {
	fmt.Println(colorHeader("IPv4 CIDR Reference Table"))
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "CIDR\tSubnet Mask\tAddresses\tUsable Hosts\n")
	fmt.Fprintf(w, "----\t-----------\t---------\t------------\n")

	for cidr := 0; cidr <= 32; cidr++ {
		mask := fmt.Sprintf("%d.%d.%d.%d",
			maskByte(cidr, 0), maskByte(cidr, 1),
			maskByte(cidr, 2), maskByte(cidr, 3))

		addrs := new(big.Int).Lsh(big.NewInt(1), uint(32-cidr))

		var usable string
		switch {
		case cidr == 32:
			usable = "1"
		case cidr == 31:
			usable = "2"
		default:
			u := new(big.Int).Sub(addrs, big.NewInt(2))
			usable = u.String()
		}

		fmt.Fprintf(w, "/%d\t%s\t%s\t%s\n", cidr, mask, addrs.String(), usable)
	}
	w.Flush()

	fmt.Println()
	fmt.Println(colorHeader("IPv6 CIDR Reference Table (common prefixes)"))
	fmt.Println()

	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "CIDR\tAddresses\tTypical Use\n")
	fmt.Fprintf(w, "----\t---------\t-----------\n")

	ipv6Prefixes := []struct {
		cidr int
		use  string
	}{
		{8, "Large allocation"},
		{16, "Large allocation"},
		{32, "ISP allocation"},
		{48, "Site allocation"},
		{56, "Subnet allocation"},
		{64, "Single subnet (standard)"},
		{96, "IPv4-mapped addresses"},
		{112, "Small subnet"},
		{126, "Point-to-point link (/126)"},
		{127, "Point-to-point link (RFC 6164)"},
		{128, "Single host"},
	}

	for _, p := range ipv6Prefixes {
		addrs := new(big.Int).Lsh(big.NewInt(1), uint(128-p.cidr))
		fmt.Fprintf(w, "/%d\t%s\t%s\n", p.cidr, addrs.String(), p.use)
	}
	w.Flush()
}

func maskByte(cidr, bytePos int) byte {
	bits := cidr - bytePos*8
	switch {
	case bits >= 8:
		return 255
	case bits <= 0:
		return 0
	default:
		return byte(0xFF << (8 - bits))
	}
}
