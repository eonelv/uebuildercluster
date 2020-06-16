// init
package netcore

import (
	"fmt"
)

func init() {
	fmt.Println("netcore.init")
	registerNetMsgBuild()
}
