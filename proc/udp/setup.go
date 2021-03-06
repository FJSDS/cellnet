package udp

import (
	"github.com/FJSDS/cellnet"
	"github.com/FJSDS/cellnet/msglog"
	"github.com/FJSDS/cellnet/peer/udp"
	"github.com/FJSDS/cellnet/proc"
)

type UDPMessageTransmitter struct {
}

func (UDPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	data := ses.Raw().(udp.DataReader).ReadData()

	msg, err = recvPacket(data)

	msglog.WriteRecvLogger(log, "udp", ses, msg)

	return
}

func (UDPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	writer := ses.(udp.DataWriter)

	msglog.WriteSendLogger(log, "udp", ses, msg)

	// 由于UDP session会被重用，这里使用peer的ContextSet为内存池分配
	return sendPacket(writer, ses.Peer().(cellnet.ContextSet), msg)
}

func init() {

	proc.RegisterProcessor("udp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(UDPMessageTransmitter))
		bundle.SetCallback(userCallback)

	})
}
