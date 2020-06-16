package sqlitedb

import (
	. "core"
	SQL "database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DataBaseMgr struct {
	DBServer *SQL.DB
}

func (this *DataBaseMgr) Execute(sql string) (SQL.Result, error) {
	return this.DBServer.Exec(sql)
}

func (this *DataBaseMgr) PreExecute(sql string, args ...interface{}) (SQL.Result, error) {
	LogInfo(sql, args)
	stmt, err := this.DBServer.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}

func (this *DataBaseMgr) Query(sql string, args ...interface{}) ([]*RowSet, error) {
	rows, err := this.DBServer.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return this.buildRowSets(rows)
}

func (this *DataBaseMgr) PreQuery(sql string, args ...interface{}) ([]*RowSet, error) {
	LogInfo(sql, args)
	stmt, err := this.DBServer.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, errq := stmt.Query(args...)

	if errq != nil {
		return nil, errq
	}
	defer rows.Close()

	return this.buildRowSets(rows)
}

func (this *DataBaseMgr) buildRowSets(rows *SQL.Rows) ([]*RowSet, error) {
	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var results []*RowSet = []*RowSet{}
	var rowSet *RowSet
	for rows.Next() {
		values := make([][]byte, len(colNames))
		scans := make([]interface{}, len(colNames))
		for i, _ := range values {
			scans[i] = &values[i]
		}
		if err := rows.Scan(scans...); err != nil {
			return nil, err
		}
		rowSet = &RowSet{}
		rowSet.Datas = make(map[string][]byte)
		rowSet.Cols = colNames
		for j, v := range values {
			key := colNames[j]
			rowSet.Datas[key] = v
		}
		results = append(results, rowSet)
	}
	return results, nil
}
