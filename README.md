# BitSplitter

BitSplitter is a CLI network subnet calculator for IPv4 and IPv6, written in Go.

## Installation

### Homebrew (macOS/Linux)

```
brew tap justinsautter/bitsplitter
brew install bitsplitter
```

### Install Script (macOS/Linux)

```
curl -sSL https://raw.githubusercontent.com/justinsautter/bitsplitter/main/install.sh | sh
```

### Install Script (Windows PowerShell)

```powershell
irm https://raw.githubusercontent.com/justinsautter/bitsplitter/main/install.ps1 | iex
```

### Build from Source

```
git clone https://github.com/justinsautter/bitsplitter.git
cd bitsplitter
make install
```

## Usage

Pass an IP address with CIDR notation. BitSplitter auto-detects IPv4 vs IPv6.

```
bitsplitter [flags] <IP/CIDR> [<IP/CIDR>...]
bitsplitter [flags] <IP> <subnet mask>
```

### Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |
| `--table` | Show CIDR reference table |
| `--version` | Show version |
| `--no-color` | Disable color output |

### IPv4

```
bitsplitter 192.168.1.0/24
```

Output:
```
# Overview (IP/CIDR)
IP address: 192.168.1.0/24

# Address Details
Address range: 192.168.1.0 - 192.168.1.255
Number of hosts: 256

Network address: 192.168.1.0
Broadcast address: 192.168.1.255

Usable range: 192.168.1.1 - 192.168.1.254
Usable hosts: 254

Reverse DNS: 1.168.192.in-addr.arpa
Hex (IP): 0xC0A80100
Hex (network): 0xC0A80100

# Mask Information
Subnet mask: 255.255.255.0
Wildcard mask: 0.0.0.255
Subnet mask (binary): 11111111.11111111.11111111.00000000
Wildcard mask (binary): 00000000.00000000.00000000.11111111

# Classification
IP type: Private
IP class: C

# Supernet
Parent CIDR: 192.168.0.0/23
Parent range: 192.168.0.0 - 192.168.1.255

# Subnet Split
  192.168.1.0/25
  192.168.1.128/25
```

### IPv6

```
bitsplitter 2001:db8::/32
```

Output:
```
# Overview (IP/CIDR)
IP address: 2001:db8::/32

# Address Details
Address range: 2001:db8:: - 2001:db8:ffff:ffff:ffff:ffff:ffff:ffff
Number of hosts: 79228162514264337593543950336

Network address: 2001:db8::
Last address: 2001:db8:ffff:ffff:ffff:ffff:ffff:ffff

Usable range: 2001:db8:: - 2001:db8:ffff:ffff:ffff:ffff:ffff:ffff
Usable hosts: 79228162514264337593543950336

Reverse DNS: 8.b.d.0.1.0.0.2.ip6.arpa
Hex (IP): 0x20010DB8000000000000000000000000
Hex (network): 0x20010DB8000000000000000000000000

# Mask Information
Subnet mask: ffff:ffff:0000:0000:0000:0000:0000:0000
Wildcard mask: 0000:0000:ffff:ffff:ffff:ffff:ffff:ffff
Subnet mask (binary): 1111111111111111:1111111111111111:...
Wildcard mask (binary): 0000000000000000:0000000000000000:...

# Classification
IP type: Public

# Supernet
Parent CIDR: 2001:db8::/31
Parent range: 2001:db8:: - 2001:db9:ffff:ffff:ffff:ffff:ffff:ffff

# Subnet Split
  2001:db8::/33
  2001:db8:8000::/33
```

### Mask Notation

```
bitsplitter 192.168.1.0 255.255.255.0
```

### Overlap Detection

Pass multiple CIDRs to check for overlaps:

```
bitsplitter 10.0.0.0/8 10.1.0.0/16 172.16.0.0/12
```

Output:
```
# Network 1: 10.0.0.0/8
  Range: 10.0.0.0 - 10.255.255.255
  Hosts: 16777216
  Usable: 16777214

# Network 2: 10.1.0.0/16
  Range: 10.1.0.0 - 10.1.255.255
  Hosts: 65536
  Usable: 65534

# Network 3: 172.16.0.0/12
  Range: 172.16.0.0 - 172.31.255.255
  Hosts: 1048576
  Usable: 1048574

# Overlap Analysis
  10.0.0.0/8 vs 10.1.0.0/16: 10.0.0.0/8 contains 10.1.0.0/16
  10.0.0.0/8 vs 172.16.0.0/12: no overlap
  10.1.0.0/16 vs 172.16.0.0/12: no overlap
```

### JSON Output

```
bitsplitter --json 192.168.1.0/24
```

### CIDR Reference Table

```
bitsplitter --table
```

## License

See [LICENSE](LICENSE) for details.
