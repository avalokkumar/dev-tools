// Package ipx provides IPv4/IPv6 subnet math.
//
// External API:
//
//	CIDRCalc(cidr string, opts Options) (Result, error)
package ipx

import (
	"fmt"
	"math/big"
	"net"

	"github.com/devforge/devforge/pkg/engine"
)

// Options tunes CIDRCalc.
type Options struct {
	// MaxList caps how many host addresses we will enumerate. 0 = no list.
	MaxList int `json:"maxList,omitempty"`
}

// Result is the subnet breakdown.
type Result struct {
	Network      string              `json:"network"`
	Broadcast    string              `json:"broadcast,omitempty"` // IPv4 only
	First        string              `json:"first"`
	Last         string              `json:"last"`
	Netmask      string              `json:"netmask"`
	Wildcard     string              `json:"wildcard,omitempty"` // IPv4 only
	Prefix       int                 `json:"prefix"`
	HostBits     int                 `json:"hostBits"`
	NumAddresses string              `json:"numAddresses"`
	UsableHosts  string              `json:"usableHosts"`
	Family       string              `json:"family"` // "ipv4" | "ipv6"
	Hosts        []string            `json:"hosts,omitempty"`
	Diagnostics  []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// CIDRCalc parses a CIDR notation and computes derived values.
func CIDRCalc(cidr string, opts Options) (Result, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "IP.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	prefix, bits := ipnet.Mask.Size()
	hostBits := bits - prefix
	res := Result{Prefix: prefix, HostBits: hostBits}
	if ip.To4() != nil {
		res.Family = "ipv4"
	} else {
		res.Family = "ipv6"
	}

	num := new(big.Int).Lsh(big.NewInt(1), uint(hostBits))
	res.NumAddresses = num.String()
	usable := new(big.Int).Set(num)
	if res.Family == "ipv4" && hostBits >= 2 {
		usable.Sub(usable, big.NewInt(2))
	}
	res.UsableHosts = usable.String()

	// Network = ipnet.IP; mask in canonical form.
	res.Network = ipnet.IP.String()
	res.Netmask = ipMaskToString(ipnet.Mask, res.Family)
	if res.Family == "ipv4" {
		res.Wildcard = wildcardString(ipnet.Mask)
	}

	first, last := firstLast(ipnet)
	res.First = first.String()
	res.Last = last.String()
	if res.Family == "ipv4" {
		res.Broadcast = last.String()
	}

	if opts.MaxList > 0 {
		res.Hosts = enumerate(first, last, opts.MaxList)
	}
	return res, nil
}

func ipMaskToString(m net.IPMask, family string) string {
	if family == "ipv4" {
		return net.IP(m).To4().String()
	}
	// IPv6 mask formatted like an address.
	return net.IP(m).String()
}

func wildcardString(m net.IPMask) string {
	wc := make([]byte, len(m))
	for i, b := range m {
		wc[i] = ^b
	}
	return net.IP(wc).String()
}

func firstLast(n *net.IPNet) (net.IP, net.IP) {
	first := n.IP
	last := make(net.IP, len(n.IP))
	for i := range n.IP {
		last[i] = n.IP[i] | ^n.Mask[i]
	}
	return first, last
}

func enumerate(first, last net.IP, limit int) []string {
	out := []string{}
	cur := append(net.IP{}, first...)
	for i := 0; i < limit; i++ {
		out = append(out, cur.String())
		if cur.Equal(last) {
			break
		}
		incr(cur)
	}
	return out
}

func incr(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			return
		}
	}
	// overflow — leave as 0s (caller is bounded)
	_ = fmt.Sprint("noop")
}
