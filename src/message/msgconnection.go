package message

import (
	. "core"
	. "def"
	"reflect"
)

type MsgConnection struct {
	ID        ObjectID
	AccountID NAME_STRING
}

func registerNetMsgConnection() {
	isSuccess := RegisterMsgFunc(CMD_CONNECTION, createNetMsgConnection)
	LogInfo("Registor message", CMD_CONNECTION)
	if !isSuccess {
		LogError("Registor CMD_BUILD faild")
	}
}

func createNetMsgConnection(cmdData *Command) NetMsg {
	netMsg := &MsgConnection{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

func (this *MsgConnection) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_CONNECTION), reflect.ValueOf(this))
}

func (this *MsgConnection) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgConnection) Process(p interface{}) {
}
