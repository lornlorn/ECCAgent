package models

import (
    "app/utils"
    "errors"

    "github.com/cihub/seelog"
)

/*
SysCron struct map to table sys_cron
*/
type SysCron struct {
    CronId      int    `xorm:"'CRON_ID' INTEGER NOT NULL UNIQUE PK"`
    CronName    string `xorm:"'CRON_NAME' VARCHAR(256) NOT NULL UNIQUE"`
    CronType    string `xorm:"'CRON_TYPE' VARCHAR(16) NOT NULL"`
    CronSpec    string `xorm:"'CRON_SPEC' VARCHAR(128) NOT NULL"`
    CronHost    string `xorm:"'CRON_HOST' VARCHAR(512)"`
    CronAuth    string `xorm:"'CRON_AUTH' VARCHAR(512)"`
    CronPrivkey string `xorm:"'CRON_PRIVKEY' VARCHAR(512)"`
    CronCmd     string `xorm:"'CRON_CMD' VARCHAR(512) NOT NULL"`
    CronStatus  string `xorm:"'CRON_STATUS' VARCHAR(16) NOT NULL"`
    CronDesc    string `xorm:"'CRON_DESC' VARCHAR(1024)"`
    CronHx      string `xorm:"'CRON_HX' VARCHAR(1024)"`
}

/*
TableName xorm mapper
SysCron struct map to table SYS_CRON
*/
func (cron SysCron) TableName() string {
    return "SYS_CRON"
}

/*
NewCron struct map to table sys_cron without column Id
*/
type NewCron struct {
    CronName    string `xorm:"VARCHAR(256) NOT NULL UNIQUE"`
    CronType    string `xorm:"VARCHAR(16) NOT NULL"`
    CronSpec    string `xorm:"VARCHAR(128) NOT NULL"`
    CronHost    string `xorm:"VARCHAR(512)"`
    CronAuth    string `xorm:"VARCHAR(512)"`
    CronPrivkey string `xorm:"VARCHAR(512)"`
    CronCmd     string `xorm:"VARCHAR(512)   NOT NULL"`
    CronStatus  string `xorm:"VARCHAR(16)   NOT NULL"`
    CronDesc    string `xorm:"VARCHAR(1024)"`
    CronHx      string `xorm:"VARCHAR(1024)"`
}

/*
TableName xorm mapper
NewCron struct map to table SYS_CRON
*/
func (cron NewCron) TableName() string {
    return "SYS_CRON"
}

// Save insert method
func (cron NewCron) Save() error {
    affected, err := utils.Engine.Insert(cron)
    if err != nil {
        return err
    }
    seelog.Debugf("%v insert : %v", affected, cron)

    return nil
}

// Delete method
func (cron NewCron) Delete() error {
    affected, err := utils.Engine.Where("cron_name=?", cron.CronName).Delete(cron)
    if err != nil {
        return err
    }
    seelog.Debugf("%v delete : %v", affected, cron)

    if affected == 0 {
        return errors.New("0 updated, nothing changed")
    }

    return nil
}

// UpdateByUUID method
func (cron NewCron) UpdateByUUID() error {
    affected, err := utils.Engine.Where("cron_name=?", cron.CronName).Update(cron)
    if err != nil {
        // seelog.Errorf("utils.Engine.ID.Update Error : %v", err)
        return err
    }
    seelog.Debugf("%v update : %v", affected, cron)

    if affected == 0 {
        return errors.New("0 updated, nothing changed")
    }

    return nil
}

/*
GetCrons func() ([]SysCron, error)
*/
func GetCrons() ([]SysCron, error) {
    crons := make([]SysCron, 0)

    if err := utils.Engine.Where("cron_status=?", "启用").Find(&crons); err != nil {
        // seelog.Errorf("utils.Engine.Find Error : %v", err)
        return nil, err
    }

    return crons, nil
}

/*
DelCronByKeys func(cronID int, cronName string, cronUUID string) error
*/
func DelCronByKeys(cronID int, cronName string) error {
    cron := new(SysCron)
    cron.CronId = cronID
    cron.CronName = cronName

    affected, err := utils.Engine.Delete(cron)
    if err != nil {
        // seelog.Errorf("utils.Engine.Delete Error : %v", err)
        return err
    }
    seelog.Debugf("%v delete : %v", affected, cron)

    return nil
}
