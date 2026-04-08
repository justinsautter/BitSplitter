package subnet

import (
	"fmt"
	"net"
)

// OverlapResult describes the relationship between two networks
type OverlapResult struct {
	A        string `json:"a"`
	B        string `json:"b"`
	Overlaps bool   `json:"overlaps"`
	Relation string `json:"relation"` // "disjoint", "a contains b", "b contains a", "identical"
}

// CheckOverlap determines the relationship between two CIDR networks
func CheckOverlap(aCIDR, bCIDR string) (OverlapResult, error) {
	_, aNet, err := net.ParseCIDR(aCIDR)
	if err != nil {
		return OverlapResult{}, fmt.Errorf("invalid CIDR %q: %v", aCIDR, err)
	}
	_, bNet, err := net.ParseCIDR(bCIDR)
	if err != nil {
		return OverlapResult{}, fmt.Errorf("invalid CIDR %q: %v", bCIDR, err)
	}

	aContainsB := aNet.Contains(bNet.IP)
	bContainsA := bNet.Contains(aNet.IP)

	result := OverlapResult{
		A: aCIDR,
		B: bCIDR,
	}

	switch {
	case aContainsB && bContainsA:
		aMask, _ := aNet.Mask.Size()
		bMask, _ := bNet.Mask.Size()
		if aMask == bMask {
			result.Overlaps = true
			result.Relation = "identical"
		} else if aMask < bMask {
			result.Overlaps = true
			result.Relation = "a contains b"
		} else {
			result.Overlaps = true
			result.Relation = "b contains a"
		}
	case aContainsB:
		result.Overlaps = true
		result.Relation = "a contains b"
	case bContainsA:
		result.Overlaps = true
		result.Relation = "b contains a"
	default:
		result.Overlaps = false
		result.Relation = "disjoint"
	}

	return result, nil
}

// CheckAllOverlaps performs pairwise overlap checking on multiple CIDRs
func CheckAllOverlaps(cidrs []string) ([]OverlapResult, error) {
	var results []OverlapResult
	for i := 0; i < len(cidrs); i++ {
		for j := i + 1; j < len(cidrs); j++ {
			result, err := CheckOverlap(cidrs[i], cidrs[j])
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}
	return results, nil
}
