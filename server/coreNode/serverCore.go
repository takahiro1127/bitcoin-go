package serverCore

import (
	"../messageManager"
	"./connectionManager"
	// "net"
	"log"
	"fmt"
)

const (
	STATE_INIT = 0
	STATE_ATANDBY = 1
	STATE_CONNECTED_TO_NETWORK = 2
	STATE_SHUTTING_DOWN = 3
)

type ServerCore struct {
	server_state int
	my_peer messageManager.Peer
	cm connectionManager.ConnectionManager
	core_node_peer messageManager.Peer
	central_host string
}

func Init(my_port, core_node_host, core_node_port string) ServerCore {
	log.Print("Initializing ServerCore")
	server_core := ServerCore{}
	server_core.server_state = STATE_INIT
	server_core.my_peer = messageManager.Peer{Host: server_core.get_ip(), Port: my_port}
	server_core.core_node_peer = messageManager.Peer{Host: core_node_host, Port: core_node_port}
	server_core.cm = connectionManager.ConnectionManager_init(server_core.my_peer, server_core.core_node_peer)
	//central_hostを仮置き
	server_core.central_host = core_node_host
	fmt.Printf("server_core\n")
	fmt.Printf("(%%#v) %#v\n", server_core)
	return server_core
}

func (server_core *ServerCore)Start() {
	log.Print("Start ServerCore")
	server_core.server_state = STATE_ATANDBY
	server_core.cm.Start()
}

func (server_core *ServerCore)Join_network() {
	if server_core.central_host != "" {
		log.Print("try to join in network")
		server_core.server_state = STATE_CONNECTED_TO_NETWORK
		server_core.cm.Join_network(server_core.my_peer)
	} else {
		log.Print("this server is runnning as genesis core node...")
	}
}

func (server_core *ServerCore)shut_down() {
	server_core.server_state = STATE_SHUTTING_DOWN
	server_core.cm.Connection_close()
}

func (server_core *ServerCore)get_my_current_state() int {
	return server_core.server_state
}

func (server_core *ServerCore)get_ip() string {
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