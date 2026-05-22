// Package netguard gates outbound network calls so DevForge cannot be used
// as an SSRF accelerant when invoked by AI agents. Resolves a hostname,
// classifies each returned IP, and refuses to dial when any IP is private
// unless the caller explicitly opts in.
//
// Private = loopback, RFC1918 (10/8, 172.16/12, 192.168/16), CGNAT
// (100.64/10), link-local IPv4 (169.254/16), unique-local IPv6 (fc00::/7),
// link-local IPv6 (fe80::/10), unspecified, multicast.
package netguard

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Result of a safe DNS resolution.
type Result struct {
	Host       string
	IPs        []net.IP
	AnyPrivate bool
	Private    []net.IP // subset of IPs that flagged private
}

// IsPrivate reports whether ip falls into any private/loopback/link-local/CGNAT
// range that should not be dialled by an MCP-driven tool.
func IsPrivate(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsPrivate() {
		return true
	}
	// CGNAT 100.64.0.0/10 is not flagged by net.IP.IsPrivate. Check explicitly.
	if v4 := ip.To4(); v4 != nil {
		if v4[0] == 100 && (v4[1]&0xc0) == 64 {
			return true
		}
		// 0.0.0.0/8 is reserved.
		if v4[0] == 0 {
			return true
		}
	}
	return false
}

// ResolveSafe resolves host and classifies IPs. timeout caps the lookup.
// Always returns the IP slice even when AnyPrivate is true so callers can
// surface the offending addresses in diagnostics.
func ResolveSafe(ctx context.Context, host string, timeout time.Duration) (Result, error) {
	if host == "" {
		return Result{}, fmt.Errorf("netguard: host is empty")
	}
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	if timeout > 30*time.Second {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := net.DefaultResolver
	addrs, err := r.LookupIPAddr(ctx, host)
	if err != nil {
		return Result{Host: host}, err
	}
	out := Result{Host: host, IPs: make([]net.IP, 0, len(addrs))}
	for _, a := range addrs {
		out.IPs = append(out.IPs, a.IP)
		if IsPrivate(a.IP) {
			out.AnyPrivate = true
			out.Private = append(out.Private, a.IP)
		}
	}
	return out, nil
}
