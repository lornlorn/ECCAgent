package models

import (
    "app/utils"

    "github.com/cihub/seelog"
)

/*
SysLog struct map to table sys_log
*/
type SysLog struct {
    LogId       int    `xorm:"INTEGER NOT NULL UNIQUE PK"`
    RunTime     string `xorm:"VARCHAR(19)   NOT NULL"`
    RunEnvs     string `xorm:"VARCHAR(512)"`
    RunCmd      string `xorm:"VARCHAR(512)   NOT NULL"`
    RunArgs     string `xorm:"VARCHAR(512)"`
    RunStatus   string `xorm:"VARCHAR(16)   NOT NULL"`
    RunMsg      string `xorm:"VARCHAR(512)"`
    LogfilePath string `xorm:"VARCHAR(512)   NOT NULL"`
    ReqSrc      string `xorm:"VARCHAR(32)   NOT NULL"`
}

/*
NewLog struct map to table sys_log without column Id
*/
type NewLog struct {
    RunTime     string `xorm:"VARCHAR(19)   NOT NULL"`
    RunEnvs     string `xorm:"VARCHAR(512)"`
    RunCmd      string `xorm:"VARCHAR(512)   NOT NULL"`
    RunArgs     string `xorm:"VARCHAR(512)"`
    RunStatus   string `xorm:"VARCHAR(16)   NOT NULL"`
    RunMsg      string `xorm:"VARCHAR(512)"`
    LogfilePath string `xorm:"VARCHAR(512)   NOT NULL"`
    ReqSrc      string `xorm:"VARCHAR(32)   NOT NULL"`
}

/*
TableName xorm mapper
NewComponent struct map to table tb_component
*/
func (nlog NewLog) TableName() string {
    return "sys_log"
}

// Save insert method
func (nlog NewLog) Save() error {
    affected, err := utils.Engine.Insert(nlog)
    if err != nil {
        return err
    }
    seelog.Debugf("%v insert : %v", affected, nlog)

    return nil
}
