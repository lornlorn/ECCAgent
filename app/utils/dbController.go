package utils

import (
    "github.com/cihub/seelog"
    "github.com/go-xorm/xorm"       // _ "github.com/lib/pq" // postgresql driver
    _ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

// Engine xorm
var Engine *xorm.Engine

// InitDB func() error
func InitDB(dbtype string, dbstr string) error {
    var err error
    // dbtype := "postgres"
    // dbstr := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, user, password, dbname)
    // dbstr := "postgres://test:test@localhost:5432/testdb?sslmode=disable"
    // dbstr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)

    seelog.Debugf("Database Type : [%v], Connect String : [%v]", dbtype, dbstr)

    Engine, err = xorm.NewEngine(dbtype, dbstr)
    if err != nil {
        return err
    }
    Engine.ShowSQL(true)

    err = Engine.Ping()
    if err != nil {
        return err
    }
    Engine.Exec("PRAGMA foreign_keys = ON")

    return nil
}
