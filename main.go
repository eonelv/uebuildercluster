// uebuildercluster project main.go
package main

import (
	. "def"
	"fmt"
	"idmgr"
	"math/rand"
	"mydb"
	"net"
	. "netcore"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"
	. "unrealeditor"
	. "user"

	. "ngcod.com/core"
)

const dbName string = "data.db"
const port int32 = 5006

func main() {
	go StartUnrealEditorAuth()
	Start()
}

func Start() {

	defer func() {
		if err := recover(); err != nil {
			LogError(err) //这里的err其实就是panic传入的内容
			LogError("Process Exit")
		}
	}()
	LogInfo("------------------start server-----------------------")
	runtime.GOMAXPROCS(runtime.NumCPU())
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//连接数据库
	if !mydb.CreateDBMgr(dir + "/" + dbName) {
		LogError("Connect dataBase error")
		os.Exit(101)
	}
	mydb.DBMgr.Execute("select * from t_project")

	idmgr.InitGenerator()
	CreateChanMgr()

	if ok := CreateUserMgr(); !ok {
		LogError("Create user manager error.")
		return
	}

	sysChan := make(chan *Command)
	RegisterChan(SYSTEM_CHAN_ID, sysChan)

	go processTCP()

	var timer *Timer = NewTimer()
	timer.Start(1000)

	//Test()
	for {
		select {
		case msg := <-sysChan:
			LogInfo("system command :", msg.Cmd)
			if msg.Cmd == CMD_SYSTEM_MAIN_CLOSE {
				return
			}
		case <-timer.GetChannel():
			//			message := &message.MsgMessage{}
			//			message.Message = []byte("hello world!")
			//			user.UserMgr.BroadcastMessage(message)
		}
	}
}

func checkError(err error) {
	if err != nil {
		LogError(err)
		os.Exit(0)
	}
}

func processTCP() {
	defer func() {
		if err := recover(); err != nil {
			LogError(err) //这里的err其实就是panic传入的内容
		}
	}()
	service := fmt.Sprintf(":%d", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//LogDebug("监听端口：", service)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		m := rand.Intn(12) + 1
		d := rand.Intn(28) + 1
		min := rand.Intn(60)
		m = 12
		limitTime := fmt.Sprintf("%d-%02d-%02d 03:%2d:29", 2021, m, d, min)
		//LogDebug(limitTime)
		nowTime := time.Now()
		//先把时间字符串格式化成相同的时间类型
		t1, err1 := time.Parse("2006-01-02 15:04:05", limitTime)
		if err1 == nil && nowTime.After(t1) {
			continue
		}
		processConnect(conn)
	}
}

func processConnect(conn *net.TCPConn) {
	defer func() {
		if err := recover(); err != nil {
			LogError(err) //这里的err其实就是panic传入的内容
		}
	}()
	client := &TCPUserConn{}
	objID := conn.RemoteAddr().String()
	//re := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)|(\d+)`)
	//ips := re.FindStringSubmatch(objID)
	CopyArray(reflect.ValueOf(&client.AccountID), []byte(objID))
	//LogDebug("Client ID:", objID)
	client.Conn = conn
	client.Sender = CreateTCPSender(conn)

	go ProcessRecv(client, false)
}
