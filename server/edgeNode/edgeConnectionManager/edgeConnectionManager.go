package edgeConnectionManager

import (
	"log"
	"net"
	"bufio"
	"fmt"
	"../../messageManager"
	"../../coreNode/coreNodeList"
)

//todo
// message managerのimport
// handle messageの分岐内の中身
// mapSetをtoSliceしてfor rangeするのがよくない

const (
	PING_INTERVAL = 0
)

type EdgeConnectionManager struct {
	peer messageManager.Peer
	core_node_set CoreNodeList.CoreNodeList
	message_manager messageManager.MessageManager
	my_c_peer messageManager.Peer
}

//確認済み
func EdgeConnectionManager_init(my_peer, core_peer messageManager.Peer) EdgeConnectionManager {
	log.Print("Initializeing EdgeConnectionManager")
	edgeConnectionManager := EdgeConnectionManager{}
	edgeConnectionManager.peer = my_peer
	edgeConnectionManager.my_c_peer = core_peer
	edgeConnectionManager.core_node_set = CoreNodeList.CoreNodeList{List: []messageManager.Peer{}}
	//ホントはpeerの配列
	//peerをhostとmyportから作る
	edgeConnectionManager.__add_peer(core_peer)
	edgeConnectionManager.message_manager = messageManager.MessageManager{}
	fmt.Printf("edgeConnectionManager\n")
	fmt.Printf("(%%#v) %#v\n", edgeConnectionManager)
	return edgeConnectionManager
}

func (edgeConnectionManager *EdgeConnectionManager)Connect_to_core_node() {
	edgeConnectionManager.__connect_to_P2PNW()
}

func (edgeConnectionManager *EdgeConnectionManager)__connect_to_P2PNW() {
	msg := edgeConnectionManager.message_manager.Build(messageManager.MSG_ADD_AS_EDGE, false, edgeConnectionManager.peer)
	log.Print("message is here↓")
	fmt.Printf("(%%#v) %#v\n", msg)
	edgeConnectionManager.send_msg_to_core_node(msg)
}

func (edgeConnectionManager *EdgeConnectionManager)Start() {
	log.Print("call start")
	go edgeConnectionManager.__connect_to_P2PNW()
	edgeConnectionManager.__wait_for_access()
}

func (edgeConnectionManager *EdgeConnectionManager)__wait_for_access() {
	log.Print("start listening")
    ln := edgeConnectionManager.CreateListener()
    edgeConnectionManager.ListenWorker(ln)
}

// net.Listenを開き、listenerを返す
func (edgeConnectionManager *EdgeConnectionManager)CreateListener() net.Listener {
	//ちょい変更、要確認
    ln, err := net.Listen("tcp", edgeConnectionManager.peer.ToString())
    if err != nil {
        log.Fatal(err)
    }
    return ln
}

func (edgeConnectionManager *EdgeConnectionManager)ListenWorker(ln net.Listener) {
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
			edgeConnectionManager.__handle_message(params)
        }()
    }
}

func (edgeConnectionManager *EdgeConnectionManager)__handle_message(data_sum string) {
	var result string
	var payload bool
	var cmd int
	var peer messageManager.Peer
	log.Print("handling")
	result, _, cmd, peer, payload = edgeConnectionManager.message_manager.Parse(data_sum)
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
		log.Print("get request yay I'm now edge" + peer.ToString())
	} else if payload {
		log.Print("get request yay I'm now edge with payload")
	} else {
		log.Print("unexpected status")
	}
}

func (edgeConnectionManager *EdgeConnectionManager)__add_peer(peer messageManager.Peer) {
	log.Print("adding peer" + peer.ToString())
	//todo py_portをpeerに変更する
	edgeConnectionManager.core_node_set.Add(peer)
}

func (edgeConnectionManager *EdgeConnectionManager)__remove_peer(peer messageManager.Peer) {
	//todo py_portをpeerに変更する
	edgeConnectionManager.core_node_set.Remove(peer)
	log.Print("remove peer" + peer.ToString())
}

func (edgeConnectionManager *EdgeConnectionManager)send_msg(peer messageManager.Peer, msg string) bool {
	conn, err := net.Dial("tcp", peer.ToString())
    if err != nil {
        return false
    }
    defer conn.Close()
	conn.Write([]byte(msg))
	log.Print("success to send message")
	return true;
}

func (edgeConnectionManager *EdgeConnectionManager)send_msg_to_core_node(message string) {
	log.Print("send msg to core node")
	for _, peer := range edgeConnectionManager.core_node_set.List {
		if edgeConnectionManager.send_msg(peer, message) {
			break
		}
	}
}

func (edgeConnectionManager *EdgeConnectionManager)check_and_renew_my_core() {
	pre_check := send_ping(edgeConnectionManager.my_c_peer)
	for _, peer := range edgeConnectionManager.core_node_set.List {
		check := send_ping(peer)
		if !pre_check {
			if check {
				edgeConnectionManager.my_c_peer = peer
			}
		} else {
			break;
		}
		pre_check = check
	}
	for _, peer := range edgeConnectionManager.core_node_set.List {
		edgeConnectionManager.send_msg(peer, "check")
	}
}

func send_ping(peer messageManager.Peer) bool {
	conn, err := net.Dial("tcp", ":" + peer.Port)
    if err != nil {
        return false;
    }
	conn.Close()
	return true;
}

func (edgeConnectionManager *EdgeConnectionManager)Connection_close() {
	msg := edgeConnectionManager.message_manager.Build(messageManager.MSG_REMOVE, false, edgeConnectionManager.my_c_peer)
	edgeConnectionManager.send_msg_to_core_node(msg)
}