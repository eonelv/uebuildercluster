package idmgr

import (
	. "core"
	. "mydb"
	"sync"
	"time"
)

const (
	ID_SERVER uint64 = 1000
	ID_CLIENT uint64 = 10000
)

type IDGenerator struct {
	IDChanSet map[uint64]chan ObjectID
	mutex     sync.RWMutex
}

var SysIDGenerator IDGenerator

func InitGenerator() {
	SysIDGenerator = IDGenerator{}
	SysIDGenerator.IDChanSet = make(map[uint64]chan ObjectID)
	SysIDGenerator.mutex = sync.RWMutex{}
	SysIDGenerator.create()
}

func generateID(idMax ObjectID, ch chan ObjectID) {
	for {
		ch <- idMax
		LogInfo("写入ID", idMax)
		idMax++
	}
}

func (this *IDGenerator) create() bool {
	idMax := ID_CLIENT
	this.IDChanSet[ID_CLIENT] = make(chan ObjectID, 1)
	go generateID(ObjectID(idMax+1), this.IDChanSet[ID_CLIENT])

	rows, err := DBMgr.Query("select max(id) as maxid from t_project")
	if err != nil || len(rows) == 0 {
		idMax := 1000
		this.IDChanSet[ID_SERVER] = make(chan ObjectID, 1)
		go generateID(ObjectID(idMax+1), this.IDChanSet[ID_SERVER])
		return false
	}
	LogDebug("row count = ", len(rows))
	for _, row := range rows {
		idMax := row.GetObjectID("maxid")
		this.IDChanSet[ID_SERVER] = make(chan ObjectID, 1)
		go generateID(ObjectID(idMax+1), this.IDChanSet[ID_SERVER])
		LogInfo("创建ID", idMax+1)
	}
	return true
}

func (this *IDGenerator) GetNextID(key uint64) ObjectID {
	var ch chan ObjectID
	var ok bool
	this.mutex.RLock()
	ch, ok = this.IDChanSet[key]
	this.mutex.RUnlock()
	if !ok {
		this.mutex.Lock()
		ch, ok = this.IDChanSet[key]
		if !ok {
			idMax := 1000
			this.IDChanSet[ID_SERVER] = make(chan ObjectID, 1)
			go generateID(ObjectID(idMax+1), this.IDChanSet[ID_SERVER])
		}
		this.mutex.Unlock()
	}
	var id ObjectID
	select {
	case id = <-ch:
	case <-time.After(20 * time.Second):
		LogError("get Next ID time out")
		return ObjectID(0)
	}

	return id
}
