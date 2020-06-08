package messageManager

import (
	"encoding/json"
	"log"
	"fmt"
)

const (
	PROTOCOL_NAME = "simple_bitcoin_protocol"
	MY_VERSION = "0.1.0"
	MSG_ADD = 0
	MSG_REMOVE = 1
	MSG_CORE_LIST = 3
	MSG_REQUEST_CORE_LIST = 3
	MSG_PING = 4
	MSG_ADD_AS_EDGE = 5
	MSG_REMOVE_EDGE = 6

	ERR_PROTOCOL_UNMATCH = 0
	ERR_VERSION_UNMATCH = 1
	OK_WITH_PAYLOAD = 2
	OK_WITHOUT_PAYLOAD = 3

	JSON_UNMARSHAL_ERROR = 500
)

type MessageManager struct{
}
type Message struct {
	Protocol string `json:"protocol"`
	Version string `json:"version"`
	Msg_type int `json:"msg_type"`
	My_port string `json:"my_port"`
	My_host string `json:"my_host"`
	Payload bool `json:"payload"`
	Peer Peer
}

type Peer struct {
	Host string
	Port string
}

func (peer *Peer)ToString() string {
	return peer.Host + ":" + peer.Port
}

func (message_manager *MessageManager)Build(msg_type int, payload bool, my_peer Peer) string {
	//todo corelistに対応する
	message := &Message{Protocol: PROTOCOL_NAME, Version: MY_VERSION, Msg_type: msg_type, Payload: payload, My_host: my_peer.Host, My_port: my_peer.Port}
	jsonBytes, err := json.Marshal(message)
	if err != nil {
        return "Json Marshal Error"
	}
	return string(jsonBytes) + "\n"
}

func (message_manager *MessageManager)Parse(msg string) (string, int, int, Peer, bool) {
	log.Print("start Parse");
	jsonBytes := ([]byte)(msg)
	data := new(Message)
	fmt.Print(data);
	if err := json.Unmarshal(jsonBytes, data); err != nil {
        return "error", JSON_UNMARSHAL_ERROR, 500, Peer{Host: "jsonerror", Port: "jsonerror"}, false
	}
	peer := Peer{Host: data.My_host, Port: data.My_port}
	if data.Protocol != PROTOCOL_NAME {
		return "error", ERR_PROTOCOL_UNMATCH, 500, peer, false
	} else if data.Version > MY_VERSION {
		return "error", ERR_VERSION_UNMATCH, 500, peer, false
	} else if data.Msg_type == MSG_CORE_LIST {
		return "ok", OK_WITH_PAYLOAD, data.Msg_type, peer, true
	} else {
		return "ok", OK_WITHOUT_PAYLOAD, data.Msg_type, peer, false
	}
}

