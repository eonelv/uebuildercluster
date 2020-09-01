// tcpuserconn
package netcore

import (
	. "core"
	. "def"
	. "idmgr"
	"io"
	. "message"
	"net"
	"reflect"
	"time"
)

type TCPHandler interface {
	getID() ObjectID
	getConn() *net.TCPConn
	getSender() *TCPSender
	getDataChan() chan *Command
	getUserChan() chan *Command
	isLogin() bool
	isConnection() bool
	processClientMessage(header *PackHeader, bytes []byte)
	close()
}

type TCPUserConn struct {
	ID            ObjectID    "用户ID, 由系统统一分配"
	AccountID     NAME_STRING "记录用户连接的账号--这里是连接的IP"
	Conn          *net.TCPConn
	Sender        *TCPSender
	dataChan      chan *Command "由TCPUserConn创建，用于登陆时交换userChan, 与User.netChan是同一对象"
	userChan      chan *Command "由User创建的channel用于网络模块传输数据包给User, 与User.recvChan是同一对象"
	_isLogin      bool
	_isConnection bool
}

func ProcessRecv(handler TCPHandler, isInner bool) {
	defer func() {
		if err := recover(); err != nil {
			LogError(err) //这里的err其实就是panic传入的内容
		}
	}()
	conn := handler.getConn()
	defer handler.close()

	for {
		//conn.SetReadDeadline(time.Now().Add(1 * 60 * time.Second))

		headerBytes := make([]byte, HEADER_LENGTH)
		_, err := io.ReadFull(conn, headerBytes)
		if err != nil {
			LogError("Read Data Error, maybe the socket is closed!  ", handler.getID())
			break
		}

		header := &PackHeader{}
		Byte2Struct(reflect.ValueOf(header), headerBytes)

		//LogDebug("Header", header.Cmd, header.Length, header.Tag, header.Version)
		bodyBytes := make([]byte, header.Length-HEADER_LENGTH)
		_, err = io.ReadFull(conn, bodyBytes)
		if err != nil {
			LogError("read data error ", err.Error())
			break
		}

		handler.processClientMessage(header, bodyBytes)
	}
}

func (client *TCPUserConn) processClientMessage(header *PackHeader, datas []byte) {
	if !client.isLogin() {
		client.processLogin(header, datas)
	} else {
		client.routMsgToUser(header, datas)
	}
}

func (client *TCPUserConn) processLogin(header *PackHeader, datas []byte) {
	defer func() {
		if err := recover(); err != nil {
			LogError(err)
		}
	}()

	if header.Cmd != CMD_CONNECTION && !client._isConnection {
		client.close()
		LogError("Wrong command", header.Cmd, " should be ", CMD_CONNECTION)
		return
	}
	if !client.isConnection() {
		go client.Sender.Start()
	}

	if header.Cmd == CMD_CONNECTION {
		client._isConnection = true

		msgConnection := &MsgConnection{}
		msgConnection.CreateByBytes(datas)
		client.Sender.Send(msgConnection)
		//LogDebug("MsgConnect成功...", msgConnection.ID, msgConnection.AccountID)
		return
	}
	if header.Cmd == CMD_REGISTER_SERVER {
		msgUserRegister := &MsgServerRegister{}
		//LogDebug("Server 注册...")
		msgUserRegister.CreateByBytes(datas)

		chanRet := make(chan ObjectID)
		go msgUserRegister.Process(chanRet)
		//等待返回UserID
		select {
		case id := <-chanRet:
			client.ID = id
			client.AccountID = msgUserRegister.Account
		case <-time.After(20 * time.Second):
			LogError("register put user channel failed:", header.Cmd)
			client.Sender.Send(msgUserRegister)
		}
	} else if header.Cmd == CMD_LOGIN {
		//LogDebug("MsgLogin成功...")
		client.ID = SysIDGenerator.GetNextID(ID_CLIENT)
		//分配一个ID
	} else {
		return
	}
	client._isLogin = true
	targetChan := GetChanByID(SYSTEM_USER_CHAN_ID)

	client.dataChan = make(chan *Command)
	msgSend := &Command{CMD_SYSTEM_USER_LOGIN, client.ID, client.dataChan, nil}
	msgSend.OtherInfo = client.Sender

	select {
	case targetChan <- msgSend:
	case <-time.After(5 * time.Second):
		LogError("loginUserToGame put user channel failed:", CMD_SYSTEM_USER_LOGIN)
	}
	client.waitLoginReturn()
}

func (client *TCPUserConn) waitLoginReturn() bool {
	msg := <-client.dataChan
	if msg.RetChan == nil {
		return false
	}
	client.userChan = msg.RetChan

	msgLogin := &MsgLogin{}
	msgLogin.UserID = client.ID
	msgLogin.Account = client.AccountID
	client.Sender.Send(msgLogin)
	return true
}

// 将消息路由到玩家处理
func (client *TCPUserConn) routMsgToUser(header *PackHeader, data []byte) bool {
	msg := &Command{header.Cmd, data, nil, nil}
	//LogDebug("routMsgToUser: ", Byte2String(client.AccountID[:]), client.ID)
	select {
	case client.userChan <- msg:
	case <-time.After(5 * time.Second):
		LogError("routMsgToUser put user channel failed:", client.ID)
		return false
	}

	return true
}

func (client *TCPUserConn) close() {
	client.Conn.Close()
	client.Sender.Close()
	close(client.dataChan)

	if !client._isLogin {
		return
	}

	userInnerChan := GetChanByID(client.ID)
	closeMsg := &Command{CMD_SYSTEM_USER_OFFLINE, nil, client.dataChan, nil}
	client._isLogin = false

	select {
	case userInnerChan <- closeMsg:
	case <-time.After(5 * time.Second):
		LogError("sendOffline put user channel failed:", client.ID)
		return
	}
	return
}

func (this *TCPUserConn) getID() ObjectID {
	return this.ID
}

func (this *TCPUserConn) getConn() *net.TCPConn {
	return this.Conn
}

func (this *TCPUserConn) getSender() *TCPSender {
	return this.Sender
}

func (this *TCPUserConn) getDataChan() chan *Command {
	return this.dataChan
}

func (this *TCPUserConn) getUserChan() chan *Command {
	return this.userChan
}

func (this *TCPUserConn) isLogin() bool {
	return this._isLogin
}

func (this *TCPUserConn) isConnection() bool {
	return this._isConnection
}
