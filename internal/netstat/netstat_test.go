package netstat

import (
	"reflect"
	"testing"
)

func TestParseNetstatOutputTCP(t *testing.T) {
	output := "" +
		"\n" +
		"Active Connections\n" +
		"\n" +
		"  Proto  Local Address          Foreign Address        State           PID\n" +
		"  TCP    0.0.0.0:135            0.0.0.0:0              LISTENING       712\n" +
		"  TCP    192.168.1.5:54321      203.0.113.10:7777      ESTABLISHED     5678\n" +
		"  TCP    192.168.1.5:54322      203.0.113.10:7778      TIME_WAIT       5678\n"

	got, err := parseNetstatOutput(output)
	if err != nil {
		t.Fatalf("parseNetstatOutput() returned error: %v", err)
	}

	want := []Connection{
		{Proto: "TCP", LocalAddr: "0.0.0.0:135", RemoteAddr: "0.0.0.0:0", State: "LISTENING", PID: 712},
		{Proto: "TCP", LocalAddr: "192.168.1.5:54321", RemoteAddr: "203.0.113.10:7777", State: "ESTABLISHED", PID: 5678},
		{Proto: "TCP", LocalAddr: "192.168.1.5:54322", RemoteAddr: "203.0.113.10:7778", State: "TIME_WAIT", PID: 5678},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseNetstatOutput() = %#v, want %#v", got, want)
	}
}

func TestParseNetstatOutputUDP(t *testing.T) {
	output := "" +
		"  Proto  Local Address          Foreign Address        State           PID\n" +
		"  UDP    0.0.0.0:68             *:*                                    912\n" +
		"  UDP    192.168.1.5:7000       *:*                                    5678\n"

	got, err := parseNetstatOutput(output)
	if err != nil {
		t.Fatalf("parseNetstatOutput() returned error: %v", err)
	}

	want := []Connection{
		{Proto: "UDP", LocalAddr: "0.0.0.0:68", RemoteAddr: "*:*", State: "", PID: 912},
		{Proto: "UDP", LocalAddr: "192.168.1.5:7000", RemoteAddr: "*:*", State: "", PID: 5678},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseNetstatOutput() = %#v, want %#v", got, want)
	}
}

func TestParseNetstatOutputMixedAndIgnoresUnknownLines(t *testing.T) {
	output := "" +
		"Active Connections\n" +
		"\n" +
		"  Proto  Local Address          Foreign Address        State           PID\n" +
		"  TCP    192.168.1.5:54321      203.0.113.10:7777      ESTABLISHED     5678\n" +
		"  UDP    192.168.1.5:7000       *:*                                    5678\n" +
		"\n" +
		"  TCP    [::]:135               [::]:0                 LISTENING       712\n"

	got, err := parseNetstatOutput(output)
	if err != nil {
		t.Fatalf("parseNetstatOutput() returned error: %v", err)
	}

	want := []Connection{
		{Proto: "TCP", LocalAddr: "192.168.1.5:54321", RemoteAddr: "203.0.113.10:7777", State: "ESTABLISHED", PID: 5678},
		{Proto: "UDP", LocalAddr: "192.168.1.5:7000", RemoteAddr: "*:*", State: "", PID: 5678},
		{Proto: "TCP", LocalAddr: "[::]:135", RemoteAddr: "[::]:0", State: "LISTENING", PID: 712},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseNetstatOutput() = %#v, want %#v", got, want)
	}
}

func TestFilterByPID(t *testing.T) {
	conns := []Connection{
		{Proto: "TCP", PID: 100},
		{Proto: "UDP", PID: 200},
		{Proto: "TCP", PID: 100},
	}

	got := FilterByPID(conns, 100)
	want := []Connection{
		{Proto: "TCP", PID: 100},
		{Proto: "TCP", PID: 100},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FilterByPID() = %#v, want %#v", got, want)
	}

	if got := FilterByPID(conns, 999); got != nil {
		t.Errorf("FilterByPID() with no matches = %#v, want nil", got)
	}
}
