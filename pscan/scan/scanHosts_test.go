package scan_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/acikgozb/cli-playground/pscan/scan"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}

	if ps.Open.String() != "closed" {
		t.Errorf("expected %q as port state but got %q", "closed", ps.Open.String())
	}

	ps.Open = true

	if ps.Open.String() != "open" {
		t.Errorf("expected %q as port state but got %q", "open", ps.Open.String())
	}
}

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name          string
		expectedState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}

	host := "localhost"
	hl := scan.HostsList{}

	hl.Add(host)

	ports := []int{}
	// Init ports, 1 open, 1 closed
	for _, tc := range testCases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}

		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)

		if tc.name == "ClosedPort" {
			// Instead of deferring the close call, we immediately close it for 2nd test case.
			ln.Close()
		}
	}

	// When
	res := scan.Run(&hl, ports)

	// Then
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %q instead\n", len(res))
	}

	if res[0].Host != host {
		t.Fatalf("expected host %q, got %q instead\n", host, res[0].Host)
	}

	if res[0].NotFound {
		t.Fatalf("expected host %q to be found\n", host)
	}

	if len(res[0].PortStates) != 2 {
		t.Fatalf("expected 2 port states, got %d instead\n", len(res[0].PortStates))
	}

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("expected port %d, got %d instead\n", ports[i], res[0].PortStates[i].Port)
		}

		if res[0].PortStates[i].Open.String() != tc.expectedState {
			t.Errorf("expected port %d to be %s\n", ports[i], tc.expectedState)
		}
	}
}

func TestRunHostNotFound(t *testing.T) {
	// This host should fail to be found unless you have it on your DNS:
	host := "389.389.389.389"

	hl := scan.HostsList{}
	hl.Add(host)

	res := scan.Run(&hl, []int{})

	// Verify the output - one Results, not found, empty PortState
	if len(res) != 1 {
		t.Fatalf("expected 1 results, got %d instead", len(res))
	}

	if res[0].Host != host {
		t.Errorf("expected host to be %s, got %s instead\n", host, res[0].Host)
	}

	if !res[0].NotFound {
		t.Errorf("expected host %q not to be found", host)
	}

	if len(res[0].PortStates) != 0 {
		t.Errorf("expected no portStates, but got %d portStates", len(res[0].PortStates))
	}
}
