package shree

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestInterfaceList(t *testing.T) {
	i, _ := net.Interfaces()
	for _, v := range i {
		addrs, _ := v.Addrs()
		for _, vv := range addrs {
			x := strings.Split(vv.String(), "/")
			ip := net.ParseIP(x[0])
			m := ip.DefaultMask()
			fmt.Print(m)
			/*var s string = ""
			if m != nil {
				s = m.String()

			}*/
			t.Log( /*v.Name, v.HardwareAddr,*/ vv.String() /*, ip, ip.IsLinkLocalUnicast(), s*/)

		}
	}
}

func TestInterfaceAddrList(t *testing.T) {
	addrs, _ := net.InterfaceAddrs()
	for _, v := range addrs {
		ipstr := strings.Split(v.String(), "/")[0]
		ip := net.ParseIP(ipstr)
		if ip == nil {
			continue
		}
		t.Log(ip.String(), ip.IsLoopback())
	}
}

// func TestInterfaceAddrNetList(t *testing.T) {
// 	addrs, _ := net.InterfaceAddrs()
// 	for _, v := range addrs {
// 		_, ipnet, _ := net.ParseCIDR(v.String())
// 		o, b := ipnet.Mask.Size()

// 		netaddr := ipnet.IP
// 		buf := new(bytes.Buffer)
// 		err := binary.Write(buf, binary.BigEndian, (1<<o)-1)
// 		if err != nil {
// 			t.Log(err)
// 		}
// 		maskedportion := buf.Bytes()

// 		buf = new(bytes.Buffer)
// 		binary.Write(buf, binary.BigEndian, (1<<b)-1)
// 		fullportion := buf.Bytes()

// 		t.Log(1<<b, 1<<o, b, o)
// 		t.Log(fullportion, maskedportion)
// 		diff := b/8 - len(maskedportion)/8
// 		for i := 0; i < diff; i++ {
// 			maskedportion = append(maskedportion, 0)

// 		}
// 		t.Log(fullportion, maskedportion)
// 		// bnetaddr := []byte(netaddr)
// 		// for i := 0; i < b/8; i++ {
// 		// 	fullportion[i] ^= maskedportion[i]
// 		// 	fullportion[i] |= bnetaddr[i]
// 		// }
// 		broadcastIP := net.IP(fullportion)
// 		t.Log(ipnet.IP.String(), broadcastIP.String(), o, b, netaddr)
// 		seq := getbroadcastaddrpattern(b-o, b)
// 		sseq := []byte(netaddr)
// 		t.Log(len(seq), len(sseq))
// 		for i := 0; i < len(seq); i++ {
// 			sseq[i] |= seq[i]
// 		}
// 		t.Log(net.IP(sseq).String())
// 	}
// }
// func TestGetBroadcastAddress(t *testing.T) {
// 	s := getBroadcastAddress()
// 	for _, v := range s {
// 		t.Log(v)
// 	}
// }
