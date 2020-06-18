// cmddef
package def

const (
	CMD_CONNECTION uint16 = 1052
	CMD_REGISTER   uint16 = 1002
	CMD_LOGIN      uint16 = 1003
	CMD_TALK       uint16 = 1006

	CMD_REGISTER_SERVER uint16 = 50007
	CMD_BIND_SERVER     uint16 = 50008
	CMD_BUILD           uint16 = 50010
	CMD_BUILD_INFO      uint16 = 50011
	CMD_PROJECT         uint16 = 50012

	CMD_ACTOR uint16 = 10008
)

const (
	CMD_SYSTEM_MAIN_CLOSE   uint16 = 10001
	CMD_SYSTEM_USER_OFFLINE uint16 = 10002
	CMD_SYSTEM_USER_LOGIN   uint16 = 10005
	CMD_SYSTEM_BROADCAST    uint16 = 10006
	CMD_SYSTEM_USER_MSG     uint16 = 10007
	CMD_SYSTEM_SERVER_BUILD uint16 = 10008
)

const (
	CMD_ROOM_CLOSE          uint16 = 20001
	CMD_ROOM_FULL           uint16 = 20002
	CMD_ROOM_START          uint16 = 20003
	CMD_ROOM_USER_READY     uint16 = 20004
	CMD_ROOM_USER_PLAY_CARD uint16 = 20005
	CMD_ROOM_USER_TEST      uint16 = 20006
)