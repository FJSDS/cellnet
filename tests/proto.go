package tests

import (
	"fmt"
	"github.com/FJSDS/cellnet"
	"github.com/FJSDS/cellnet/codec"
	_ "github.com/FJSDS/cellnet/codec/binary"
	"github.com/FJSDS/cellnet/util"
	"reflect"
)

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    int(util.StringHash("tests.TestEchoACK")),
	})
}
