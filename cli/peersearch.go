package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

//This file contains basics of Peer Discovery mechanism

const keyNeighbours = "neighbourPeers"
const bROKENREtries = 5
const idleConnectionTimeout = 5 //secs

type peerInfo struct {
	pid      uint64
	addr     net.Addr
	portList []portMapping
}
type portMapping struct {
	port    uint16
	appname string //Application associated with it
}

var peerMap map[uint64]net.Addr

type packetType byte

const (
	hELLO packetType = iota + 0
	hELLO_RCV
	hELLO_DEC

	lIST_PORTS
	pORTS
	pORTS_DEC

	cLOSE
	cLOSE_GRACEFULLY
)

//States for state machine
const (
	STATE_0 = iota + 0 //Initial State
	STATE_1
	STATE_2
	STATE_3 //Dead State
)

type discoveryPacket struct {
	PacketType packetType
	Pid        uint64

	Addr         net.Addr
	PortMappings []portMapping
	//	IsClient   bool //Specifies if the this message is from client or server
}

//paket size is 1+4 = 5 Bytes
func (p *discoveryPacket) marshal() []byte {

	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	encoder.Encode(*p)
	return w.Bytes()
}

func (p *discoveryPacket) unmarshal(b []byte) error {
	reader := bytes.NewReader(b)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(p)
}

func getBroadcastAddress() (broadcastIp []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Print(COLOR_RED)
		fmt.Println("Couldn't find network addresses")
		resetConsoleColor()
	}
	for _, v := range addrs {
		ipstr := v.String()
		_, ipnet, _ := net.ParseCIDR(ipstr)
		netip := ipnet.IP
		masked, size := ipnet.Mask.Size()
		patt := getbroadcastaddrpattern(size-masked, size)
		bnetip := []byte(netip)
		for i := 0; i < len(bnetip); i++ {
			bnetip[i] |= patt[i]
		}
		broadcastIp = append(broadcastIp, net.IP(bnetip).String())
	}
	return
}

func getbroadcastaddrpattern(unmaskedbits, addrSize int) []byte {
	var b []byte
	var x = int(unmaskedbits / 8)
	for i := 0; i < x; i++ {
		b = append(b, (1<<8)-1)
	}
	if unmaskedbits%8 != 0 {
		x := 1<<(unmaskedbits%8) - 1
		b = append(b, byte(x))
	}

	var ip = make([]byte, addrSize/8)
	for i := 0; i < len(b); i++ {
		ip[i] |= b[i]
	}
	//Let's reverse this pattern
	addrSize /= 8
	for i := 0; i < int(addrSize/2); i++ {
		var c byte
		c = ip[i]
		ip[i] = ip[addrSize-1-i]
		ip[addrSize-1-i] = c
	}
	return ip
}

func managePeerDiscoveryListener(discoveredPeer chan peerInfo) error {
	listener, err := net.Listen("udp", joinHost("", shreePort))

	if err != nil {
		//handling the error
		return err
	}
	//Putting the connection in socket collection
	socketCollection.add("peerdiscoverylistener", listener)
	go func() {
		defer listener.Close()
		for {
			inconn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleDiscoveryConnections(inconn, discoveredPeer, false)

		}
	}()

	return nil
}

//handleDiscoveryConnections, handles converstion b/w discovery server and client
//discoveredPeer, sends live discovered peer data to console
//isClient parameter tells whether this conversation is initiated by client or server
func handleDiscoveryConnections(inconn net.Conn, discoveredPeer chan peerInfo, isClient bool) {
	defer inconn.Close() //Closes the given connection
	var pi = peerInfo{}
	var breakFlag bool
	decoder := gob.NewDecoder(inconn)
	encoder := gob.NewEncoder(inconn)
	p := new(discoveryPacket)
	var cs uint = STATE_0 //Initial state at the of sessional state begin
	var currentPacketType packetType

	if isClient { //Being a client current connection will initiate conversation
		//Don't no why I m specific about it
		p.PacketType = hELLO
		p.Pid = currentPeer.pid
		encoder.Encode(p)

	}
	retries := 0
	for {
		err := decoder.Decode(p)
		if err != nil {
			nerr, ok := err.(net.Error)
			if ok {
				if nerr.Timeout() { //Checking for read timeout
					cs = STATE_3
				}
			}
			if retries > bROKENREtries {
				cs = STATE_3 //Setting current state to dead state
			} else {
				//Broken packet
				retries += 1
				continue

			}
		}
		retries = 0 //Reseting retries for valid packs
		currentPacketType = p.PacketType
		switch cs { //Current State
		case STATE_0:
			//connection sent hello/hello_rcv message
			//Along with this it will send its basic info i.e. peer id

			pi.pid = p.Pid

			switch p.PacketType {
			case hELLO: //Let's send response back
				p.PacketType = hELLO_RCV
				pi.addr = p.Addr                                                //Clients address
				p.Addr, _ = net.ResolveUDPAddr("udp4", joinHost("", shreePort)) //Server's address
				//Sending port information
			case hELLO_RCV:
				p.PacketType = lIST_PORTS
				//Let's retrived ports of the local host
				p.PortMappings = getMyPortMappings()

			}
		case STATE_1:
			switch p.PacketType {
			case lIST_PORTS:
				pi.portList = p.PortMappings
				//Everthing is perfect so far
				p.PacketType = pORTS
				//Let's retrived ports of the local host
				p.PortMappings = getMyPortMappings()
			case pORTS:
				p.PacketType = cLOSE_GRACEFULLY

			case pORTS_DEC:
				p.PacketType = cLOSE
			}

		case STATE_2:
			//Now this is the final state,
			//it means that tx has been complete and we can persist the port mapping and peer details
			switch p.PacketType {
			case cLOSE_GRACEFULLY:
				neigh, err := localdb.Get([]byte(keyNeighbours), nil)
				var peers []peerInfo
				if err == nil {
					json.Unmarshal(neigh, &peers)
					peers = append(peers, pi)
				}
				b, _ := json.Marshal(peers)
				localdb.Put([]byte(keyNeighbours), b, nil)
				p.PacketType = cLOSE_GRACEFULLY
				discoveredPeer <- pi
				breakFlag = true
			}
		case STATE_3:
			//At this point we will just send close state message to the opposite side
			p.PacketType = cLOSE

			breakFlag = true
		}
		if breakFlag {
			break
		}
		p.Pid = currentPeer.pid //Setting the pid of the packet

		cs = nextState(cs, currentPacketType, p.PacketType)
		encoder.Encode(p)
		inconn.SetReadDeadline(time.Now().Add(time.Second * idleConnectionTimeout))
	}
}

func getMyPortMappings() []portMapping {
	v, err := localdb.Get([]byte(keyLocalpfw), nil)
	var pm = make([]portMapping, 0)
	if err == nil {
		var ports [][2]uint16
		err = json.Unmarshal(v, &ports)
		if err == nil {
			//valid port marshling

			for _, v := range ports {
				pm = append(pm, portMapping{v[0], "Local Port Connection"})
			}
		}
	}

	return pm
}

//nextState, This functions implements nextState function for the discovery protocol fsm
func nextState(currentState uint, eventpacketType, responseEventpacketType packetType) (n uint) {
	n = currentState //This means if not a valid input supplied connection will be in its initial state
	switch currentState {
	case STATE_0:
		switch eventpacketType {
		case hELLO:
			switch responseEventpacketType {
			case hELLO_RCV:
				n = STATE_1
			case hELLO_DEC:
				n = STATE_3
			}
		case hELLO_RCV:
			switch responseEventpacketType {
			case lIST_PORTS:
				n = STATE_1

			}
		case hELLO_DEC:
			switch responseEventpacketType {
			case cLOSE:
				n = STATE_3
			}
		}

	case STATE_1:
		switch eventpacketType {
		case lIST_PORTS:
			switch responseEventpacketType {
			case pORTS:
				n = STATE_2
			case pORTS_DEC:
				n = STATE_3
			}
		case pORTS_DEC:
			switch responseEventpacketType {
			case cLOSE:
				n = STATE_3

			}
		case pORTS:
			switch responseEventpacketType {
			case cLOSE_GRACEFULLY:
				n = STATE_2
			}
		}

		//State 2 and 3 doesn't do any thing

	}
	return
}

func managePeerDiscoveryResponder() {
	for _, v := range getBroadcastAddress() {
		conn, err := net.Dial("udp", joinHost(v, shreePort))
		if err != nil {
			log.Print("Couldn't dial udp connection due to ", err.Error())
			continue
		}
		go handleDiscoveryConnections(conn, nil, true) //Handles client side of conversation

	}
}

//This code starts peer discovery listeners and manages responding
func startPeerDiscovery(currentpeer *peer) <-chan peerInfo {
	if peerMap == nil {
		peerMap = make(map[uint64]net.Addr)
	}
	var discoveredPeer = make(chan peerInfo, 10)

	err := managePeerDiscoveryListener(discoveredPeer) //Enables listening of manages, kind of server deployment
	if err != nil {
		print(COLOR_RED)
		println("Couldn't start peer discovery mechanism:\t", err.Error())
		resetConsoleColor()
	}
	go func() {
		//Setting timer to close peer discovery server
		defer socketCollection.close("peerdiscoverylistener")
		<-time.After(time.Minute)
	}()
	//Let's broadcast some packets
	managePeerDiscoveryResponder()
	return discoveredPeer
}

var print = fmt.Print
var println = fmt.Println
