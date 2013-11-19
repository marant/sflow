package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net"
	"net/http"
	"os"
	"sflow"
)

var (
	listenPort = "1234"
	listenAddr = "localhost"
)

type Server struct {
	pkts chan *sflow.Packet
}

func main() {
	server := Server{}
	server.pkts = make(chan *sflow.Packet)

	http.Handle("/data", websocket.Handler(server.wsHandler))
	http.Handle("/", http.FileServer(http.Dir("./www/")))

	go listenForPackets(server.pkts)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func listenForPackets(pkts chan *sflow.Packet) {
	conn, err := net.ListenPacket("udp", listenAddr+":"+listenPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}

	buf := make([]byte, sflow.SFLOW_PACKET_LENGTH)
	for {
		_, _, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			continue
		}

		pkt, err := sflow.Unmarshal(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			continue
		}

		pkts <- pkt
	}
}

func (s *Server) wsHandler(ws *websocket.Conn) {
loop:
	for pkt := range s.pkts {
		buf, err := pkt.Marshal()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			continue
		}

		err = websocket.Message.Send(ws, buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			break loop
		}
	}

	fmt.Println("Connection lost")
}
