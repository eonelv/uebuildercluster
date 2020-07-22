package mydb

import (
	. "core"
	SQL "database/sql"
	. "db/sqlitedb"
	"os"
	"strings"
)

const initSql string = `CREATE TABLE [t_project] (
  [id] INTEGER NOT NULL PRIMARY KEY,
  [name] VARCHAR(255) NOT NULL,
  [projectName] VARCHAR(255) NOT NULL,
  [host] VARCHAR(255) NOT NULL,
  [account] VARCHAR(255) NOT NULL,
  [svn] VARCHAR(255) NOT NULL,
  [member] VARCHAR(255) NOT NULL,
  [buildstep] VARCHAR(1024),
  [serverState] INTEGER DEFAULT 0);
  insert into t_project (id, name, projectName, host, account, svn, member, buildstep) values (1000, 'test', 'projectName', 'host', 'account', 'svn', 'member', 'buildstep');`

var DBMgr DataBaseMgr

func CreateDBMgr(path string) bool {
	dbExist := CreateDB(path)

	db, err := SQL.Open("sqlite3", path)
	//	db, err := SQL.Open("mysql", "ouyang:ouyang@tcp(192.168.0.10:3306)/zentao?charset=utf8")
	if err != nil {
		LogError("DataBase Connect Error %s \n", err.Error())
		return false
	}
	DBMgr = DataBaseMgr{}
	DBMgr.DBServer = db
	if !dbExist {
		InitDB()
	}
	return true
}

func DBExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDB(path string) bool {
	isDBExist, _ := DBExists(path)
	if !isDBExist {
		fileDest, errDest := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModeAppend)
		if errDest != nil {
			LogError(errDest)
			return false
		}
		defer fileDest.Close()
		return false
	}
	return true
}

func InitDB() bool {
	LogInfo("create database and init db")
	var result bool
	result = initSQL()
	return result
}

func initSQL() bool {

	lineArray := strings.Split(initSql, ";")

	var errSQL error
	var errorlog string = "SQL init error: "
	var hasError bool
	for _, line := range lineArray {
		if strings.TrimSpace(line) == "" {
			continue
		}
		_, errSQL = DBMgr.Execute(line)
		if errSQL != nil {
			hasError = true
			errorlog += line
		}
	}
	if hasError {
		LogError(errorlog)
	}
	return true
}
