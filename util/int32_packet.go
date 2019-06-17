package util

import (
	"encoding/binary"
	"github.com/FJSDS/cellnet"
	"github.com/FJSDS/cellnet/codec"
	"io"
)



const (
	intSize = 4
)

// 接收Length-Type-Value格式的封包流程
func RecvInt32LTVPacket(reader io.Reader, maxPacketSize int) (msg interface{},intValue int32, err error) {
	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, dataSize+intSize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(reader, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < dataSize+intSize {
		return nil,0, ErrMinPacket
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint16(sizeBuffer)
	msgID := binary.LittleEndian.Uint16(sizeBuffer[dataSize:])
	intValue =int32(binary.LittleEndian.Uint32(sizeBuffer[headSize:]))
	if size < dataSize+intSize{
		return nil,0, ErrMinPacket
	}
	if size >= uint16(maxPacketSize) {
		return nil,0, ErrMaxPacket
	}

	// 分配包体大小
	body := make([]byte, size-(dataSize+intSize))

	// 读取包体数据
	_, err = io.ReadFull(reader, body)

	// 发生错误时返回
	if err != nil {
		return
	}
	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(int(msgID), body)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil,0, err
	}

	return
}

// 发送Length-Type-Value格式的封包流程
func SendInt32LTVPacket(writer io.Writer, ctx cellnet.ContextSet, data interface{}) error {

	var (
		msgData []byte
		msgID   int
		meta    *cellnet.MessageMeta
		intValue int32
	)

	switch m := data.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	case *cellnet.IntPacket:
		intValue = m.IntValue
		var err error

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(data, ctx)

		if err != nil {
			return err
		}

		msgID = meta.ID
	default: // 发普通编码包
		var err error

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(data, ctx)

		if err != nil {
			return err
		}

		msgID = meta.ID
	}

	pkt := make([]byte, dataSize+msgIDSize+len(msgData))

	// Length
	binary.LittleEndian.PutUint16(pkt, uint16(headSize+len(msgData)))

	// Type
	binary.LittleEndian.PutUint16(pkt[dataSize:], uint16(msgID))
	//int32Value
	binary.LittleEndian.PutUint32(pkt[headSize:], uint32(intValue))
	// Value
	copy(pkt[headSize+intSize:], msgData)

	// 将数据写入Socket
	err := WriteFull(writer, pkt)

	// Codec中使用内存池时的释放位置
	if meta != nil {
		codec.FreeCodecResource(meta.Codec, msgData, ctx)
	}

	return err
}
