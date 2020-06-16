package netcore

import (
	. "core"
	. "def"
	"mydb"
	"reflect"
	"time"
)

func (this *MsgBuild) Process(p interface{}) {
	puser, ok := p.(*User)
	if !ok {
		return
	}
	switch this.Action {
	case QUERY:
		this.query(puser)
	case BUILD:
		this.build(puser)
	}
}

func (this *MsgBuild) query(user *User) {
	rows, err := mydb.DBMgr.PreQuery("select id, name, projectName, host, member,serverState from t_project where id > 1000")
	if err != nil {
		LogError("query error. ", err)
		return
	}
	this.Num = byte(len(rows))
	var totalData []byte = []byte{}
	for _, v := range rows {
		project := &Project{}
		project.ID = v.GetObjectID("id")
		CopyArray(reflect.ValueOf(&project.Name), []byte(v.GetString("name")))
		CopyArray(reflect.ValueOf(&project.ProjectName), []byte(v.GetString("projectName")))
		CopyArray(reflect.ValueOf(&project.Host), []byte(v.GetString("host")))
		CopyArray(reflect.ValueOf(&project.Member), []byte(v.GetString("member")))
		project.ServerState = v.GetInt32("serverState")

		data, _ := Struct2Bytes(reflect.ValueOf(project))
		totalData = append(totalData, data...)
	}
	this.PData = totalData
	user.Sender.Send(this)
}

func (this *MsgBuild) build(user *User) {
	project := &Project{}
	Byte2Struct(reflect.ValueOf(project), this.PData)

	//0 未启动 2 使用中
	rows, err := mydb.DBMgr.PreQuery("select name, projectName, host, member, serverState, svn from t_project where id = ? and serverState = ?",
		project.ID, 1)
	if err != nil || len(rows) == 0 {
		LogError("no server stand by. ", err)
		return
	}
	this.Num = 1
	var totalData []byte = []byte{}
	for _, v := range rows {
		CopyArray(reflect.ValueOf(&project.Name), []byte(v.GetString("name")))
		CopyArray(reflect.ValueOf(&project.ProjectName), []byte(v.GetString("projectName")))
		CopyArray(reflect.ValueOf(&project.Host), []byte(v.GetString("host")))
		CopyArray(reflect.ValueOf(&project.Member), []byte(v.GetString("member")))
		project.ServerState = v.GetInt32("serverState")
		CopyArray(reflect.ValueOf(&project.SVN), []byte(v.GetString("svn")))

		LogDebug("准备发送数据到编译服务器", v.GetString("name"), v.GetString("member"), v.GetString("svn"))
		data, _ := Struct2Bytes(reflect.ValueOf(project))
		totalData = append(totalData, data...)
		break
	}

	this.PData = totalData

	ch := GetChanByID(project.ID)

	msgSend := &Command{CMD_SYSTEM_SERVER_BUILD, 0, nil, nil}
	msgSend.OtherInfo = this

	select {
	case ch <- msgSend:
	case <-time.After(5 * time.Second):
		LogError("Request server to build. put user channel failed:", CMD_SYSTEM_SERVER_BUILD)
	}
}

func (this *MsgBuildInfo) Process(p interface{}) {
	_, err := mydb.DBMgr.PreExecute("update t_project set serverState=? where id=?", this.ServerState, this.ID)
	LogDebug("Update server state to:", this.ServerState)
	if err != nil {
		LogError("Update serverState Error:", err)
		return
	}
}
