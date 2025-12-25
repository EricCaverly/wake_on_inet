package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type MacAddress [6]byte

type MagicPacket struct {
	Header  [6]byte
	Payload [16]MacAddress
}

func build_magic_packet(mac_addr MacAddress) (MagicPacket, error) {
	var mp MagicPacket

	// Build the header
	mp.Header = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	// Multiply MAC address 16 times into payload
	for i := range mp.Payload {
		mp.Payload[i] = mac_addr
	}

	return mp, nil
}

func convert_magic_to_bytes(mp MagicPacket) ([]byte, error) {
	// Make a buffer of bytes
	var buf bytes.Buffer

	// Write the contents of the struct into the buffer with BigEndian
	err := binary.Write(&buf, binary.BigEndian, mp)
	if err != nil {
		return nil, err
	}

	// Return back the buffer as bytes
	return buf.Bytes(), nil
}

func wake_pc(mac_addr string, broadcast_ip string) error {

	// Parse MAC
	hw_addr, err := net.ParseMAC(mac_addr)
	if err != nil {
		return err
	}

	// Form the destination IP:PORT
	dst_ip, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", broadcast_ip, 9))
	if err != nil {
		return err
	}

	// Make struct
	mp, err := build_magic_packet(MacAddress(hw_addr))
	if err != nil {
		return err
	}

	// Convert struct to bytes
	data, err := convert_magic_to_bytes(mp)
	if err != nil {
		return err
	}

	// Dial BROADCAST:9 with UDP
	conn, err := net.DialUDP("udp", nil, dst_ip)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Write the buffer of bytes containing FF...MAC...
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	if n != 102 {
		return fmt.Errorf("more data was sent than expected")
	}

	return nil
}
