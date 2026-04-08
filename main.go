package main

import (
	"bitsplitter/subnet"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: bitsplitter <IP/CIDR>")
		fmt.Println("Examples:")
		fmt.Println("  bitsplitter 192.168.0.0/24")
		fmt.Println("  bitsplitter 2001:db8::/32")
		os.Exit(1)
	}

	info, err := subnet.CalculateSubnet(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n# Overview (IP/CIDR)")
	fmt.Printf("IP address: %s/%d\n", info.IP, info.CIDR)

	fmt.Println("\n# Address Details")
	fmt.Printf("Address range: %s\n", info.AddressRange)
	fmt.Printf("Number of hosts: %s\n\n", info.NumHosts.String())
	fmt.Printf("Network address: %s\n", info.NetworkAddr)

	if info.IsIPv6 {
		fmt.Printf("Last address: %s\n\n", info.LastAddr)
	} else {
		fmt.Printf("Broadcast address: %s\n\n", info.LastAddr)
	}

	fmt.Printf("Usable range: %s\n", info.UsableRange)
	fmt.Printf("Usable hosts: %s\n", info.UsableHosts.String())

	fmt.Println("\n# Mask Information")
	fmt.Printf("Subnet mask (binary): %s\n", info.SubnetMaskBin)
	fmt.Printf("Wildcard mask (binary): %s\n", info.WildcardMaskBin)

	fmt.Println("\n# Classification")
	fmt.Printf("IP type: %s\n", map[bool]string{true: "Private", false: "Public"}[info.IsPrivate])
	if !info.IsIPv6 {
		fmt.Printf("IP class: %s\n", info.IPClass)
	}
}
