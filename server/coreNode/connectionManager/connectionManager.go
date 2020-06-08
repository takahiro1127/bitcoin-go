package connectionManager

import (
	"log"
	"net"
	"bufio"
	"fmt"
	"../../messageManager"
	"../coreNodeList"
	"../../edgeNode/edgeNodeList"
)

//todo
// message managerのimport
// handle messageの分岐内の中身
// mapSetをtoSliceしてfor rangeするのがよくない

const (
	PING_INTERVAL = 0
)

type ConnectionManager struct {
	peer messageManager.Peer
	core_node_set CoreNodeList.CoreNodeList
	edge_node_set EdgeNodeList.EdgeNodeList
	message_manager messageManager.MessageManager
	my_c_peer messageManager.Peer
	call_back func
}

func ConnectionManager_init(my_peer, core_peer messageManager.Peer, call_cack func) ConnectionManager {
	log.Print("Initializeing ConnectionManager")
	connectionManager := ConnectionManager{}
	connectionManager.peer = my_peer
	connectionManager.core_node_set = CoreNodeList.CoreNodeList{List: []messageManager.Peer{}}
	connectionManager.edge_node_set = EdgeNodeList.EdgeNodeList{List: []messageManager.Peer{}}
	connectionManager.call_back = call_back
	//ホントはpeerの配列
	//peerをhostとmyportから作る
	if core_peer != my_peer {
		connectionManager.__add_peer(core_peer)
	}
	connectionManager.message_manager = messageManager.MessageManager{}
	fmt.Printf("connectionManager\n")
	fmt.Printf("(%%#v) %#v\n", connectionManager)
	return connectionManager
}

func (connectionManager *ConnectionManager)__add_peer(peer messageManager.Peer) {
	log.Print("adding peer " + peer.ToString())
	//todo py_portをpeerに変更する
	connectionManager.core_node_set.Add(peer)
	log.Print("added and now is below\n")
	fmt.Printf("(%%#v) %#v\n", connectionManager.core_node_set.GetList())
}

func (connectionManager *ConnectionManager)__remove_peer(peer messageManager.Peer) {
	//todo py_portをpeerに変更する
	connectionManager.core_node_set.Remove(peer)
	log.Print("remove peer " + peer.ToString())
}

func (connectionManager *ConnectionManager)__add_edge_peer(peer messageManager.Peer) {
	log.Print("adding edge peer " + peer.ToString())
	//todo py_portをpeerに変更する
	connectionManager.edge_node_set.Add(peer)
	log.Print("added and now is below\n")
	fmt.Printf("(%%#v) %#v\n", connectionManager.edge_node_set.GetList())
}

func (connectionManager *ConnectionManager)__remove_edge_peer(peer messageManager.Peer) {
	//todo py_portをpeerに変更する
	connectionManager.edge_node_set.Remove(peer)
	log.Print("remove edge peer " + peer.ToString())
}

// net.Listenを開き、listenerを返す
//50030は適当だけど、変更したらclient.goのport番号も変更して上げる必要あり。
func (connectionManager *ConnectionManager)CreateListener() net.Listener {
	//ちょい変更、要確認
    ln, err := net.Listen("tcp", connectionManager.peer.ToString())
    if err != nil {
        log.Fatal(err)
    }
    return ln
}

func (connectionManager *ConnectionManager)ListenWorker(ln net.Listener) {
    for {
		log.Print("waiting for connection")
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go func() {
			log.Print("connected by ....")
			params, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Print("got error");
				log.Print(params);
                log.Fatal(err)
			}
			log.Print("get params");
			log.Print(params);
			connectionManager.__handle_message(params)
        }()
    }
}

func (connectionManager *ConnectionManager)__wait_for_access() {
	log.Print("start listening")
    ln := connectionManager.CreateListener()
    connectionManager.ListenWorker(ln)
}

func (connectionManager *ConnectionManager)__handle_message(data_sum string) {
	var result string
	var payload bool
	var cmd int
	var peer messageManager.Peer
	log.Print("handling")
	result, _, cmd, peer, payload = connectionManager.message_manager.Parse(data_sum)
	log.Print("result")
	log.Print(result)
	log.Print("cmd")
	log.Print(cmd)
	log.Print("payload")
	log.Print(payload)
	if result == "error" {
		log.Print("errorです")
		return
	} else if !payload {
		if cmd == messageManager.MSG_ADD {
			log.Print("Add node request was received")
			connectionManager.__add_peer(peer)
			if peer != connectionManager.peer {
				//todo ここで、新しく追加されたpeerを他のpeerに追加するメッセージを作る。
				//初めて受け取ったpeerであるかどうかの情報が必要となる
				//そうでなければ無限ループが発生してしまう。
				//removeについても似た処理を行う必要がある。
				msg := connectionManager.message_manager.Build(messageManager.MSG_CORE_LIST, false, connectionManager.peer)
				connectionManager.send_msg_to_all_peer(msg)
			}
		} else if cmd == messageManager.MSG_REMOVE {
			log.Print("REMOVE request was recieved")
			connectionManager.__remove_peer(peer)
			msg := connectionManager.message_manager.Build(messageManager.MSG_CORE_LIST, false, connectionManager.peer)
			connectionManager.send_msg_to_all_peer(msg)
		} else if cmd == messageManager.MSG_REQUEST_CORE_LIST {
			log.Print("List for core nodes was recieved")
		} else if cmd == messageManager.MSG_ADD_AS_EDGE {
			log.Print("add Edge node")
			connectionManager.__add_edge_peer(peer)
			msg := connectionManager.message_manager.Build(messageManager.MSG_CORE_LIST, false, connectionManager.peer)
			connectionManager.send_msg(peer, msg)
		} else if cmd == messageManager.MSG_ADD_AS_EDGE {
			log.Print("remove edge node")
			connectionManager.__remove_edge_peer(peer)
		} else if cmd == messageManager.MSG_NEW_TRANSACTION {
			log.Print("NEW TRANSACTION")
		} else if cmd == messageManager.MSG_NEW_BLOCK {
			log.Print("NEW BLOCK")
		} else if cmd == messageManager.RSP_FULL_CHAIN {
			log.Print("RSP_FULL_CHAIN")
		} else if cmd == messageManager.MSG_ENHANCED {
			log.Print("MSG ENHANCED")
		} else {
			connectionManager.call_back(data_sum, peer)
		}
	} else if payload {
		if cmd == messageManager.MSG_CORE_LIST {
			log.Print("refresh the core node list...")
			// payloadでpeerの配列を受け取り、いれる
		} else {
			log.Print("unknown command")
			connectionManager.call_back(data_sum, peer)
		}
	} else {
		log.Print("unexpected status")
	}
}

func (connectionManager *ConnectionManager)send_msg(peer messageManager.Peer, msg string) {
	conn, err := net.Dial("tcp", peer.ToString())
    if err != nil {
        connectionManager.__remove_peer(peer)
        return
    }
    defer conn.Close()
    conn.Write([]byte(msg))
}

func (connectionManager *ConnectionManager)send_msg_to_all_peer(message string) {
	log.Print("send msg to all peer")
	for _, peer := range connectionManager.core_node_set.List {
		connectionManager.send_msg(peer, message)
	}
}

func (connectionManager *ConnectionManager)__check_peers_connection() {
	log.Print("check now")
	for _, peer := range connectionManager.core_node_set.List {
		if !__is_alive(peer) {
			connectionManager.core_node_set.Remove(peer)
		}
	}
	for _, peer := range connectionManager.core_node_set.List {
		connectionManager.send_msg(peer, "check")
	}
}

func __is_alive(peer messageManager.Peer) bool {
	conn, err := net.Dial("tcp", ":" + peer.Port)
    if err != nil {
        return false;
    }
	conn.Close()
	return true;
}

func (connectionManager *ConnectionManager)Start() {
	log.Print("call start")
	connectionManager.__wait_for_access()
}

func (connectionManager *ConnectionManager)Join_network(peer messageManager.Peer) {
	connectionManager.my_c_peer = peer
	log.Print("set my_c_peer as" + connectionManager.my_c_peer.ToString())
	connectionManager.__connect_to_P2PNW(peer)
}

func (connectionManager *ConnectionManager)__connect_to_P2PNW(peer messageManager.Peer) {
	msg := connectionManager.message_manager.Build(messageManager.MSG_ADD, false, peer)
	log.Print("message is here↓")
	fmt.Printf("(%%#v) %#v\n", msg)
	connectionManager.send_msg_to_all_peer(msg)
}

func (connectionManager *ConnectionManager)Connection_close() {
	msg := connectionManager.message_manager.Build(messageManager.MSG_REMOVE, false, connectionManager.my_c_peer)
	connectionManager.send_msg(connectionManager.my_c_peer, msg)
}