package CoreNodeList

import (
	"log"
	"../../messageManager"
)

type CoreNodeList struct{
	List []messageManager.Peer
}
//rangeで回す事によって、同時処理に弊害が起きる可能性が高い。

func (core_node_set *CoreNodeList)Add(peer messageManager.Peer) {
	log.Print("adding peer" + peer.ToString())
	if (core_node_set.find(peer) == -1) {
		core_node_set.List = append(core_node_set.List, peer)
	} else {
		log.Print("already exist")
	}
}

func (core_node_set *CoreNodeList)Remove(peer messageManager.Peer) {
	log.Print("removing peer" + peer.ToString())
	index := core_node_set.find(peer)
	if (index == -1) {
		log.Print("already removed")
	} else {
		core_node_set.delete(index)
	}
}

func (core_node_set *CoreNodeList)find(peer messageManager.Peer) int {
	index := -1
    for i, v := range core_node_set.List {
        if v == peer {
            index = i
        }
	}
	return index
}


func (core_node_set *CoreNodeList)delete(index int) {
    res := []messageManager.Peer{}
    for i, v := range core_node_set.List {
        if i == index {
            continue
        }
        res = append(res, v)
    }
    core_node_set.List = res
}

func (core_node_set *CoreNodeList)overwrite(new_list []messageManager.Peer) {
	log.Print("List will be overwrite")
	core_node_set.List = new_list
}

func (core_node_set *CoreNodeList)GetList() []messageManager.Peer {
	return core_node_set.List
}