// serverconfig
package cfg

import (
	"fmt"
	"os"
	"path/filepath"

	. "ngcod.com/core"
	"ngcod.com/utils"
)

type ProcessActiveData struct {
	PID       int    //进程ID
	ParentPID int    //父进程ID
	Name      string //进程名
}

type TConfigData struct {
	Path           string //启动程序路徑
	ProcessName    string //原始進程名
	ActProcessName string //需要检查的进程名
	NewWindow      bool
}

type Config struct {
	BuilderHome string
	Datas       map[string]TConfigData
	ActiveDatas map[string]*ProcessActiveData
}

const serverconfig string = `{
	"uebuildtoolart":{"path":"E:/uebuilderhome/uebuilder/UEBuilderTool.exe", "process":"uebuildtool"}
}`

func (this *Config) ReadConfig() error {
	if this.ActiveDatas == nil {
		this.ActiveDatas = make(map[string]*ProcessActiveData)
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		this.BuilderHome = dir
	} else {
		this.BuilderHome = "E:/golang/uebuildtool"
	}

	configHome := fmt.Sprintf("%s/config", this.BuilderHome)
	utils.PathExistAndCreate(configHome)
	configFileName := configHome + "/config.json"

	oldJson, err := utils.ReadJson(configFileName)
	if err != nil {
		LogError("Read config Json Failed! 1.")
		utils.WriteFile([]byte(serverconfig), configFileName)
		return err
	}
	ConfigDatas := oldJson.MustMap()

	this.Datas = make(map[string]TConfigData)
	for k, _ := range ConfigDatas {
		configData := TConfigData{}
		configData.Path = utils.GetTableString(ConfigDatas, k, "path")
		configData.ProcessName = utils.GetTableString(ConfigDatas, k, "process")
		configData.ActProcessName = utils.GetTableString(ConfigDatas, k, "process2")
		configData.NewWindow = utils.GetTableBool(ConfigDatas, k, "newWindow")
		this.Datas[k] = configData
		LogDebug("Config Data. ", configData)
	}
	return nil
}
