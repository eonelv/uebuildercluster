// msgproject
package message

import (
	. "def"
	"mydb"
	"reflect"
	. "user"

	. "ngcod.com/core"
)

func registerNetMsgProject() {
	isSuccess := RegisterMsgFunc(CMD_PROJECT, createNetMsgProject)
	LogInfo("Registor message", CMD_PROJECT)
	if !isSuccess {
		LogError("Registor CMD_BUILD faild")
	}
}

func createNetMsgProject(cmdData *Command) NetMsg {
	netMsg := &MsgProject{}
	netMsg.CreateByBytes(cmdData.Message.([]byte))
	return netMsg
}

const (
	QUERY_PROJECT uint16 = 1
	SAVE_PROJECT  uint16 = 2
)

type MsgProject struct {
	Action uint16
	Num    uint16 "总的Project数量"
	PData  []byte "所有Project"
}

type Project struct {
	ID          ObjectID
	Name        [255]byte
	ProjectName [255]byte
	Host        [255]byte
	Account     [255]byte
	Member      [255]byte
	BuildStep   [1024]byte
	SVN         [255]byte
	Desc        [1024]byte
	ServerState int32
}

func (this *MsgProject) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_PROJECT), reflect.ValueOf(this))
}

func (this *MsgProject) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgProject) Process(p interface{}) {
	puser, ok := p.(*User)
	if !ok {
		return
	}

	switch this.Action {
	case QUERY_PROJECT:
		this.query(puser)
		break
	case SAVE_PROJECT:
		this.save(puser)
		break
	}
}

func (this *MsgProject) query(user *User) {
	sql := "select id, name, projectName, host, account, member, buildstep, svn, desc, serverState from t_project where id > 1000"
	rows, err := mydb.DBMgr.PreQuery(sql)
	if err != nil {
		LogError("query error. ", err)
		return
	}
	this.Num = uint16(len(rows))
	var totalData []byte = []byte{}
	for _, v := range rows {
		project := &Project{}
		project.ID = v.GetObjectID("id")
		CopyArray(reflect.ValueOf(&project.Name), []byte(v.GetString("name")))
		CopyArray(reflect.ValueOf(&project.ProjectName), []byte(v.GetString("projectName")))
		CopyArray(reflect.ValueOf(&project.Host), []byte(v.GetString("host")))
		CopyArray(reflect.ValueOf(&project.Account), []byte(v.GetString("account")))
		CopyArray(reflect.ValueOf(&project.Member), []byte(v.GetString("member")))
		CopyArray(reflect.ValueOf(&project.BuildStep), []byte(v.GetString("buildstep")))
		CopyArray(reflect.ValueOf(&project.SVN), []byte(v.GetString("svn")))
		CopyArray(reflect.ValueOf(&project.Desc), []byte(v.GetString("desc")))
		project.ServerState = v.GetInt32("serverState")

		data, _ := Struct2Bytes(reflect.ValueOf(project))
		totalData = append(totalData, data...)
	}
	this.PData = totalData
	user.Sender.Send(this)
}

func (this *MsgProject) save(user *User) {
	project := &Project{}
	Byte2Struct(reflect.ValueOf(project), this.PData)
	_, err := mydb.DBMgr.PreExecute("update t_project set name=?, member=?, desc=? where id=?",
		Byte2String(project.Name[:]), Byte2String(project.Member[:]), Byte2String(project.Desc[:]), project.ID)
	if err != nil {
		LogError("Save Project Error:", err)
		return
	}

	//user.Sender.Send(this)
}
