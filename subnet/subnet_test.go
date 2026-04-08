package subnet

import (
	"testing"
)

func TestCalculateSubnetIPv4(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantCIDR      int
		wantNumHosts  string
		wantUsable    string
		wantNetwork   string
		wantBroadcast string
		wantClass     string
		wantPrivate   bool
	}{
		{
			name:          "/24 private class C",
			input:         "192.168.1.0/24",
			wantCIDR:      24,
			wantNumHosts:  "256",
			wantUsable:    "254",
			wantNetwork:   "192.168.1.0",
			wantBroadcast: "192.168.1.255",
			wantClass:     "C",
			wantPrivate:   true,
		},
		{
			name:          "/16 private class B",
			input:         "172.16.0.0/16",
			wantCIDR:      16,
			wantNumHosts:  "65536",
			wantUsable:    "65534",
			wantNetwork:   "172.16.0.0",
			wantBroadcast: "172.16.255.255",
			wantClass:     "B",
			wantPrivate:   true,
		},
		{
			name:          "/8 private class A",
			input:         "10.0.0.0/8",
			wantCIDR:      8,
			wantNumHosts:  "16777216",
			wantUsable:    "16777214",
			wantNetwork:   "10.0.0.0",
			wantBroadcast: "10.255.255.255",
			wantClass:     "A",
			wantPrivate:   true,
		},
		{
			name:          "/32 single host",
			input:         "192.168.1.1/32",
			wantCIDR:      32,
			wantNumHosts:  "1",
			wantUsable:    "1",
			wantNetwork:   "192.168.1.1",
			wantBroadcast: "192.168.1.1",
			wantClass:     "C",
			wantPrivate:   true,
		},
		{
			name:          "/31 point-to-point",
			input:         "192.168.1.0/31",
			wantCIDR:      31,
			wantNumHosts:  "2",
			wantUsable:    "2",
			wantNetwork:   "192.168.1.0",
			wantBroadcast: "192.168.1.1",
			wantClass:     "C",
			wantPrivate:   true,
		},
		{
			name:          "public class A",
			input:         "8.8.8.0/24",
			wantCIDR:      24,
			wantNumHosts:  "256",
			wantUsable:    "254",
			wantNetwork:   "8.8.8.0",
			wantBroadcast: "8.8.8.255",
			wantClass:     "A",
			wantPrivate:   false,
		},
		{
			name:          "class D multicast",
			input:         "224.0.0.0/4",
			wantCIDR:      4,
			wantNumHosts:  "268435456",
			wantUsable:    "268435454",
			wantNetwork:   "224.0.0.0",
			wantBroadcast: "239.255.255.255",
			wantClass:     "D",
			wantPrivate:   false,
		},
		{
			name:          "class E reserved",
			input:         "240.0.0.0/4",
			wantCIDR:      4,
			wantNumHosts:  "268435456",
			wantUsable:    "268435454",
			wantNetwork:   "240.0.0.0",
			wantBroadcast: "255.255.255.255",
			wantClass:     "E",
			wantPrivate:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := CalculateSubnet(tt.input)
			if err != nil {
				t.Fatalf("CalculateSubnet(%q) error: %v", tt.input, err)
			}

			if info.CIDR != tt.wantCIDR {
				t.Errorf("CIDR = %d, want %d", info.CIDR, tt.wantCIDR)
			}
			if info.NumHosts.String() != tt.wantNumHosts {
				t.Errorf("NumHosts = %s, want %s", info.NumHosts.String(), tt.wantNumHosts)
			}
			if info.UsableHosts.String() != tt.wantUsable {
				t.Errorf("UsableHosts = %s, want %s", info.UsableHosts.String(), tt.wantUsable)
			}
			if info.NetworkAddr.String() != tt.wantNetwork {
				t.Errorf("NetworkAddr = %s, want %s", info.NetworkAddr, tt.wantNetwork)
			}
			if info.LastAddr.String() != tt.wantBroadcast {
				t.Errorf("LastAddr = %s, want %s", info.LastAddr, tt.wantBroadcast)
			}
			if info.IPClass != tt.wantClass {
				t.Errorf("IPClass = %s, want %s", info.IPClass, tt.wantClass)
			}
			if info.IsPrivate != tt.wantPrivate {
				t.Errorf("IsPrivate = %v, want %v", info.IsPrivate, tt.wantPrivate)
			}
			if info.IsIPv6 {
				t.Error("IsIPv6 = true, want false")
			}
		})
	}
}

func TestCalculateSubnetIPv4BinaryMask(t *testing.T) {
	info, err := CalculateSubnet("192.168.1.0/24")
	if err != nil {
		t.Fatal(err)
	}
	want := "11111111.11111111.11111111.00000000"
	if info.SubnetMaskBin != want {
		t.Errorf("SubnetMaskBin = %s, want %s", info.SubnetMaskBin, want)
	}
	wantWild := "00000000.00000000.00000000.11111111"
	if info.WildcardMaskBin != wantWild {
		t.Errorf("WildcardMaskBin = %s, want %s", info.WildcardMaskBin, wantWild)
	}
}

func TestCalculateSubnetIPv6(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantCIDR     int
		wantNumHosts string
		wantUsable   string
		wantNetwork  string
		wantLastAddr string
		wantPrivate  bool
	}{
		{
			name:         "/32 documentation prefix",
			input:        "2001:db8::/32",
			wantCIDR:     32,
			wantNumHosts: "79228162514264337593543950336", // 2^96
			wantUsable:   "79228162514264337593543950336",
			wantNetwork:  "2001:db8::",
			wantLastAddr: "2001:db8:ffff:ffff:ffff:ffff:ffff:ffff",
			wantPrivate:  false,
		},
		{
			name:         "/128 single host",
			input:        "::1/128",
			wantCIDR:     128,
			wantNumHosts: "1",
			wantUsable:   "1",
			wantNetwork:  "::1",
			wantLastAddr: "::1",
			wantPrivate:  false,
		},
		{
			name:         "/64 common allocation",
			input:        "2001:db8:abcd::/64",
			wantCIDR:     64,
			wantNumHosts: "18446744073709551616", // 2^64
			wantUsable:   "18446744073709551616",
			wantNetwork:  "2001:db8:abcd::",
			wantLastAddr: "2001:db8:abcd:0:ffff:ffff:ffff:ffff",
			wantPrivate:  false,
		},
		{
			name:         "/127 point-to-point",
			input:        "2001:db8::/127",
			wantCIDR:     127,
			wantNumHosts: "2",
			wantUsable:   "2",
			wantNetwork:  "2001:db8::",
			wantLastAddr: "2001:db8::1",
			wantPrivate:  false,
		},
		{
			name:         "/48 standard site allocation",
			input:        "2001:db8:abcd::/48",
			wantCIDR:     48,
			wantNumHosts: "1208925819614629174706176", // 2^80
			wantUsable:   "1208925819614629174706176",
			wantNetwork:  "2001:db8:abcd::",
			wantLastAddr: "2001:db8:abcd:ffff:ffff:ffff:ffff:ffff",
			wantPrivate:  false,
		},
		{
			name:         "ULA private",
			input:        "fd00::/8",
			wantCIDR:     8,
			wantNumHosts: "1329227995784915872903807060280344576", // 2^120
			wantUsable:   "1329227995784915872903807060280344576",
			wantNetwork:  "fd00::",
			wantLastAddr: "fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			wantPrivate:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := CalculateSubnet(tt.input)
			if err != nil {
				t.Fatalf("CalculateSubnet(%q) error: %v", tt.input, err)
			}

			if !info.IsIPv6 {
				t.Error("IsIPv6 = false, want true")
			}
			if info.CIDR != tt.wantCIDR {
				t.Errorf("CIDR = %d, want %d", info.CIDR, tt.wantCIDR)
			}
			if info.NumHosts.String() != tt.wantNumHosts {
				t.Errorf("NumHosts = %s, want %s", info.NumHosts.String(), tt.wantNumHosts)
			}
			if info.UsableHosts.String() != tt.wantUsable {
				t.Errorf("UsableHosts = %s, want %s", info.UsableHosts.String(), tt.wantUsable)
			}
			if info.NetworkAddr.String() != tt.wantNetwork {
				t.Errorf("NetworkAddr = %s, want %s", info.NetworkAddr, tt.wantNetwork)
			}
			if info.LastAddr.String() != tt.wantLastAddr {
				t.Errorf("LastAddr = %s, want %s", info.LastAddr, tt.wantLastAddr)
			}
			if info.IsPrivate != tt.wantPrivate {
				t.Errorf("IsPrivate = %v, want %v", info.IsPrivate, tt.wantPrivate)
			}
			if info.IPClass != "" {
				t.Errorf("IPClass = %q, want empty for IPv6", info.IPClass)
			}
		})
	}
}

func TestCalculateSubnetIPv6BinaryMask(t *testing.T) {
	info, err := CalculateSubnet("2001:db8::/32")
	if err != nil {
		t.Fatal(err)
	}
	want := "1111111111111111:1111111111111111:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000"
	if info.SubnetMaskBin != want {
		t.Errorf("SubnetMaskBin = %s, want %s", info.SubnetMaskBin, want)
	}
}

func TestDecimalMask(t *testing.T) {
	tests := []struct {
		input       string
		wantSubnet  string
		wantWild    string
	}{
		{"192.168.1.0/24", "255.255.255.0", "0.0.0.255"},
		{"10.0.0.0/8", "255.0.0.0", "0.255.255.255"},
		{"172.16.0.0/16", "255.255.0.0", "0.0.255.255"},
		{"192.168.1.0/32", "255.255.255.255", "0.0.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			info, err := CalculateSubnet(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if info.SubnetMask != tt.wantSubnet {
				t.Errorf("SubnetMask = %s, want %s", info.SubnetMask, tt.wantSubnet)
			}
			if info.WildcardMask != tt.wantWild {
				t.Errorf("WildcardMask = %s, want %s", info.WildcardMask, tt.wantWild)
			}
		})
	}
}

func TestReverseDNS(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"192.168.1.0/24", "1.168.192.in-addr.arpa"},
		{"10.0.0.0/8", "10.in-addr.arpa"},
		{"172.16.0.0/16", "16.172.in-addr.arpa"},
		{"192.168.1.128/25", "128/25.1.168.192.in-addr.arpa"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			info, err := CalculateSubnet(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if info.ReverseDNS != tt.want {
				t.Errorf("ReverseDNS = %s, want %s", info.ReverseDNS, tt.want)
			}
		})
	}
}

func TestReverseDNSv6(t *testing.T) {
	info, err := CalculateSubnet("2001:db8::/32")
	if err != nil {
		t.Fatal(err)
	}
	want := "8.b.d.0.1.0.0.2.ip6.arpa"
	if info.ReverseDNS != want {
		t.Errorf("ReverseDNS = %s, want %s", info.ReverseDNS, want)
	}
}

func TestHexRepresentation(t *testing.T) {
	info, err := CalculateSubnet("192.168.1.0/24")
	if err != nil {
		t.Fatal(err)
	}
	if info.HexIP != "0xC0A80100" {
		t.Errorf("HexIP = %s, want 0xC0A80100", info.HexIP)
	}
	if info.HexNetwork != "0xC0A80100" {
		t.Errorf("HexNetwork = %s, want 0xC0A80100", info.HexNetwork)
	}
}

func TestParentNetwork(t *testing.T) {
	tests := []struct {
		input      string
		wantParent string
	}{
		{"192.168.1.0/24", "192.168.0.0/23"},
		{"10.0.0.0/8", "10.0.0.0/7"},
		{"0.0.0.0/0", "N/A"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			info, err := CalculateSubnet(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if info.ParentCIDR != tt.wantParent {
				t.Errorf("ParentCIDR = %s, want %s", info.ParentCIDR, tt.wantParent)
			}
		})
	}
}

func TestSubnetSplit(t *testing.T) {
	info, err := CalculateSubnet("192.168.1.0/24")
	if err != nil {
		t.Fatal(err)
	}
	if len(info.Subnets) != 2 {
		t.Fatalf("expected 2 subnets, got %d", len(info.Subnets))
	}
	if info.Subnets[0] != "192.168.1.0/25" {
		t.Errorf("Subnets[0] = %s, want 192.168.1.0/25", info.Subnets[0])
	}
	if info.Subnets[1] != "192.168.1.128/25" {
		t.Errorf("Subnets[1] = %s, want 192.168.1.128/25", info.Subnets[1])
	}
}

func TestSubnetSplitMaxCIDR(t *testing.T) {
	info, err := CalculateSubnet("192.168.1.1/32")
	if err != nil {
		t.Fatal(err)
	}
	if len(info.Subnets) != 0 {
		t.Errorf("expected no subnets for /32, got %d", len(info.Subnets))
	}
}

func TestParseMaskNotation(t *testing.T) {
	tests := []struct {
		ip, mask string
		want     string
		wantErr  bool
	}{
		{"192.168.1.0", "255.255.255.0", "192.168.1.0/24", false},
		{"10.0.0.0", "255.0.0.0", "10.0.0.0/8", false},
		{"172.16.0.0", "255.255.0.0", "172.16.0.0/16", false},
		{"192.168.1.0", "invalid", "", true},
		{"invalid", "255.255.255.0", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.ip+" "+tt.mask, func(t *testing.T) {
			got, err := ParseMaskNotation(tt.ip, tt.mask)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestCalculateSubnetInvalidInput(t *testing.T) {
	_, err := CalculateSubnet("not-a-cidr")
	if err == nil {
		t.Error("expected error for invalid input, got nil")
	}
}
