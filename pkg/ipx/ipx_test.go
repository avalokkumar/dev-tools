package ipx

import (
	"testing"
)

func TestCIDR_IPv4_24(t *testing.T) {
	t.Parallel()
	r, _ := CIDRCalc("192.168.1.0/24", Options{})
	if r.Network != "192.168.1.0" {
		t.Fatalf("net = %s", r.Network)
	}
	if r.First != "192.168.1.0" || r.Last != "192.168.1.255" {
		t.Fatalf("first/last = %s..%s", r.First, r.Last)
	}
	if r.Netmask != "255.255.255.0" {
		t.Fatalf("mask = %s", r.Netmask)
	}
	if r.Wildcard != "0.0.0.255" {
		t.Fatalf("wildcard = %s", r.Wildcard)
	}
	if r.NumAddresses != "256" {
		t.Fatalf("num = %s", r.NumAddresses)
	}
	if r.UsableHosts != "254" {
		t.Fatalf("usable = %s", r.UsableHosts)
	}
}

func TestCIDR_IPv4_30(t *testing.T) {
	t.Parallel()
	r, _ := CIDRCalc("10.0.0.0/30", Options{MaxList: 8})
	if len(r.Hosts) != 4 {
		t.Fatalf("hosts = %v", r.Hosts)
	}
}

func TestCIDR_IPv6(t *testing.T) {
	t.Parallel()
	r, _ := CIDRCalc("2001:db8::/64", Options{})
	if r.Family != "ipv6" {
		t.Fatalf("family = %s", r.Family)
	}
	if r.Prefix != 64 || r.HostBits != 64 {
		t.Fatalf("prefix/host = %d/%d", r.Prefix, r.HostBits)
	}
	if r.NumAddresses == "" {
		t.Fatalf("empty num")
	}
}

func TestCIDR_Invalid(t *testing.T) {
	t.Parallel()
	r, _ := CIDRCalc("not-cidr", Options{})
	if r.Diagnostics[0].Code != "IP.PARSE" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}
