package edgeCore

import (
	"../messageManager"
	"./edgeConnectionManager"
	// "net"
	"log"
	"fmt"
)

const (
	STATE_INIT = 0
	STATE_ATANDBY = 1
	STATE_SHUTTING_DOWN = 2
)

type EdgeCore struct {
	server_state int
	my_peer messageManager.Peer
	cm edgeConnectionManager.EdgeConnectionManager
	core_node_peer messageManager.Peer
	central_host string
}

func Init(my_port, core_node_host, core_node_port string) EdgeCore {
	log.Print("Initializing EdgeCore")
	edge_core := EdgeCore{}
	edge_core.server_state = STATE_INIT
	edge_core.my_peer = messageManager.Peer{Host: edge_core.get_ip(), Port: my_port}
	edge_core.core_node_peer = messageManager.Peer{Host: core_node_host, Port: core_node_port}
	edge_core.cm = edgeConnectionManager.EdgeConnectionManager_init(edge_core.my_peer, edge_core.core_node_peer)
	//central_hostを仮置き
	edge_core.central_host = core_node_host
	fmt.Printf("edge_core\n")
	fmt.Printf("(%%#v) %#v\n", edge_core)
	return edge_core
}

func (edge_core *EdgeCore)Start() {
	log.Print("Start EdgeCore")
	edge_core.server_state = STATE_ATANDBY
	edge_core.cm.Start()
}

func (edge_core *EdgeCore)shut_down() {
	edge_core.server_state = STATE_SHUTTING_DOWN
	edge_core.cm.Connection_close()
}

func (edge_core *EdgeCore)get_my_current_state() int {
	return edge_core.server_state
}

func (edge_core *EdgeCore)get_ip() string {
	// addrs, err := net.InterfaceAddrs()
    // if err != nil {
    //     fmt.Printf("Dial error: %s\n", err)
    //     return "error"
	// }
	// adress := ""
	// for _, addrs := range addrs {
	// 	adress = addrs.String()
	// }
	// fmt.Println(adress)
	// return adress
	return "127.0.0.1"
}