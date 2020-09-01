// msgfileinfo
package message

import (
	. "core"
	. "def"
	"reflect"
	"user"
)

func registerNetMsgFileInfo() {
	isSuccess := RegisterMsgFunc(CMD_FILE, createNetMsgNetFileInfo)
	LogInfo("Registor message", CMD_FILE)
	if !isSuccess {
		LogError("Registor CMD_FILE faild")
	}
}

func createNetMsgNetFileInfo(cmdData *Command) NetMsg {
	netMsg := &MsgFile{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

const (
	QUERY_FILE  uint16 = 1
	LIST_FILE   uint16 = 2
	REMOVE_FILE uint16 = 3
)

type MsgFileInfo struct {
	IsDir    bool
	Size     int64
	FileName [1024]byte
}

type MsgFile struct {
	ProjectID ObjectID
	Action    uint16
	Num       uint16
	PData     []byte //MsgFileInfo
}

func (this *MsgFile) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_FILE), reflect.ValueOf(this))
}

func (this *MsgFile) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgFile) Process(p interface{}) {
	if this.Action == LIST_FILE {
		user.UserMgr.BroadcastMessage(this)
	} else {

		fileInfo := &MsgFileInfo{}
		Byte2Struct(reflect.ValueOf(fileInfo), this.PData[:])
		ch := GetChanByID(this.ProjectID)

		msgSend := &Command{CMD_FILE, 0, nil, nil}
		msgSend.OtherInfo = this

		SendCommand(ch, msgSend, 10)
	}
}
