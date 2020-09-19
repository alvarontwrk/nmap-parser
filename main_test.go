package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetPorts(t *testing.T) {
	cases := []struct {
		inFile string
		ports  []int
	}{
		{"testdata/top500.xml", []int{21, 22, 139, 445}},
		{"testdata/allports.xml", []int{21, 22, 139, 445, 3632}},
		{"testdata/local.xml", []int{25, 111, 139, 445, 2049, 35111, 40495, 49577, 57135}},
	}

	for _, c := range cases {
		in, err := os.Open(c.inFile)
		if err != nil {
			t.Fatalf("failed to open input file: %s", err)
		}

		fileBytes := readerToBytes(in)
		ports := getPorts(fileBytes)
		if !reflect.DeepEqual(ports, c.ports) {
			t.Logf("want: %v", c.ports)
			t.Logf("have: %v", ports)
			t.Errorf("failed getting ports from %s", c.inFile)
		}
	}
}

func TestPortsToNmap(t *testing.T) {
	cases := []struct {
		inFile string
		out    string
	}{
		{"testdata/top500.xml", "-p21,22,139,445"},
		{"testdata/allports.xml", "-p21,22,139,445,3632"},
		{"testdata/local.xml", "-p25,111,139,445,2049,35111,40495,49577,57135"},
	}

	for _, c := range cases {
		in, err := os.Open(c.inFile)
		if err != nil {
			t.Fatalf("failed to open input file: %s", err)
		}

		fileBytes := readerToBytes(in)
		ports := getPorts(fileBytes)
		portsString := portsToNmap(ports)
		if portsString != c.out {
			t.Logf("want: %v", c.out)
			t.Logf("have: %v", portsString)
			t.Errorf("failed getting ports' string for nmap from %s", c.inFile)
		}
	}
}

func TestGetAddress(t *testing.T) {
	cases := []struct {
		inFile  string
		address string
	}{
		{"testdata/top500.xml", "10.10.10.3"},
		{"testdata/allports.xml", "10.10.10.3"},
		{"testdata/local.xml", "127.0.0.1"},
	}

	for _, c := range cases {
		in, err := os.Open(c.inFile)
		if err != nil {
			t.Fatalf("failed to open input file: %s", err)
		}

		fileBytes := readerToBytes(in)
		address := getAddress(fileBytes)
		if address != c.address {
			t.Logf("want: %v", c.address)
			t.Logf("have: %v", address)
			t.Errorf("failed getting address from %s", c.inFile)
		}
	}
}

func TestGetHostnames(t *testing.T) {
	cases := []struct {
		inFile    string
		hostnames []string
	}{
		{"testdata/top500.xml", []string{"lame.htb"}},
		{"testdata/allports.xml", []string{"lame.htb"}},
		{"testdata/local.xml", []string{"localhost"}},
	}

	for _, c := range cases {
		in, err := os.Open(c.inFile)
		if err != nil {
			t.Fatalf("failed to open input file: %s", err)
		}

		fileBytes := readerToBytes(in)
		hostnames := getHostnames(fileBytes)
		if !reflect.DeepEqual(hostnames, c.hostnames) {
			t.Logf("want: %v", c.hostnames)
			t.Logf("have: %v", hostnames)
			t.Errorf("failed getting hostnames from %s", c.inFile)
		}
	}
}
