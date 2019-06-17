package tcp

import (
	"github.com/FJSDS/cellnet"
	"github.com/FJSDS/cellnet/proc"
)

func init() {

	proc.RegisterProcessor("tcp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(TCPMessageTransmitter))
		bundle.SetHooker(new(MsgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))

	})

	proc.RegisterProcessor("tcp.int32ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(TCPInt32MessageTransmitter))
		bundle.SetHooker(new(MsgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))

	})
}
