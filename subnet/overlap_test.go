package subnet

import "testing"

func TestCheckOverlap(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		overlaps bool
		relation string
	}{
		{
			name:     "disjoint",
			a:        "10.0.0.0/8",
			b:        "172.16.0.0/12",
			overlaps: false,
			relation: "disjoint",
		},
		{
			name:     "a contains b",
			a:        "10.0.0.0/8",
			b:        "10.1.0.0/16",
			overlaps: true,
			relation: "a contains b",
		},
		{
			name:     "b contains a",
			a:        "10.1.0.0/16",
			b:        "10.0.0.0/8",
			overlaps: true,
			relation: "b contains a",
		},
		{
			name:     "identical",
			a:        "192.168.1.0/24",
			b:        "192.168.1.0/24",
			overlaps: true,
			relation: "identical",
		},
		{
			name:     "adjacent not overlapping",
			a:        "192.168.0.0/24",
			b:        "192.168.1.0/24",
			overlaps: false,
			relation: "disjoint",
		},
		{
			name:     "IPv6 disjoint",
			a:        "2001:db8::/32",
			b:        "2001:db9::/32",
			overlaps: false,
			relation: "disjoint",
		},
		{
			name:     "IPv6 containment",
			a:        "2001:db8::/32",
			b:        "2001:db8:abcd::/48",
			overlaps: true,
			relation: "a contains b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CheckOverlap(tt.a, tt.b)
			if err != nil {
				t.Fatalf("CheckOverlap error: %v", err)
			}
			if result.Overlaps != tt.overlaps {
				t.Errorf("Overlaps = %v, want %v", result.Overlaps, tt.overlaps)
			}
			if result.Relation != tt.relation {
				t.Errorf("Relation = %s, want %s", result.Relation, tt.relation)
			}
		})
	}
}

func TestCheckAllOverlaps(t *testing.T) {
	results, err := CheckAllOverlaps([]string{"10.0.0.0/8", "10.1.0.0/16", "172.16.0.0/12"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 pairwise results, got %d", len(results))
	}
	// 10/8 vs 10.1/16 should overlap
	if !results[0].Overlaps {
		t.Error("expected 10/8 and 10.1/16 to overlap")
	}
	// 10/8 vs 172.16/12 should not
	if results[1].Overlaps {
		t.Error("expected 10/8 and 172.16/12 to not overlap")
	}
	// 10.1/16 vs 172.16/12 should not
	if results[2].Overlaps {
		t.Error("expected 10.1/16 and 172.16/12 to not overlap")
	}
}
