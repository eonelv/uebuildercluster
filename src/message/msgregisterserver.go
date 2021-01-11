package message

import (
	. "def"
	. "idmgr"
	"mydb"
	"reflect"
	"strings"
	"time"

	. "ngcod.com/core"
)

type MsgServerRegister struct {
	Host        NAME_STRING "IP地址"
	Account     NAME_STRING "服务器所在的目录名字"
	UserID      ObjectID    "MsgConnection返回的ID"
	SVN         [1024]byte  "SVN地址"
	Member      [1024]byte  "通知的用户列表"
	ProjectName NAME_STRING "项目名称"
	IsServer    bool        "是编译服务器还是用户"
}

func registerNetMsgRegisterServer() {
	isSuccess := RegisterMsgFunc(CMD_REGISTER_SERVER, createNetMsgRegisterServer)
	LogInfo("Registor message", CMD_REGISTER_SERVER)
	if !isSuccess {
		LogError("Registor CMD_BUILD faild")
	}
}

func createNetMsgRegisterServer(cmdData *Command) NetMsg {
	netMsg := &MsgServerRegister{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

func (this *MsgServerRegister) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_REGISTER_SERVER), reflect.ValueOf(this))
}

func (this *MsgServerRegister) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgServerRegister) Process(p interface{}) {
	retChan := p.(chan ObjectID)

	rowsNum, err := mydb.DBMgr.PreQuery("select id from t_project where host = ? and account = ?",
		Byte2String(this.Host[:]), Byte2String(this.Account[:]))
	if err != nil {
		LogError(err)
		this.UserID = 1
		return
	}

	if len(rowsNum) != 0 {
		//不更新字段， 只有第一次注册服务器时，使用本地的config/config.json配置
		//以后统一在管理服务器修改
		/*
			rowsResult, err1 := mydb.DBMgr.PreExecute("update t_project set projectName=?, member=? where host=?",
				Byte2String(this.ProjectName[:]), Byte2String(this.Member[:]), Byte2String(this.Account[:]))
			if err1 != nil {
				LogError(err1)
				this.UserID = 2
				return
			}
			if num, _ := rowsResult.RowsAffected(); num == 0 {
				LogError(err1)
				this.UserID = 2
				return
			}
		*/
		this.UserID = rowsNum[0].GetObjectID("id")
		select {
		case retChan <- this.UserID:
		case <-time.After(20 * time.Second):
			LogError("MsgUserRegister send error")
		}
		return
	}

	id := SysIDGenerator.GetNextID(ID_SERVER)

	if id == 0 {
		LogError("generate ID faild")
		this.UserID = 3
		return
	}
	sql := "insert into t_project (id, name, projectName, host, account, svn, member, serverState) values (?,?,?,?,?,?,?,?)"

	members := Byte2String(this.Member[:])
	members = strings.ReplaceAll(members, `,`, "-")
	rowsResult, err1 := mydb.DBMgr.PreExecute(sql,
		id, Byte2String(this.ProjectName[:]), Byte2String(this.ProjectName[:]),
		Byte2String(this.Host[:]), Byte2String(this.Account[:]),
		Byte2String(this.SVN[:]), members, 0)
	if err1 != nil {
		LogError(err1)
		this.UserID = 2
		return
	}
	if num, _ := rowsResult.RowsAffected(); num == 0 {
		LogError(err1)
		this.UserID = 2
		return
	}
	this.UserID = id
	select {
	case retChan <- id:
	case <-time.After(20 * time.Second):
		LogError("MsgServerRegister send error")
	}
}
