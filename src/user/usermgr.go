package user

import (
	. "def"
	"time"

	. "ngcod.com/core"
)

var UserMgr UserManager

type UserManager struct {
	users      map[ObjectID]*User
	servers    map[NAME_STRING]*User
	systemChan chan *Command
}

func CreateUserMgr() bool {
	UserMgr = UserManager{}
	UserMgr.systemChan = make(chan *Command)
	UserMgr.users = make(map[ObjectID]*User)
	UserMgr.servers = make(map[NAME_STRING]*User)
	go startRecv(&UserMgr)
	return true
}

func startRecv(userMgr *UserManager) {
	RegisterChan(SYSTEM_USER_CHAN_ID, UserMgr.systemChan)
	defer UnRegisterChan(SYSTEM_USER_CHAN_ID)
	for {
		select {
		case msg := <-UserMgr.systemChan:
			userMgr.processMsg(msg)
		}
	}
}

func (this *UserManager) processMsg(msg *Command) {
	switch msg.Cmd {
	case CMD_SYSTEM_USER_LOGIN:
		this.processUserLogin(msg)
	case CMD_BIND_SERVER:
		this.processBindServer(msg)
	case CMD_SYSTEM_BROADCAST:
		this.processBroadCast(msg.Message.(NetMsg))
	}
}

func (this *UserManager) processUserLogin(msg *Command) {
	id := msg.Message.(ObjectID)

	oldUser, exist := this.users[id]

	if exist && !oldUser.IsServer && oldUser.Status == USER_STATUS_OFFLINE {
		exist = false
	}
	if !exist {
		oldUser = CreateUser(id)
	}
	this.users[id] = oldUser

	select {
	case oldUser.innerChan <- msg:
	case <-time.After(10 * time.Second):
		LogError("new user busy :", id)
		return
	}
}

func (this *UserManager) processBindServer(msg *Command) {
	id := msg.Message.(ObjectID)
	account := msg.OtherInfo.(NAME_STRING)
	AccountID := Byte2String(account[:])
	LogDebug("注册服务器. ", "ID=", id, ", Account=", AccountID)
	u, ok := this.users[id]
	if ok {
		this.servers[account] = u
	}
	select {
	case u.innerChan <- msg:
	case <-time.After(10 * time.Second):
		return
	}
}

func (this *UserManager) processBroadCast(msg NetMsg) {
	for _, u := range this.users {
		if u.Status == USER_STATUS_OFFLINE {
			continue
		}
		u.Sender.Send(msg)
		//LogDebug("BroadcastMessage to User,", u.ID, u.Status)
	}
}

func (this *UserManager) BroadcastMessage(msg NetMsg) {
	cmd := &Command{CMD_SYSTEM_BROADCAST, msg, nil, nil}
	select {
	case this.systemChan <- cmd:
	case <-time.After(20 * time.Second):
		return
	}
}
