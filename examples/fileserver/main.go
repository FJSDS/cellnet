package main

import (
	"flag"
	"fmt"
	"github.com/FJSDS/cellnet"
	"github.com/FJSDS/cellnet/peer"
	_ "github.com/FJSDS/cellnet/peer/http"
	"github.com/FJSDS/cellnet/proc"
	_ "github.com/FJSDS/cellnet/proc/http"
)

var shareDir = flag.String("share", ".", "folder to share")
var port = flag.Int("port", 9091, "listen port")

func main() {

	flag.Parse()

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("http.Acceptor", "httpfile", fmt.Sprintf(":%d", *port), nil).(cellnet.HTTPAcceptor)
	p.SetFileServe(".", *shareDir)

	proc.BindProcessorHandler(p, "http", nil)

	p.Start()
	queue.StartLoop()

	queue.Wait()
}
