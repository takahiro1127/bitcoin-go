package main

import (
	"./coreNode"
	"./edgeNode"
	"flag"
	"log"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if args[0] == "server" {
		//corenodeは奇数
		if args[1] == "2" {
			log.Print("main start as a second node")
			my_p2p_server := serverCore.Init("5446", "127.0.0.1", "5445")
			//後半は2ではない場合のserverを指している。
			my_p2p_server.Join_network()
			my_p2p_server.Start()
		} else if args[1] == "3" {
			log.Print("main start as a second node")
			my_p2p_server := serverCore.Init("9949", "127.0.0.1", "5445")
			//後半は2ではない場合のserverを指している。
			my_p2p_server.Join_network()
			my_p2p_server.Start()
		} else {
			my_p2p_server := serverCore.Init("5445", "127.0.0.1", "5445")
			//後半は自分自身を指している
			my_p2p_server.Start()
		}
	} else if args[0] == "client" {
		//edgenodeは偶数
		if args[1] == "2" {
			log.Print("main start as a second node")
			my_p2p_server := edgeCore.Init("4444", "127.0.0.1", "5445")
			my_p2p_server.Start()
		} else if args[1] == "3" {
			log.Print("main start as a second node")
			my_p2p_server := edgeCore.Init("6666", "127.0.0.1", "9943")
			my_p2p_server.Start()
		} else {
			my_p2p_server := edgeCore.Init("8888", "127.0.0.1", "5445")
			my_p2p_server.Start()
		}
	}
}