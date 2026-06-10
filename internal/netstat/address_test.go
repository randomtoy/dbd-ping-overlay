package netstat

import (
	"errors"
	"testing"
)

func TestExtractHost(t *testing.T) {
	tests := []struct {
		name         string
		addr         string
		wantHost     string
		wantPort     int
		wantWildcard bool
		wantErr      bool
	}{
		{
			name:     "ipv4",
			addr:     "203.0.113.10:7777",
			wantHost: "203.0.113.10",
			wantPort: 7777,
		},
		{
			name:     "ipv4 zero address",
			addr:     "0.0.0.0:0",
			wantHost: "0.0.0.0",
			wantPort: 0,
		},
		{
			name:     "ipv6 bracketed",
			addr:     "[2001:db8::1]:443",
			wantHost: "2001:db8::1",
			wantPort: 443,
		},
		{
			name:         "wildcard",
			addr:         "*:*",
			wantWildcard: true,
		},
		{
			name:    "missing port",
			addr:    "203.0.113.10",
			wantErr: true,
		},
		{
			name:    "ipv6 missing closing bracket",
			addr:    "[2001:db8::1:443",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port, err := ExtractHost(tt.addr)

			switch {
			case tt.wantWildcard:
				if !errors.Is(err, ErrWildcardAddress) {
					t.Fatalf("ExtractHost(%q) error = %v, want ErrWildcardAddress", tt.addr, err)
				}
			case tt.wantErr:
				if err == nil {
					t.Fatalf("ExtractHost(%q) error = nil, want error", tt.addr)
				}
			default:
				if err != nil {
					t.Fatalf("ExtractHost(%q) returned unexpected error: %v", tt.addr, err)
				}
				if host != tt.wantHost {
					t.Errorf("ExtractHost(%q) host = %q, want %q", tt.addr, host, tt.wantHost)
				}
				if port != tt.wantPort {
					t.Errorf("ExtractHost(%q) port = %d, want %d", tt.addr, port, tt.wantPort)
				}
			}
		})
	}
}

func TestIsPublicIPv4(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"203.0.113.10", true},
		{"8.8.8.8", true},
		{"1.2.3.4", true},

		{"127.0.0.1", false},
		{"10.0.0.5", false},
		{"10.255.255.255", false},
		{"172.16.0.1", false},
		{"172.31.255.255", false},
		{"192.168.1.5", false},
		{"169.254.1.1", false},
		{"0.0.0.0", false},
		{"224.0.0.1", false},       // multicast
		{"255.255.255.255", false}, // reserved/broadcast

		{"not-an-ip", false},
		{"2001:db8::1", false}, // IPv6, out of scope for IsPublicIPv4
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			if got := IsPublicIPv4(tt.ip); got != tt.want {
				t.Errorf("IsPublicIPv4(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}
