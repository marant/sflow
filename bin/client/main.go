package main

import (
	"flag"
	"fmt"
	"github.com/miekg/pcap"
	"net"
	"os"
	"sflow"
)

// errors
var (
	errNoSuchDevice = fmt.Errorf("No such device")
)

var (
	iface  string
	server = "localhost:1234"
)

func usage() {
	flag.PrintDefaults()
}

func main() {
	flag.StringVar(&iface, "iface", "", "Interface to read data from")
	flag.Usage = usage

	flag.Parse()

	if len(os.Args) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	pkts := make(chan *sflow.Packet)

	conn, err := connectToServer(server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't connect to server")
		os.Exit(-1)
	}

	go startCapture(iface, pkts)

	for pkt := range pkts {
		buf, err := pkt.Marshal()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}

		fmt.Println(pkt)

		written := 0
		for written < len(buf) {
			nbytes, err := conn.Write(buf[written:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
				break
			}
			written += nbytes
		}
	}
}

func connectToServer(server string) (net.Conn, error) {
	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func startCapture(name string, pkts chan *sflow.Packet) {
	handler, err := pcap.OpenLive(name, 65535, true, 500)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No such device: %s\n", name)
		os.Exit(-1)
	}

	for pkt, r := handler.NextEx(); r >= 0; pkt, r = handler.NextEx() {
		if r == 0 {
			continue
		}

		sflowPkt, err := parseSflowPacket(pkt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}

		pkts <- sflowPkt
	}
}

func parseSflowPacket(pkt *pcap.Packet) (*sflow.Packet, error) {
	sflowPkt := new(sflow.Packet)
	pkt.Decode()

	for _, hdr := range pkt.Headers {
		if iphdr, ok := hdr.(*pcap.Iphdr); ok {
			// if both ips are nil, this is not an IP packet
			// and we can ignore this
			if iphdr.SrcIp == nil && iphdr.DestIp == nil {
				continue
			}

			sflowPkt.Src = iphdr.SrcIp
			sflowPkt.Dst = iphdr.DestIp
		}

		if tcphdr, ok := hdr.(*pcap.Tcphdr); ok {
			sflowPkt.SrcPort = tcphdr.SrcPort
			sflowPkt.DstPort = tcphdr.DestPort
			sflowPkt.Protocol = sflow.PROTOCOL_TCP
		} else if udphdr, ok := hdr.(*pcap.Udphdr); ok {
			sflowPkt.SrcPort = udphdr.SrcPort
			sflowPkt.DstPort = udphdr.DestPort
			sflowPkt.Protocol = sflow.PROTOCOL_UDP
		}
	}

	return sflowPkt, nil
}
