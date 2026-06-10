package netstat

import "testing"

func TestDefaultSelectorPrefersEstablishedTCP(t *testing.T) {
	conns := []Connection{
		{Proto: "TCP", RemoteAddr: "203.0.113.10:7777", State: "ESTABLISHED", PID: 1},
		{Proto: "UDP", RemoteAddr: "198.51.100.20:9000", PID: 1},
	}

	got, ok := (DefaultSelector{}).Select(conns)
	if !ok {
		t.Fatal("Select() ok = false, want true")
	}
	if got.IP != "203.0.113.10" {
		t.Errorf("Select() IP = %q, want %q", got.IP, "203.0.113.10")
	}
}

func TestDefaultSelectorIgnoresPrivateAndWildcardAddresses(t *testing.T) {
	conns := []Connection{
		{Proto: "TCP", RemoteAddr: "192.168.1.1:443", State: "ESTABLISHED", PID: 1},
		{Proto: "UDP", RemoteAddr: "*:*", PID: 1},
		{Proto: "TCP", RemoteAddr: "203.0.113.10:7777", State: "ESTABLISHED", PID: 1},
	}

	got, ok := (DefaultSelector{}).Select(conns)
	if !ok {
		t.Fatal("Select() ok = false, want true")
	}
	if got.IP != "203.0.113.10" {
		t.Errorf("Select() IP = %q, want %q", got.IP, "203.0.113.10")
	}
}

func TestDefaultSelectorIgnoresNonEstablishedTCP(t *testing.T) {
	conns := []Connection{
		{Proto: "TCP", RemoteAddr: "203.0.113.10:7777", State: "TIME_WAIT", PID: 1},
		{Proto: "TCP", RemoteAddr: "203.0.113.10:7777", State: "LISTENING", PID: 1},
	}

	if _, ok := (DefaultSelector{}).Select(conns); ok {
		t.Fatal("Select() ok = true, want false")
	}
}

func TestDefaultSelectorFrequencyBreaksTies(t *testing.T) {
	conns := []Connection{
		{Proto: "UDP", RemoteAddr: "198.51.100.20:9000", PID: 1},
		{Proto: "UDP", RemoteAddr: "198.51.100.30:9001", PID: 1},
		{Proto: "UDP", RemoteAddr: "198.51.100.30:9001", PID: 1},
	}

	got, ok := (DefaultSelector{}).Select(conns)
	if !ok {
		t.Fatal("Select() ok = false, want true")
	}
	if got.IP != "198.51.100.30" {
		t.Errorf("Select() IP = %q, want %q", got.IP, "198.51.100.30")
	}
}

func TestDefaultSelectorEmpty(t *testing.T) {
	if _, ok := (DefaultSelector{}).Select(nil); ok {
		t.Fatal("Select(nil) ok = true, want false")
	}
}
