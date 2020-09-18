package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type NmapRun struct {
	XMLName xml.Name `xml:"nmaprun"`
	Version string   `xml:"version,attr"`
	Hosts   []Host   `xml:"host"`
}

type Host struct {
	XMLName xml.Name `xml:"host"`
	Ports   Ports    `xml:"ports"`
	Address Address  `xml:"address"`
}

type Address struct {
	XMLName xml.Name `xml:"address"`
	Addr    string   `xml:"addr,attr"`
}

type Ports struct {
	XMLName   xml.Name   `xml:"ports"`
	OpenPorts []OpenPort `xml:"port"`
}

type OpenPort struct {
	XMLName  xml.Name `xml:"port"`
	Protocol string   `xml:"protocol,attr"`
	PortID   int      `xml:"portid,attr"`
	Service  Service  `xml:"service"`
}

type Service struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr"`
}

var nmaprun NmapRun

func getPorts(bytes []byte) []int {
	nmaprun = NmapRun{}
	/*
		buf := new(bytes.Buffer)
		burf.ReadFrom(r)
	*/
	xml.Unmarshal(bytes, &nmaprun)
	var openPorts []int
	for _, openPort := range nmaprun.Hosts[0].Ports.OpenPorts {
		openPorts = append(openPorts, openPort.PortID)
	}
	return openPorts
}

func intSliceToString(intElements []int, glue string) string {
	var stringElements []string
	for _, elem := range intElements {
		stringElements = append(stringElements, strconv.Itoa(elem))
	}
	return strings.Join(stringElements, glue)
}

func portsToNmap(ports []int) string {
	return "-p" + intSliceToString(ports, ",")
}

func getAddress(bytes []byte) string {
	nmaprun = NmapRun{}
	/*
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		print(len(buf.Bytes()))
	*/
	xml.Unmarshal(bytes, &nmaprun)
	return nmaprun.Hosts[0].Address.Addr
}

func readerToBytes(r io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.Bytes()
}

func init() {
	flag.Usage = func() {
		h := "Parse XML output from Nmap (file or stdin) to get reusable information\n\n"

		h += "Usage:\n"
		h += "  nmap-parser [OPTIONS] [FILE|URL|-]\n\n"

		h += "Options:\n"
		h += "  -a, --address  Get address\n"
		h += "  -p, --ports    List open ports\n"
		h += "\n"

		fmt.Fprintf(os.Stderr, h)
	}
}

func main() {
	var (
		portFlag    bool
		addressFlag bool
	)

	flag.BoolVar(&addressFlag, "address", false, "")
	flag.BoolVar(&addressFlag, "a", false, "")
	flag.BoolVar(&portFlag, "ports", false, "")
	flag.BoolVar(&portFlag, "p", false, "")

	flag.Parse()

	var rawInput io.Reader
	filename := flag.Arg(0)
	if filename == "" || filename == "-" {
		rawInput = os.Stdin
	} else {
		r, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(-1)
		}
		rawInput = r
	}

	fileBytes := readerToBytes(rawInput)

	if portFlag {
		ports := getPorts(fileBytes)
		portsAsString := portsToNmap(ports)
		fmt.Println(portsAsString)
		os.Exit(0)
	} else if addressFlag {
		address := getAddress(fileBytes)
		fmt.Println(address)
		os.Exit(0)
	} else {
		address := getAddress(fileBytes)
		ports := getPorts(fileBytes)
		portsAsString := intSliceToString(ports, ", ")
		fmt.Fprintf(os.Stdout, "[+] Host:\t%s\n", address)
		fmt.Fprintf(os.Stdout, "[+] Ports:\t%s\n", portsAsString)
	}
}
