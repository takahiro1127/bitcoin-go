package EdgeNodeList

import (
	"log"
	"../../messageManager"
)

type EdgeNodeList struct{
	List []messageManager.Peer
}
//rangeで回す事によって、同時処理に弊害が起きる可能性が高い。

func (edge_node_set *EdgeNodeList)Add(peer messageManager.Peer) {
	log.Print("adding edge peer" + peer.ToString())
	if (edge_node_set.find(peer) == -1) {
		edge_node_set.List = append(edge_node_set.List, peer)
	} else {
		log.Print("already exist")
	}
}

func (edge_node_set *EdgeNodeList)Remove(peer messageManager.Peer) {
	log.Print("removing edge peer" + peer.ToString())
	index := edge_node_set.find(peer)
	if (index == -1) {
		log.Print("already removed")
	} else {
		edge_node_set.delete(index)
	}
}

func (edge_node_set *EdgeNodeList)find(peer messageManager.Peer) int {
	index := -1
    for i, v := range edge_node_set.List {
        if v == peer {
            index = i
        }
	}
	return index
}


func (edge_node_set *EdgeNodeList)delete(index int) {
    res := []messageManager.Peer{}
    for i, v := range edge_node_set.List {
        if i == index {
            continue
        }
        res = append(res, v)
    }
    edge_node_set.List = res
}

func (edge_node_set *EdgeNodeList)overwrite(new_list []messageManager.Peer) {
	log.Print("List will be overwrite")
	edge_node_set.List = new_list
}

func (edge_node_set *EdgeNodeList)GetList() []messageManager.Peer {
	return edge_node_set.List
}