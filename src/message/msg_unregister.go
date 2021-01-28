// msg_unregister.go
package message

import (
	. "def"
	"mydb"
	"reflect"
	"user"

	. "ngcod.com/core"
)

type MsgUnRegister struct {
	ID     ObjectID
	UserID ObjectID
	State  byte
}

func registerNetMsgUnRegister() {
	isSuccess := RegisterMsgFunc(CMD_UNREGISTER_SERSVER, createNetMsgUnRegister)
	LogInfo("Registor message", CMD_UNREGISTER_SERSVER)
	if !isSuccess {
		LogError("Registor CMD_UNREGISTER_SERSVER faild")
	}
}

func createNetMsgUnRegister(cmdData *Command) NetMsg {
	netMsg := &MsgUnRegister{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

func (this *MsgUnRegister) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_UNREGISTER_SERSVER), reflect.ValueOf(this))
}

func (this *MsgUnRegister) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgUnRegister) Process(p interface{}) {
	//转发给编译服务器
	if this.State == 1 {
		rows, err := mydb.DBMgr.PreQuery("select * from t_project where id=?", this.ID)
		if err != nil {
			LogError("remove database Error:", err, this.ID)
			return
		}
		if len(rows) == 0 {
			LogError("Project not exist")
			return
		}
		var serverState int32
		for _, v := range rows {
			serverState = v.GetInt32("serverState")
			break
		}
		//正在编译
		if serverState == 2 || serverState == 100 {
			LogError("Project is building now, can't remove")
			return
		}
		//正在删除
		if serverState == 1000 {
			LogError("Project is deleting now, don't submit again")
			return
		}

		_, err = mydb.DBMgr.PreExecute("update t_project set serverState = 1000 where id=?", this.ID)
		if err != nil {
			LogError("update serverState = 1000 Error:", err, this.ID)
			return
		}

		ch := GetChanByID(this.ID)
		msgSend := &Command{CMD_UNREGISTER_SERSVER, 0, nil, nil}
		msgSend.OtherInfo = this
		err = SendCommand(ch, msgSend, 10)

		if err != nil {
			LogError(err)
		}

	} else if this.State == 2 { //删除数据库，转发3给编译服务器，发送消息给用户
		_, err := mydb.DBMgr.PreExecute("delete from t_project where id=?", this.ID)
		if err != nil {
			LogError("remove database Error:", err, this.ID)
			return
		}
		this.State = 3
		user.UserMgr.BroadcastMessage(this)
	}

}
