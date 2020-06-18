package message

import (
	"fmt"
)

func init() {
	fmt.Println("message.init")
	registerNetMsgLogin()
	registerNetMsgConnection()
	registerNetMsgRegisterServer()
	registerNetMsgBindServer()
	registerNetMsgProject()
}
