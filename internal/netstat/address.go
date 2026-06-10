package netstat

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ErrWildcardAddress is returned by ExtractHost for the netstat placeholder
// address "*:*", which UDP sockets report when they have not connected to a
// specific remote endpoint.
var ErrWildcardAddress = errors.New("wildcard address")

// ExtractHost splits a netstat-style "host:port" address into its host and
// port components. It accepts plain IPv4 addresses ("203.0.113.10:7777")
// and bracketed IPv6 addresses ("[2001:db8::1]:7777"). The UDP placeholder
// "*:*" returns ErrWildcardAddress.
func ExtractHost(addr string) (host string, port int, err error) {
	if addr == "*:*" {
		return "", 0, ErrWildcardAddress
	}

	if strings.HasPrefix(addr, "[") {
		end := strings.Index(addr, "]")
		if end == -1 {
			return "", 0, fmt.Errorf("invalid address %q: missing closing bracket", addr)
		}
		host = addr[1:end]
		rest := addr[end+1:]
		if !strings.HasPrefix(rest, ":") {
			return "", 0, fmt.Errorf("invalid address %q: missing port", addr)
		}
		port, err = strconv.Atoi(rest[1:])
		if err != nil {
			return "", 0, fmt.Errorf("invalid address %q: %w", addr, err)
		}
		return host, port, nil
	}

	idx := strings.LastIndex(addr, ":")
	if idx == -1 {
		return "", 0, fmt.Errorf("invalid address %q: missing port separator", addr)
	}

	host = addr[:idx]
	port, err = strconv.Atoi(addr[idx+1:])
	if err != nil {
		return "", 0, fmt.Errorf("invalid address %q: %w", addr, err)
	}

	return host, port, nil
}

// nonPublicIPv4Blocks lists the IPv4 ranges treated as non-public: loopback,
// private (RFC 1918), link-local, multicast, and other reserved space.
var nonPublicIPv4Blocks = []*net.IPNet{
	{IP: net.IPv4(0, 0, 0, 0).To4(), Mask: net.CIDRMask(8, 32)},      // "this" network
	{IP: net.IPv4(10, 0, 0, 0).To4(), Mask: net.CIDRMask(8, 32)},     // RFC 1918
	{IP: net.IPv4(127, 0, 0, 0).To4(), Mask: net.CIDRMask(8, 32)},    // loopback
	{IP: net.IPv4(169, 254, 0, 0).To4(), Mask: net.CIDRMask(16, 32)}, // link-local
	{IP: net.IPv4(172, 16, 0, 0).To4(), Mask: net.CIDRMask(12, 32)},  // RFC 1918
	{IP: net.IPv4(192, 168, 0, 0).To4(), Mask: net.CIDRMask(16, 32)}, // RFC 1918
	{IP: net.IPv4(224, 0, 0, 0).To4(), Mask: net.CIDRMask(4, 32)},    // multicast
	{IP: net.IPv4(240, 0, 0, 0).To4(), Mask: net.CIDRMask(4, 32)},    // reserved
}

// IsPublicIPv4 reports whether host is a routable, public IPv4 address. It
// returns false for non-IPv4 addresses (including IPv6) and for any address
// in a private, loopback, link-local, multicast, or reserved range.
func IsPublicIPv4(host string) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	for _, block := range nonPublicIPv4Blocks {
		if block.Contains(ip4) {
			return false
		}
	}

	return true
}
