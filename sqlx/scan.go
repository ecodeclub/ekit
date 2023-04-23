package sqlx

import (
	"database/sql"
	"reflect"
	"time"
)

func ScanRows(rows *sql.Rows) ([]any, error) {
	//defer rows.Close()
	colsInfo, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	colsData := make([]any, 0, len(colsInfo))
	for _, colInfo := range colsInfo {
		typ := colInfo.ScanType()
		// 保险起见，循环的去除指针
		for typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		newData := reflect.New(typ).Interface()
		colsData = append(colsData, newData)
	}
	err = rows.Scan(colsData...)
	if err != nil {
		return nil, err
	}
	// 去掉reflect.New的指针
	for i := 0; i < len(colsData); i++ {
		colsData[i] = reflect.ValueOf(colsData[i]).Elem().Interface()
	}
	return colsData, nil
}

func ScanAll(rows *sql.Rows) ([][]any, error) {
	// TODO 暂定结果数32
	res := make([][]any, 0, 32)
	for rows.Next() {
		cols, err := ScanRows(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, cols)
	}
	return res, nil
}

func ScanRows_(rows *sql.Rows) ([]interface{}, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columns := make([]interface{}, len(colTypes))
	columnPtrs := make([]interface{}, len(colTypes))
	for i := range columns {
		columnPtrs[i] = &columns[i]
	}
	for i, colType := range colTypes {
		switch colType.DatabaseTypeName() {
		case "INT", "INTEGER", "TINYINT", "SMALLINT", "MEDIUMINT":
			columns[i] = new(int64)
		case "BIGINT":
			columns[i] = new(uint64)
		case "FLOAT", "DOUBLE":
			columns[i] = new(float64)
		case "DECIMAL":
			columns[i] = new(string)
		case "DATE", "TIME", "YEAR", "DATETIME", "TIMESTAMP":
			columns[i] = new(time.Time)
		default:
			columns[i] = new(string)
		}
	}
	results := make([]interface{}, len(columns))
	for i := range results {
		results[i] = &columns[i]
	}
	for rows.Next() {
		err := rows.Scan(results...)
		if err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return columns, nil
}
