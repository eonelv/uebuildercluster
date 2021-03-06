// msg_netreport
package message

import (
	. "def"
	"mydb"
	"reflect"
	"user"

	. "ngcod.com/core"
)

func registerNetMsgNetReport() {
	isSuccess := RegisterMsgFunc(CMD_NET_REPORT, createNetMsgNetReport)
	LogInfo("Registor message", CMD_NET_REPORT)
	if !isSuccess {
		LogError("Registor CMD_NET_REPORT faild")
	}
}

func createNetMsgNetReport(cmdData *Command) NetMsg {
	netMsg := &MsgNetReport{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

type MsgNetReport struct {
	ID      ObjectID
	Message [1024]byte
}

func (this *MsgNetReport) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_NET_REPORT), reflect.ValueOf(this))
}

func (this *MsgNetReport) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgNetReport) Process(p interface{}) {
	_, err := mydb.DBMgr.PreExecute("update t_project set buildstep=? where id=?", Byte2String(this.Message[:]), this.ID)
	if err != nil {
		LogError("Update serverState Error:", err)
		return
	}
	user.UserMgr.BroadcastMessage(this)
}
