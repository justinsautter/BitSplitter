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
			if info.IPClass != "N/A" {
				t.Errorf("IPClass = %s, want N/A", info.IPClass)
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

func TestCalculateSubnetInvalidInput(t *testing.T) {
	_, err := CalculateSubnet("not-a-cidr")
	if err == nil {
		t.Error("expected error for invalid input, got nil")
	}
}
