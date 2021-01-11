// msgtick
package message

import (
	. "def"
	"reflect"
	. "user"

	. "ngcod.com/core"
)

type MsgTick struct {
}

func registerNetMsgTick() {
	isSuccess := RegisterMsgFunc(CMD_TICK, createNetMsgTick)
	LogInfo("Registor message", CMD_TICK)
	if !isSuccess {
		LogError("Registor message", CMD_TICK, "failed")
	}
}

func createNetMsgTick(cmdData *Command) NetMsg {
	netMsg := &MsgTick{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

func (this *MsgTick) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_TICK), reflect.ValueOf(this))
}

func (this *MsgTick) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgTick) Process(p interface{}) {
	puser, ok := p.(*User)
	if !ok {
		return
	}
	puser.Sender.Send(this)
}
