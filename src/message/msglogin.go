// msglogin
package message

import (
	. "def"
	"reflect"

	. "ngcod.com/core"
)

type MsgLogin struct {
	Account  NAME_STRING "MAC地址"
	UserID   ObjectID    "MsgConnection返回的ID"
	IsServer bool        "是编译服务器还是用户"
}

func registerNetMsgLogin() {
	isSuccess := RegisterMsgFunc(CMD_LOGIN, createNetMsgLogin)
	LogInfo("Registor message", CMD_LOGIN)
	if !isSuccess {
		LogError("Registor CMD_LOGIN faild")
	}
}

func createNetMsgLogin(cmdData *Command) NetMsg {
	netMsg := &MsgLogin{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

func (this *MsgLogin) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_LOGIN), reflect.ValueOf(this))
}

func (this *MsgLogin) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgLogin) Process(p interface{}) {

}
