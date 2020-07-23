package netcore

import (
	. "core"
	. "def"
	. "message"
	"mydb"
	"reflect"
	"time"
	. "user"
)

const (
	ServerStateNone     int32 = 0
	ServerStateIdle     int32 = 1
	ServerStateBuilding int32 = 2
)

const (
	ServerErrorNone  int32 = 0
	ServerErrorBegin int32 = 2

	ServerErrorBuilding        int32 = 100
	ServerErrorOtherBuilding   int32 = 200
	ServerErrorNoServerUseable int32 = 300
	ServerErrorChanError       int32 = 400
	ServerErrorDB              int32 = 500
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
	rows, err := mydb.DBMgr.PreQuery("select id, name, projectName, host, member, serverState from t_project where id > 1000")
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
		CopyArray(reflect.ValueOf(&project.Member), []byte(v.GetString("member")))
		project.ServerState = v.GetInt32("serverState")

		data, _ := Struct2Bytes(reflect.ValueOf(project))
		totalData = append(totalData, data...)
	}
	this.PData = totalData
	user.Sender.Send(this)
}

func (this *MsgBuild) build(user *User) {
	LogDebug("开始编译")
	project := &Project{}
	Byte2Struct(reflect.ValueOf(project), this.PData)

	//0 未启动 2 使用中
	rows, err := mydb.DBMgr.PreQuery("select id from t_project where host = ? and serverState = ?",
		Byte2String(project.Host[:]), ServerStateBuilding)
	if err != nil || len(rows) != 0 {
		if len(rows) != 0 {
			if rows[0].GetObjectID("id") == project.ID {
				LogError("server is building. ", err)
				this.sendBack(user, project, ServerErrorBuilding)
			} else {
				LogError("another server is building. ", err)
				this.sendBack(user, project, ServerErrorOtherBuilding)
			}
		} else {
			LogError("another server is building. ", err)
			this.sendBack(user, project, ServerErrorDB)
		}
		return
	}

	//0 未启动 2 使用中
	rows, err = mydb.DBMgr.PreQuery("select name, projectName, host, account, member, serverState, svn from t_project where id = ? and serverState = ?",
		project.ID, ServerStateIdle)
	if err != nil || len(rows) == 0 {
		this.sendBack(user, project, ServerErrorNoServerUseable)
		LogError("no server stand by. ", err)
		return
	}
	this.Num = 1
	var totalData []byte = []byte{}
	for _, v := range rows {
		CopyArray(reflect.ValueOf(&project.Name), []byte(v.GetString("name")))
		CopyArray(reflect.ValueOf(&project.ProjectName), []byte(v.GetString("projectName")))
		CopyArray(reflect.ValueOf(&project.Host), []byte(v.GetString("host")))
		CopyArray(reflect.ValueOf(&project.Account), []byte(v.GetString("account")))
		CopyArray(reflect.ValueOf(&project.Member), []byte(v.GetString("member")))
		CopyArray(reflect.ValueOf(&project.SVN), []byte(v.GetString("svn")))
		project.ServerState = v.GetInt32("serverState")

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
		//this.sendBack(user, project, ServerErrorIdle)
		LogDebug("1. 发送消息到编译服务器(成功)")
	case <-time.After(5 * time.Second):
		this.sendBack(user, project, ServerErrorChanError)
		LogError("Request server to build. put user channel failed:", CMD_SYSTEM_SERVER_BUILD)
	}
}

func (this *MsgBuild) sendBack(u *User, project *Project, state int32) {
	msgBuildInfo := &MsgBuildInfo{}
	msgBuildInfo.ID = project.ID
	msgBuildInfo.UserID = this.UserID
	msgBuildInfo.Host = project.Host
	msgBuildInfo.Name = project.Name
	msgBuildInfo.ProjectName = project.ProjectName
	msgBuildInfo.ServerState = state

	u.Sender.Send(msgBuildInfo)
}

func (this *MsgBuildInfo) Process(p interface{}) {

	targetUserCh := GetChanByID(this.UserID)
	msgSend := &Command{CMD_SYSTEM_SERVER_BUILD, 0, nil, this}
	err := SendCommand(targetUserCh, msgSend, 10)
	if err != nil {
		LogError(err)
	}

	_, err = mydb.DBMgr.PreExecute("update t_project set serverState=?, buildstep=? where id=?", this.ServerState, "", this.ID)
	LogDebug("Update server state to:", this.ServerState)
	if err != nil {
		LogError("Update serverState Error:", err)
		return
	}
}
