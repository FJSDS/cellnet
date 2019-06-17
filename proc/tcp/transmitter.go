package tcp

import (
	"github.com/FJSDS/cellnet"
	"github.com/FJSDS/cellnet/util"
	"io"
	"net"
)

type TCPMessageTransmitter struct {
}

type socketOpt interface {
	MaxPacketSize() int
	ApplySocketReadTimeout(conn net.Conn, callback func())
	ApplySocketWriteTimeout(conn net.Conn, callback func())
}

func (TCPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	opt := ses.Peer().(socketOpt)

	if conn, ok := reader.(net.Conn); ok {

		// 有读超时时，设置超时
		opt.ApplySocketReadTimeout(conn, func() {

			msg, err = util.RecvLTVPacket(reader, opt.MaxPacketSize())

		})
	}

	return
}

func (TCPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) (err error) {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	opt := ses.Peer().(socketOpt)

	// 有写超时时，设置超时
	opt.ApplySocketWriteTimeout(writer.(net.Conn), func() {

		err = util.SendLTVPacket(writer, ses.(cellnet.ContextSet), msg)

	})

	return
}

type TCPInt32MessageTransmitter struct {
}


func (TCPInt32MessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	opt := ses.Peer().(socketOpt)

	if conn, ok := reader.(net.Conn); ok {

		// 有读超时时，设置超时
		opt.ApplySocketReadTimeout(conn, func() {

			msgTemp, intValue,errTemp := util.RecvInt32LTVPacket(reader, opt.MaxPacketSize())
			msg= &cellnet.IntPacket{
				IntValue:intValue,
				Msg:msgTemp,
			}
			err = errTemp
		})
	}

	return
}

func (TCPInt32MessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) (err error) {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	opt := ses.Peer().(socketOpt)

	// 有写超时时，设置超时
	opt.ApplySocketWriteTimeout(writer.(net.Conn), func() {
		err = util.SendInt32LTVPacket(writer, ses.(cellnet.ContextSet), msg)
	})

	return
}
