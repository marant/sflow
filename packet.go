package sflow

import (
	"fmt"
	"net"
)

var (
	ErrNotIP         = fmt.Errorf("Not an IP packet")
	ErrNotEnoughData = fmt.Errorf("Not data to parse a sflow packet")
)

const (
	SFLOW_PACKET_LENGTH = 13
	PROTOCOL_TCP        = 0x01
	PROTOCOL_UDP        = 0x02
)

const (
	posSrc     = 0
	posDst     = 4
	posProto   = 8
	posSrcPort = 9
	posDstPort = 11
)

type Packet struct {
	Src      net.IP
	Dst      net.IP
	Protocol byte
	SrcPort  uint16
	DstPort  uint16
}

func (p *Packet) Marshal() ([]byte, error) {
	buf := make([]byte, 13)

	if p.Src == nil || p.Dst == nil {
		return nil, ErrNotIP
	}

	buf[posSrc], buf[posSrc+1], buf[posSrc+2], buf[posSrc+3] = p.Src[3], p.Src[2], p.Src[1], p.Src[0]
	buf[posDst], buf[posDst+1], buf[posDst+2], buf[posDst+3] = p.Dst[3], p.Dst[2], p.Dst[1], p.Dst[0]
	buf[posProto] = p.Protocol
	buf[posSrcPort], buf[posSrcPort+1] = byte(p.SrcPort>>8), byte(p.SrcPort)
	buf[posDstPort], buf[posDstPort+1] = byte(p.DstPort>>8), byte(p.DstPort)

	return buf, nil

}

func Unmarshal(data []byte) (*Packet, error) {
	if len(data) < SFLOW_PACKET_LENGTH {
		return nil, ErrNotEnoughData
	}

	pkt := new(Packet)
	pkt.Src = make([]byte, 4)
	pkt.Dst = make([]byte, 4)

	pkt.Src[0], pkt.Src[1], pkt.Src[2], pkt.Src[3] = data[posSrc+3], data[posSrc+2], data[posSrc+1], data[posSrc]
	pkt.Dst[0], pkt.Dst[1], pkt.Dst[2], pkt.Dst[3] = data[posDst+3], data[posDst+2], data[posDst+1], data[posDst]
	pkt.Protocol = data[posProto]
	pkt.SrcPort = uint16(data[posSrcPort])<<8 | uint16(data[posSrcPort+1])
	pkt.DstPort = uint16(data[posDstPort])<<8 | uint16(data[posDstPort+1])

	return pkt, nil
}
