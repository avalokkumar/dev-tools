package netx

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/devforge/devforge/internal/netguard"
	"github.com/devforge/devforge/pkg/engine"
)

// DNSLookupOptions tunes DNSLookup.
type DNSLookupOptions struct {
	// Type is one of "a", "aaaa", "mx", "txt", "cname", "ns", "ptr", "all".
	// Default "all".
	Type string `json:"type,omitempty"`
	// TimeoutSec caps the lookup duration. Default 3, max 30.
	TimeoutSec int `json:"timeoutSec,omitempty"`
	// AllowPrivate permits resolved private/loopback IPs in the result. Default false.
	AllowPrivate bool `json:"allowPrivate,omitempty"`
}

// DNSLookupResult holds the records.
type DNSLookupResult struct {
	Host        string              `json:"host"`
	A           []string            `json:"a,omitempty"`
	AAAA        []string            `json:"aaaa,omitempty"`
	MX          []DNSMX             `json:"mx,omitempty"`
	TXT         []string            `json:"txt,omitempty"`
	CNAME       string              `json:"cname,omitempty"`
	NS          []string            `json:"ns,omitempty"`
	PTR         []string            `json:"ptr,omitempty"`
	AnyPrivate  bool                `json:"anyPrivate"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// DNSMX is one MX record entry.
type DNSMX struct {
	Host string `json:"host"`
	Pref int    `json:"pref"`
}

// DNSLookup resolves DNS records for host. Honors the netguard private-IP
// rule: when AllowPrivate=false and any A/AAAA result is private, the engine
// emits DNS.PRIVATE_BLOCKED at SevError and clears the IP fields.
func DNSLookup(host string, opts DNSLookupOptions) (DNSLookupResult, error) {
	host = strings.TrimSpace(host)
	if host == "" {
		return DNSLookupResult{Diagnostics: []engine.Diagnostic{{
			Code: "DNS.EMPTY", Message: "host is empty", Severity: engine.SevError,
		}}}, nil
	}
	timeout := time.Duration(opts.TimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	if timeout > 30*time.Second {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	r := net.DefaultResolver
	res := DNSLookupResult{Host: host}
	typ := strings.ToLower(opts.Type)
	if typ == "" {
		typ = "all"
	}

	if typ == "a" || typ == "aaaa" || typ == "all" {
		ips, err := r.LookupIPAddr(ctx, host)
		if err != nil {
			res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
				Code: "DNS.LOOKUP", Message: err.Error(), Severity: engine.SevError,
			})
		}
		for _, a := range ips {
			if netguard.IsPrivate(a.IP) {
				res.AnyPrivate = true
			}
			if a.IP.To4() != nil && (typ == "a" || typ == "all") {
				res.A = append(res.A, a.IP.String())
			} else if typ == "aaaa" || typ == "all" {
				res.AAAA = append(res.AAAA, a.IP.String())
			}
		}
		if res.AnyPrivate && !opts.AllowPrivate {
			res.A, res.AAAA = nil, nil
			res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
				Code: "DNS.PRIVATE_BLOCKED",
				Message: fmt.Sprintf("host %q resolves to a private IP; pass allowPrivate=true to view", host),
				Severity: engine.SevError,
			})
		}
	}
	if typ == "mx" || typ == "all" {
		mxs, err := r.LookupMX(ctx, host)
		if err == nil {
			for _, mx := range mxs {
				res.MX = append(res.MX, DNSMX{Host: strings.TrimSuffix(mx.Host, "."), Pref: int(mx.Pref)})
			}
			sort.Slice(res.MX, func(i, j int) bool { return res.MX[i].Pref < res.MX[j].Pref })
		}
	}
	if typ == "txt" || typ == "all" {
		txt, err := r.LookupTXT(ctx, host)
		if err == nil {
			res.TXT = txt
		}
	}
	if typ == "cname" || typ == "all" {
		cname, err := r.LookupCNAME(ctx, host)
		if err == nil {
			res.CNAME = strings.TrimSuffix(cname, ".")
		}
	}
	if typ == "ns" || typ == "all" {
		ns, err := r.LookupNS(ctx, host)
		if err == nil {
			for _, n := range ns {
				res.NS = append(res.NS, strings.TrimSuffix(n.Host, "."))
			}
			sort.Strings(res.NS)
		}
	}
	if typ == "ptr" {
		// PTR only meaningful for IPs; net.LookupAddr handles either.
		names, err := r.LookupAddr(ctx, host)
		if err == nil {
			for _, n := range names {
				res.PTR = append(res.PTR, strings.TrimSuffix(n, "."))
			}
		}
	}
	return res, nil
}
