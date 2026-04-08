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

### IPv4

```
bitsplitter 172.16.27.0/18
```

Output:
```
# Overview (IP/CIDR)
IP address: 172.16.0.0/18

# Address Details
Address range: 172.16.0.0 - 172.16.63.255
Number of hosts: 16384

Network address: 172.16.0.0
Broadcast address: 172.16.63.255

Usable range: 172.16.0.1 - 172.16.63.254
Usable hosts: 16382

# Mask Information
Subnet mask (binary): 11111111.11111111.11000000.00000000
Wildcard mask (binary): 00000000.00000000.00111111.11111111

# Classification
IP type: Private
IP class: B
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

# Mask Information
Subnet mask (binary): 1111111111111111:1111111111111111:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000
Wildcard mask (binary): 0000000000000000:0000000000000000:1111111111111111:1111111111111111:1111111111111111:1111111111111111:1111111111111111:1111111111111111

# Classification
IP type: Public
```

## License

See [LICENSE](LICENSE) for details.
