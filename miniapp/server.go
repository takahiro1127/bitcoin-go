package miniapp

import (
	"net"
	"log"
	"bufio"
)

// net.Listenを開き、listenerを返す
//50030は適当だけど、変更したらclient.goのport番号も変更して上げる必要あり。
func CreateListener() net.Listener {
    ln, err := net.Listen("tcp", ":50030")
    if err != nil {
        log.Fatal(err)
    }
    return ln
}

func ListenWorker(ln net.Listener) {
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go func() {
            message, err := bufio.NewReader(conn).ReadString('\n')
            if err != nil {
                log.Fatal(err)
            }
            //現状メッセージ表示のみ
            log.Print("Message Received:", string(message))
            // newmessage := strings.ToUpper(message)
            conn.Write([]byte("Message Received:" + string(message) + "\n"))
            if message == "stop" {
                conn.Close()
            }
        }()
    }
}


// func main() {
//     log.Print("立ててる")
//     ln := CreateListener()
//     ListenWorker(ln)
// }