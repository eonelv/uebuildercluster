// unrealeditorAuth
package unrealeditor

import (
	"fmt"
	"io"
	"net"
	. "netcore"
	"reflect"

	. "ngcod.com/core"
)

var Tags []string = []string{
	"Xzylo0Ba1", "Handsome005Boy", "Release138MySoul", "Beautiful365Gril", "ABzslo0Bo8",
	"EoneMat002zsl", "ABlvlo0Bllp8", "ABDfa", "MMcdn8875", "BBCCofe",
	"98knxjgnkds", "lkjhgbn5264hj", "Coffee", "CheckCode", "ToWork",
	"EveryoneIsNo1", "ShowTime", "AlwaysGoto", "ThankAndThink", "LoveIsDirty"}

type MsgUnrealAuth struct {
	Index int64
	Value [100]byte
}

func (this *MsgUnrealAuth) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(20000), reflect.ValueOf(this))
}

func (this *MsgUnrealAuth) CreateByBytes(bytes []byte) (bool, int) {
	return Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgUnrealAuth) Process(p interface{}) {

}

func StartUnrealEditorAuth() {
	defer func() {
		if err := recover(); err != nil {
			LogError(err) //这里的err其实就是panic传入的内容
		}
	}()

	service := fmt.Sprintf(":%d", 6001)
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

	client.Conn = conn
	client.Sender = CreateTCPSender(conn)
	go client.Sender.Start()

	go processRecv(client)
}

func processRecv(handler *TCPUserConn) {
	defer func() {
		if err := recover(); err != nil {
			LogError(err) //这里的err其实就是panic传入的内容
		}
	}()
	conn := handler.Conn
	defer handler.Close()

	for {
		headerBytes := make([]byte, HEADER_LENGTH)
		_, err := io.ReadFull(conn, headerBytes)
		if err != nil {
			//LogError("Read Data Error, maybe the socket is closed!  ")
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
		message := &MsgUnrealAuth{}
		message.CreateByBytes(bodyBytes)
		CopyArray(reflect.ValueOf(&message.Value), []byte(getCode(message.Index)))
		handler.Sender.Send(message)
	}
}

func getCode(Value int64) string {
	index := getIndex(Value)
	if index == -1 {
		return ""
	}
	if index == 13 {
		return ""
	}
	return Tags[index]
}

func getIndex(Value int64) int64 {
	var temp int64 = -1
	if Value == -1 {
		temp = -1
	} else {
		if Value%2 == 0 {
			temp = Value % 7
		} else {
			temp = Value % 15
		}
		if temp != 0 {
			Value += Value / temp
		}
		temp = Value % 20
	}
	return temp
}

func checkError(err error) {
	if err != nil {
		LogError(err)
	}
}
