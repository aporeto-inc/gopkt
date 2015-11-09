/*
 * Network packet analysis framework.
 *
 * Copyright (c) 2014, Alessandro Ghedini
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
 * IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
 * THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
 * PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
 * CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
 * EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
 * PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package filter_test

import "log"
import "testing"

import "github.com/ghedo/go.pkt/filter"
import "github.com/ghedo/go.pkt/packet"

var test_eth_arp = []byte{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x4c, 0x72, 0xb9, 0x54, 0xe5, 0x3d,
	0x08, 0x06, 0x00, 0x01, 0x08, 0x00, 0x06, 0x04, 0x00, 0x01, 0x4c, 0x72,
	0xb9, 0x54, 0xe5, 0x3d, 0xc0, 0xa8, 0x01, 0x87, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xc1, 0x1b, 0xd0, 0x25,
}

var test_eth_vlan_arp = []byte{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x4c, 0x72, 0xb9, 0x54, 0xe5, 0x3d,
	0x81, 0x00, 0x00, 0x87, 0x08, 0x06, 0x00, 0x01, 0x08, 0x00, 0x06, 0x04,
	0x00, 0x01, 0x4c, 0x72, 0xb9, 0x54, 0xe5, 0x3d, 0xc0, 0xa8, 0x01, 0x87,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc1, 0x1b, 0xd0, 0x25,
}

var test_eth_ipv4_udp = []byte{
	0x00, 0x21, 0x96, 0x6e, 0xf0, 0x70, 0x4c, 0x72, 0xb9, 0x54, 0xe5, 0x3d,
	0x08, 0x00, 0x45, 0x00, 0x00, 0x1c, 0x00, 0x01, 0x00, 0x00, 0x40, 0x11,
	0x27, 0x60, 0xc0, 0xa8, 0x01, 0x87, 0xc1, 0x1b, 0xd0, 0x25, 0xa2, 0x5a,
	0x20, 0x92, 0x00, 0x08, 0xe9, 0x80,
}

var test_eth_ipv4_tcp = []byte{
	0x00, 0x21, 0x96, 0x6e, 0xf0, 0x70, 0x4c, 0x72, 0xb9, 0x54, 0xe5, 0x3d,
	0x08, 0x00, 0x45, 0x00, 0x00, 0x28, 0x00, 0x01, 0x00, 0x00, 0x40, 0x06,
	0x27, 0x5f, 0xc0, 0xa8, 0x01, 0x87, 0xc1, 0x1b, 0xd0, 0x25, 0xa2, 0x5a,
	0x20, 0x92, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50, 0x02,
	0x20, 0x00, 0x79, 0x85, 0x00, 0x00,
}

func TestMatch(t *testing.T) {
	arp, err := filter.Compile("arp", packet.Eth)
	if err != nil {
		t.Fatalf("Error compiling arp")
	}

	if !arp.Validate() {
		t.Fatalf("Invalid filter ARP\n%s", arp)
	}

	udp, err := filter.Compile("udp", packet.Eth)
	if err != nil {
		t.Fatalf("Error compiling udp")
	}

	port, err := filter.Compile("port 8338", packet.Eth)
	if err != nil {
		t.Fatalf("Error compiling port")
	}

	if !arp.Match(test_eth_arp) {
		t.Fatalf("ARP mismatch")
	}

	if arp.Match(test_eth_ipv4_udp) {
		t.Fatalf("ARP matched (but it shouldn't have)")
	}

	if !udp.Match(test_eth_ipv4_udp) {
		t.Fatalf("ARP mismatch")
	}

	if udp.Match(test_eth_ipv4_tcp) {
		t.Fatalf("UDP matched (but it shouldn't have)")
	}

	if !port.Match(test_eth_ipv4_udp) {
		t.Fatalf("UDP port mismatch")
	}

	if !port.Match(test_eth_ipv4_tcp) {
		t.Fatalf("TCP port mismatch")
	}

	if port.Match(test_eth_vlan_arp) {
		t.Fatalf("Port matched (but it shouldn't have)")
	}
}

func BenchmarkMatch(b *testing.B) {
	port, _ := filter.Compile("port 8338", packet.Eth)

	for n := 0; n < b.N; n++ {
		port.Match(test_eth_ipv4_tcp);
	}
}

func ExampleFilter() {
	// Match UDP or TCP packets on top of Ethernet
	flt, err := filter.Compile("udp or tcp", packet.Eth)
	if err != nil {
		log.Fatal(err)
	}

	if flt.Match([]byte("random data")) {
		log.Println("MATCH!!!")
	}
}
