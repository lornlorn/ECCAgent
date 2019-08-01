package scheduler

import (
    "app/models"
    "app/utils"
    "github.com/cihub/seelog"
    "net"
    "time"
)

func CronPortScan(cron models.SysCron) error {
    UUID := utils.GetUniqueID()
    seelog.Infof("[%v]PORT->[%v] Begin ...", UUID, cron.CronCmd)

    conn, err := net.DialTimeout("tcp", cron.CronCmd, 5*time.Second)
    if err != nil {
        seelog.Errorf("[%v]ERROR : %v", UUID, err.Error())
        utils.SendHXMsg("端口探测失败通知", cron.CronHx, cron.CronCmd)
        return err
    }
    defer conn.Close()

    seelog.Infof("[%v]>>> Check [%v] OK! <<<", UUID, cron.CronCmd)
    seelog.Infof("[%v]PORT->[%v] Finish ...", UUID, cron.CronCmd)

    return nil
}
