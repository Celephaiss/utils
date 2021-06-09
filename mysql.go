package utils

import (
	"fmt"
)

func ConstructDataSourceName(host string, port int, user, pass, database string) string {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, pass, host, port, database)
	// a2_only_read:a2_only_read@tcp(10.160.86.130:3306)/a2?charset=utf8
	return dataSourceName
}

//func QueryMySQL(query, dataSourceName string) (*sqlx.DB, error) {
//
//	db, err := sqlx.Open("mysql", dataSourceName)
//
//	if err != nil {
//		return nil, err
//	}
//
//	//defer func(db *sqlx.DB) {
//	//	err := db.Close()
//	//	if err != nil {
//	//		return
//	//	}
//	//}(db)
//
//	return db, nil
//	//var results []*interface{}
//
//	//err = db.Select(&results, query, Yesterday())
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//return results, nil
//}
